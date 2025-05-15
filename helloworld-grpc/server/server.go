package main

import (
	"context"
	"log"
	"net"

	pb "helloworld/proto/helloworld"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + req.Name}, nil
}

// 匿名类型转换，确保 server 实现了 GreeterServer 接口
var _ pb.GreeterServer = (*server)(nil)

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Println("gRPC 服务器启动...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
