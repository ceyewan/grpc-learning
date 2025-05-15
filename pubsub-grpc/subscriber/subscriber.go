package main

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "pubsub/proto/pubsub"
)

func subscribeToTopic(client pb.PubSubClient, topic string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Subscribing to topic: %s", topic)
	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{
		Topic: topic,
	})
	if err != nil {
		log.Printf("Could not subscribe to topic %s: %v", topic, err)
		return
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message from topic %s: %v", topic, err)
			break
		}
		log.Printf("Received message from topic %s: %s", topic, msg.Message)
	}
}

func main() {
	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewPubSubClient(conn)
	var wg sync.WaitGroup
	topics := []string{"topic1", "topic2", "topic3"}
	// 创建三个协程分别订阅不同的主题
	for _, topic := range topics {
		wg.Add(1)
		go subscribeToTopic(client, topic, &wg)
	}
	// 等待所有协程完成（实际上不会结束，除非出错）
	wg.Wait()
}
