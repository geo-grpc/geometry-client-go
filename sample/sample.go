package main

import (
	"flag"
	"google.golang.org/grpc/testdata"
	"google.golang.org/grpc/credentials"
	"log"
	"google.golang.org/grpc"
	pb "geometry-client-go/epl/geometry"
	"context"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:8980", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
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
	//sample := pb.Operator
	//sample := pb.NewRouteGuideClient(conn)
	geometry_string := []string{"MULTIPOLYGON (((40 40, 20 45, 45 30, 40 40)), ((20 35, 45 20, 30 5, 10 10, 10 30, 20 35), (30 20, 20 25, 20 15, 30 20)))"}
	serviceGeometry := pb.GeometryBagData{GeometryStrings: geometry_string, GeometryEncodingType:pb.GeometryEncodingType_wkt}
	cutterGeometry := pb.GeometryBagData{GeometryStrings: []string{"LINESTRING(0 0, 45 45)"}, GeometryEncodingType:pb.GeometryEncodingType_wkt}

	operatorCut := &pb.OperatorRequest{LeftGeometryBag:&serviceGeometry, RightGeometryBag:&cutterGeometry}
	//operatorCut := &pb.OperatorRequest{
	//	PrimaryGeoms:&pb.OperatorRequest_LeftGeometryBag{&serviceGeometry},
	//	SecondaryGeoms:&pb.OperatorRequest_RightGeometryBag{&cutterGeometry}}
	//operator := pb.OperatorRequest{PrimaryGeoms:serviceGeometry, SecondaryGeoms:cutterGeometry}
	//operator := pb.OperatorRequest{&serviceGeometry, &cutterGeometry}
	//operator := pb.OperatorRequest{LeftGeometryBag:&serviceGeometry, RightGeometryBag:&cutterGeometry}
	operatorCut.OperatorType = pb.ServiceOperatorType_Cut
	operatorCut.ResultsEncodingType = pb.GeometryEncodingType_geojson

	operatorResult, err := client.ExecuteOperation(context.Background(), operatorCut)


	log.Println(operatorResult.GeometryBag.GeometryStrings[0])

	operatorEquals := pb.OperatorRequest{ LeftNestedRequest:operatorCut, RightGeometryBag:&serviceGeometry, OperatorType:pb.ServiceOperatorType_Equals}
	//operatorEquals := pb.OperatorRequest{PrimaryGeoms:&pb.OperatorRequest_LeftNestedRequest{operatorCut}, SecondaryGeoms:&pb.OperatorRequest_RightGeometryBag{&serviceGeometry}}
	operatorResultEquals, err := client.ExecuteOperation(context.Background(), &operatorEquals)
	log.Println(operatorResultEquals.SpatialRelationship)
}
