package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/amb1s1/gonetconfig/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to %v, error: %v", address, err)
	}
	defer conn.Close()
	c := pb.NewGoNetConfigServiceClient(conn)

	request := &pb.ConfigGenRequest{
		ConfigGenType: pb.ConfigGenType_CGT_FULL.Enum(),
		Device: &pb.Device{
			Name:   proto.String("router1"),
			Vendor: pb.Vendor_VD_CISCO.Enum(),
			Model:  pb.Model_MD_UNKNOW.Enum(),
		},
	}
	response, err := c.GetConfigGen(ctx, request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("results: ", response.GetConfigFeature().GetConfiguration())
	fmt.Println("status: ", response.GetStatus())
}
