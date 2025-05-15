package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
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
	// 注册 RPC 服务
	if err := rpc.Register(new(HelloService)); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// 在 /jsonrpc 路径上配置一个 HTTP 处理函数
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		// io.ReadWriteCloser 是一个接口，表示既可以读取又可以写入的对象
		// 将 HTTP 请求体（r.Body）作为读取端，将 HTTP 响应写入端（w）作为写入端
		// 使用匿名结构体封装 ReadCloser 和 Writer，组合成一个 ReadWriteCloser
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		// 基于 conn 创建一个 JSON 编解码器，并处理 RPC 请求
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	// 设置优雅退出
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务器
	log.Println("HTTP JSON-RPC server started, listening on port 1234...")
	go http.ListenAndServe(":1234", nil)

	// 等待退出信号
	<-shutdown
	log.Println("Server shutting down...")
}
