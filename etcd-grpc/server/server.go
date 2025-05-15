package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"helloworld/etcd"
	pb "helloworld/proto/helloworld"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
	port string
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + req.Name + " from port " + s.port}, nil
}

// 匿名类型转换，确保 server 实现了 GreeterServer 接口
var _ pb.GreeterServer = (*server)(nil)

func main() {
	port := flag.String("port", "1234", "服务端口")
	flag.Parse()
	addr := ":" + *port

	if err := etcd.InitDefaultClient(nil); err != nil {
		log.Fatalf("etcd 初始化失败: %v", err)
	}
	defer etcd.CloseDefaultClient()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	// 注册到 etcd
	registry, err := etcd.NewServiceRegistry(nil)
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}

	// 注册服务
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = registry.Register(ctx, "greater-service", "instance-"+*port, "localhost:"+*port)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{port: *port})
	log.Printf("gRPC 服务器启动于 %s...", addr)

	// 优雅退出，注销 etcd
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("收到退出信号，关闭服务器...")
		listener.Close()
		os.Exit(0)
	}()

	if err := s.Serve(listener); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
