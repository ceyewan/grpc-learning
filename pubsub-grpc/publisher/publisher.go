package main

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "pubsub/proto/pubsub"
)

// 生成随机消息
func generateRandomMessage(r *rand.Rand) string {
	words := []string{
		"你好", "世界", "消息", "发布", "订阅", "gRPC", "Golang",
		"微服务", "分布式", "通信", "实时", "数据", "流", "推送",
		"事件", "系统", "云原生", "应用", "接口", "协议",
	}

	var sb strings.Builder
	wordCount := 3 + r.Intn(5) // 随机生成单词数量

	for i := 0; i < wordCount; i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(words[r.Intn(len(words))])
	}

	return sb.String()
}

func publishMessage(client pb.PubSubClient, topic, message string) {
	resp, err := client.Publish(context.Background(), &pb.PublishRequest{
		Topic:   topic,
		Message: message,
	})
	if err != nil {
		log.Printf("发布到主题 %s 失败: %v", topic, err)
		return
	}
	log.Printf("发布到主题 %s 成功: %v, 消息内容: %s", topic, resp.Success, message)
}

func main() {
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接: %v", err)
	}
	defer conn.Close()
	client := pb.NewPubSubClient(conn)
	topics := []string{"topic1", "topic2", "topic3"}
	log.Println("开始轮流发布消息到不同主题...")
	// 无限循环，每隔一段时间发送一条消息
	for i := 0; ; i++ {
		// 选择主题 (轮询方式)
		topic := topics[i%len(topics)]
		// 生成随机消息
		message := generateRandomMessage(r)
		// 发布消息
		publishMessage(client, topic, message)
		// 等待1-3秒
		sleepTime := 1000 + rand.Intn(2000)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}
