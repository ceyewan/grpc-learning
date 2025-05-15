package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "helloworld/proto/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "gRPC Client"})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}
	fmt.Println("服务端返回:", resp.Message)
}
