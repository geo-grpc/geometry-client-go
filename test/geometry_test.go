package test

import (
	"flag"
	"testing"
	pb "geo-grpc/geometry-client-go/epl/geometry"
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
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func TestNestedRequests(t *testing.T) {
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
	/*

        OperatorRequest nestedLeft = OperatorRequest
                .newBuilder()
                .setLeftNestedRequest(serviceOpLeft)
                .setOperatorType(ServiceOperatorType.ConvexHull)
                .setResultSpatialReference(spatialReferenceGall)
                .build();

        OperatorRequest serviceOpRight = OperatorRequest
                .newBuilder()
                .setLeftGeometryBag(geometryBagRight)
                .setOperatorType(ServiceOperatorType.GeodesicBuffer)
                .setBufferParams(OperatorRequest.BufferParams.newBuilder().addDistances(1000).setUnionResult(false).build())
                .setOperationSpatialReference(spatialReferenceWGS)
                .build();
        OperatorRequest nestedRight = OperatorRequest
                .newBuilder()
                .setLeftNestedRequest(serviceOpRight)
                .setOperatorType(ServiceOperatorType.ConvexHull)
                .setResultSpatialReference(spatialReferenceGall)
                .build();

        OperatorRequest operatorRequestContains = OperatorRequest
                .newBuilder()
                .setLeftNestedRequest(nestedLeft)
                .setRightNestedRequest(nestedRight)
                .setOperatorType(ServiceOperatorType.Contains)
                .setOperationSpatialReference(spatialReferenceMerc)
                .build();

        GeometryOperatorsGrpc.GeometryOperatorsBlockingStub stub = GeometryOperatorsGrpc.newBlockingStub(inProcessChannel);
        OperatorResult operatorResult = stub.executeOperation(operatorRequestContains);
        Map<Integer, Boolean> map = operatorResult.getRelateMapMap();

        assertTrue(map.get(0));
	        SpatialReferenceData spatialReferenceNAD = SpatialReferenceData.newBuilder().setWkid(4269).build();
        SpatialReferenceData spatialReferenceMerc = SpatialReferenceData.newBuilder().setWkid(3857).build();
        SpatialReferenceData spatialReferenceWGS = SpatialReferenceData.newBuilder().setWkid(4326).build();
        SpatialReferenceData spatialReferenceGall = SpatialReferenceData.newBuilder().setWkid(54016).build();
	 */


	spatialReferenceWGS := pb.SpatialReferenceData{Wkid:4326}
	spatialReferenceNAD := pb.SpatialReferenceData{Wkid:4269}//SpatialReferenceData.newBuilder().setWkid(4269).build();
	spatialReferenceMerc := pb.SpatialReferenceData{Wkid:3857}
	spatialReferenceGall := pb.SpatialReferenceData{Wkid:54016}

	//var polyline = "MULTILINESTRING ((-120 -45, -100 -55, -90 -63, 0 0, 1 1, 100 25, 170 45, 175 65))";
	geometry_string := []string{"MULTILINESTRING ((-120 -45, -100 -55, -90 -63, 0 0, 1 1, 100 25, 170 45, 175 65))"}
	lefGeometryBag := pb.GeometryBagData{
		Wkt: geometry_string,
		GeometryEncodingType:pb.GeometryEncodingType_wkt,
		SpatialReference:&spatialReferenceNAD}

	operatorLeft := pb.OperatorRequest{
		GeometryBag:&lefGeometryBag,
		OperatorType:pb.ServiceOperatorType_Buffer,
		BufferParams:&pb.BufferParams{Distances:[]float64{.5}},
		ResultSpatialReference:&spatialReferenceWGS}

	operatorNestedLeft := pb.OperatorRequest{
		LeftGeometryRequest:&operatorLeft,
		OperatorType:pb.ServiceOperatorType_ConvexHull,
		ResultSpatialReference:&spatialReferenceGall}

	rightGeometryBag := pb.GeometryBagData{
		Wkt: geometry_string,
		GeometryEncodingType:pb.GeometryEncodingType_wkt,
		SpatialReference:&spatialReferenceNAD}

	operatorRight := pb.OperatorRequest{
		GeometryBag:&rightGeometryBag,
		OperatorType:pb.ServiceOperatorType_GeodesicBuffer,
		BufferParams:&pb.BufferParams{
			Distances:[]float64{1000},
			UnionResult:false},
		OperationSpatialReference:&spatialReferenceWGS}

	operatorNestedRight := pb.OperatorRequest{
		LeftGeometryRequest:&operatorRight,
		OperatorType:pb.ServiceOperatorType_ConvexHull,
		ResultSpatialReference:&spatialReferenceGall}

	operatorContains := pb.OperatorRequest{
		LeftGeometryRequest:&operatorNestedLeft,
		RightGeometryRequest:&operatorNestedRight,
		OperatorType:pb.ServiceOperatorType_Contains,
		OperationSpatialReference:&spatialReferenceMerc}
	operatorResultEquals, err := client.ExecuteOperation(context.Background(), &operatorContains)
	log.Println(operatorResultEquals.SpatialRelationship)

	result := operatorResultEquals.RelateMap[0]

	if result != true {
		t.Errorf("left nested request geometry should container right geometry nested request\n")
	}
}