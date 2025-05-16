package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "pubsub/proto/pubsub"
)

// 生成有意义的消息
func generateMeaningfulMessage(topic string, index int) string {
	switch topic {
	case "topic1":
		return fmt.Sprintf("技术新闻 #%d: Go 1.22 发布了新的性能改进", index)
	case "topic2":
		return fmt.Sprintf("市场动态 #%d: 云计算服务需求持续增长", index)
	case "topic3":
		return fmt.Sprintf("系统通知 #%d: 服务器将在下周进行例行维护", index)
	default:
		return fmt.Sprintf("消息 #%d 发送到 %s", index, topic)
	}
}

func publishMessages(client pb.PubSubClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号监听，以便优雅地关闭流
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("接收到终止信号，正在关闭...")
		cancel()
	}()

	topics := []string{"topic1", "topic2", "topic3"}
	var wg sync.WaitGroup

	// 为每个主题创建一个协程发布消息
	for _, topic := range topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()

			// 创建流
			stream, err := client.Publish(ctx)
			if err != nil {
				log.Fatalf("创建发布流失败: %v", err)
			}

			log.Printf("开始向主题 %s 发布消息...", topic)

			// 发送10条消息到指定主题
			for i := 1; i <= 10; i++ {
				select {
				case <-ctx.Done():
					// 上下文被取消，结束发送
					return
				default:
					// 生成有意义的消息
					message := generateMeaningfulMessage(topic, i)

					// 发送消息
					req := &pb.PublishRequest{
						Topic:   topic,
						Message: message,
					}
					if err := stream.Send(req); err != nil {
						log.Printf("发送消息失败: %v", err)
						return
					}

					log.Printf("已发送消息到主题 %s: %s", topic, message)

					// 短暂等待
					time.Sleep(500 * time.Millisecond)
				}
			}

			// 完成发送，关闭流
			resp, err := stream.CloseAndRecv()
			if err != nil {
				log.Printf("关闭流时出错: %v", err)
			} else {
				log.Printf("向主题 %s 发布结束。成功: %v, 消息数: %d",
					topic, resp.Success, resp.MessageCount)
			}
		}(topic)
	}

	// 等待所有发布完成
	wg.Wait()
	log.Println("所有主题的消息发布完成")
}

func main() {
	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接: %v", err)
	}
	defer conn.Close()

	client := pb.NewPubSubClient(conn)
	publishMessages(client)
}
