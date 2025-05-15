package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"helloworld/etcd"
	pb "helloworld/proto/helloworld"
)

func main() {
	// 初始化etcd客户端
	if err := etcd.InitDefaultClient(nil); err != nil {
		log.Fatalf("etcd 初始化失败: %v", err)
	}
	defer etcd.CloseDefaultClient()

	// 创建服务发现实例
	discovery, err := etcd.NewServiceDiscovery(nil)
	if err != nil {
		log.Fatalf("Failed to create discovery: %v", err)
	}

	// 获取服务连接
	conn, err := discovery.GetConnection(context.Background(), "greater-service")
	if err != nil {
		log.Fatalf("Failed to connect to service: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	n := 20 // 并发请求数
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			resp, err := client.SayHello(c, &pb.HelloRequest{Name: fmt.Sprintf("gRPC Client %d", i)})
			if err != nil {
				log.Printf("调用失败: %v", err)
				return
			}
			fmt.Printf("[%d] 服务端返回: %s\n", i, resp.Message)
		}(i)
	}
	wg.Wait()
}
