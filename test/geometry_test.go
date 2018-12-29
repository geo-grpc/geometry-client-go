/*
Copyright 2017 Echo Park Labs

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

For additional information, contact:

email: info@echoparklabs.io
*/

package test

import (
	"flag"
	"testing"
	pb "geometry-client-go/epl/protobuf/geometry"
	"google.golang.org/grpc/testdata"
	"google.golang.org/grpc/credentials"
	"log"
	"google.golang.org/grpc"
	"context"

)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:8980", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com",
		"The server name use to verify the hostname returned by TLS handshake")
)

func TestGeometryRequests(t *testing.T) {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGeometryOperatorsClient(conn)

	spatialReferenceWGS := pb.SpatialReferenceData{Wkid:4326}
	spatialReferenceNAD27 := pb.SpatialReferenceData{Wkid:4267}
	spatialReferenceMerc := pb.SpatialReferenceData{Wkid:3857}
	spatialReferenceGall := pb.SpatialReferenceData{Wkid:54016}
	spatialReferenceMoll := pb.SpatialReferenceData{Proj4:"+proj=moll +lon_0=0 +x_0=0 +y_0=0 +datum=WGS84 +units=m +no_defs"}

	geometry_string := []string{"MULTILINESTRING ((-120 -45, -100 -55, -90 -63, 0 0, 1 1, 100 25, 170 45, 175 65))"}

	// define a geometry from an array of wkt strings (in this case only one geometry) with spatial reference nad83
	lefGeometryBag := pb.GeometryBagData{
		Wkt: geometry_string,
		SpatialReference:&spatialReferenceNAD27}

	// define a buffer opertor on geometry then buffered with distance size .5 degrees (I know horrible), and then
	// the result is transformed to WGS84
	operatorLeft := pb.OperatorRequest{
		GeometryBag:&lefGeometryBag,
		OperatorType:pb.ServiceOperatorType_Buffer,
		BufferParams:&pb.BufferParams{Distances:[]float64{.5}},
		ResultSpatialReference:&spatialReferenceWGS}

	// nest the result of previous buffer operation as the geometry input for this operation,
	// project that resulting buffered geometry to World Mollweide, perform the convex hull operation then
	// project that result to World Gall Stereo
	operatorNestedLeft := pb.OperatorRequest{
		GeometryRequest:&operatorLeft,
		OperatorType:pb.ServiceOperatorType_ConvexHull,
		OperationSpatialReference:&spatialReferenceMoll,
		ResultSpatialReference:&spatialReferenceGall}

	// define a geometry from wkt string with spatial reference nad83
	rightGeometryBag := pb.GeometryBagData{
		Wkt: geometry_string,
		SpatialReference:&spatialReferenceNAD27}

	// Project the geometry to WGS84 and then perform a geodesic buffer of the input geometry with a distance of 1000 meters.
	// The parameters for the geodesic buffer will be derived from the WGS84 spatial reference (the resulting spatial reference will be wgs84)
	operatorRight := pb.OperatorRequest{
		GeometryBag:&rightGeometryBag,
		OperatorType:pb.ServiceOperatorType_GeodesicBuffer,
		BufferParams:&pb.BufferParams{
			Distances:[]float64{1000},
			UnionResult:false},
		OperationSpatialReference:&spatialReferenceWGS}

	// Perform a convex hull operation on the previous buffer operation's result. And then project to Gall
	operatorNestedRight := pb.OperatorRequest{
		GeometryRequest:&operatorRight,
		OperatorType:pb.ServiceOperatorType_ConvexHull,
		ResultSpatialReference:&spatialReferenceGall}

	// take each of the nested buffer + convex hull and test that the non-geodesic contains the geodesic
	operatorContains := pb.OperatorRequest{
		LeftGeometryRequest:&operatorNestedLeft,
		RightGeometryRequest:&operatorNestedRight,
		OperatorType:pb.ServiceOperatorType_Contains,
		OperationSpatialReference:&spatialReferenceMerc}
	operatorResultEquals, err := client.ExecuteOperation(context.Background(), &operatorContains)

	result := operatorResultEquals.RelateMap[0]

	if result != true {
		t.Errorf("left nested request geometry should container right geometry nested request\n")
	}
}