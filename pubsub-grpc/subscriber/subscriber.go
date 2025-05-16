package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "pubsub/proto/pubsub"
)

func subscribeToTopic(client pb.PubSubClient, topic string, ctx context.Context) error {
	log.Printf("开始订阅主题: %s", topic)

	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{
		Topic: topic,
	})
	if err != nil {
		return fmt.Errorf("无法订阅主题 %s: %v", topic, err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("取消订阅主题 %s", topic)
			return nil
		default:
			msg, err := stream.Recv()
			if err != nil {
				return fmt.Errorf("从主题 %s 接收消息时出错: %v", topic, err)
			}
			log.Printf("从主题 %s 接收到消息: %s", topic, msg.Message)
		}
	}
}

func main() {
	// 定义命令行参数
	var topic string
	flag.StringVar(&topic, "topic", "", "要订阅的主题名称 (必需)")
	flag.Parse()

	// 检查是否提供了主题参数
	if topic == "" {
		fmt.Println("错误: 必须指定要订阅的主题")
		fmt.Println("用法: subscriber -topic=<主题名称>")
		fmt.Println("示例: subscriber -topic=topic1")
		os.Exit(1)
	}

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听终止信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("接收到终止信号，正在关闭...")
		cancel()
	}()

	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewPubSubClient(conn)

	// 订阅指定的主题
	if err := subscribeToTopic(client, topic, ctx); err != nil {
		log.Fatalf("订阅失败: %v", err)
	}

	log.Println("订阅已关闭")
}
