package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

type HelloService struct{}

type Request struct {
	Name string
}

type Response struct {
	Message string
}

func (h *HelloService) SayHello(req Request, resp *Response) error {
	if req.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	resp.Message = fmt.Sprintf("Hello, %s", req.Name)
	return nil
}

func main() {
	if err := rpc.Register(new(HelloService)); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	// 监听端口，监听的是 TCP 协议 1234 端口
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to listen on port 1234: %v", err)
		return
	}
	defer listener.Close()

	// 优雅退出
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	log.Println("RPC server started, listening on port 1234...")
	go func() {
		for {
			// 等待客户端连接，连接成功后返回一个 net.Conn 对象
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			// 调用 rpc.ServeConn 处理连接
			go rpc.ServeConn(conn) // 处理 RPC 连接
		}
	}()

	<-shutdown
	log.Println("Server shutting down...")
}
