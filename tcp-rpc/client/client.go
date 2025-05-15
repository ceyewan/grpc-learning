package main

import (
	"log"
	"net/rpc"
	"time"
)

type Request struct {
	Name string
}

type Response struct {
	Message string
}

func main() {
	// 连接 RPC 服务器
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Failed to connect to RPC server: %v", err)
	}
	defer client.Close()
	log.Println("Connected to RPC server successfully")

	// 同步调用 RPC 方法
	syncRequest := Request{Name: "Sync World"}
	var syncResponse Response
	err = client.Call("HelloService.SayHello", syncRequest, &syncResponse)
	if err != nil {
		log.Fatalf("Synchronous RPC call failed: %v", err)
	}
	log.Printf("Synchronous call succeeded, response: %s", syncResponse.Message)

	// 异步调用 RPC 方法
	asyncRequest := Request{Name: "Async World"}
	var asyncResponse Response

	// 发起异步调用
	call := client.Go("HelloService.SayHello", asyncRequest, &asyncResponse, nil)

	// 执行其他任务，用 time.Sleep 模拟
	log.Println("Performing other tasks while waiting for async response...")
	time.Sleep(time.Second)

	// 等待 RPC 结果返回
	replyCall := <-call.Done
	if replyCall.Error != nil {
		log.Fatalf("Asynchronous RPC call failed: %v", replyCall.Error)
	}
	log.Printf("Asynchronous call succeeded, response: %s", asyncResponse.Message)
}
