package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/amb1s1/gonetconfig/proto"
	"google.golang.org/grpc"
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
		ConfigGenType: pb.ConfigGenType_CGT_FULL,
		Device: &pb.Device{
			Name:   "router1",
			Vendor: pb.Vendor_VD_CISCO,
			Model:  pb.Model_MD_UNKNOW,
		},
	}
	response, err := c.GetConfigGen(ctx, request)
	if err != nil {
		log.Fatal(err)
	}
	for _, config := range response.GetConfigFeature() {
		fmt.Println("results: ", config.Configuration)
		fmt.Println("Feature: ", config.Name)
	}
}
