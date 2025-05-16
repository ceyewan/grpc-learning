package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "pubsub/proto/pubsub"
)

type pubSubServer struct {
	pb.UnimplementedPubSubServer
	redisClient *redis.Client
}

var _ pb.PubSubServer = (*pubSubServer)(nil)

func NewPubSubServer() *pubSubServer {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	return &pubSubServer{
		redisClient: redisClient,
	}
}

func (s *pubSubServer) Publish(stream pb.PubSub_PublishServer) error {
	ctx := stream.Context()
	messageCount := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 客户端已发送完所有消息
			return stream.SendAndClose(&pb.PublishResponse{
				Success:      true,
				MessageCount: int32(messageCount),
			})
		}
		if err != nil {
			return status.Errorf(codes.Internal, "接收消息失败: %v", err)
		}

		// 发布消息到Redis
		err = s.redisClient.Publish(ctx, req.Topic, req.Message).Err()
		if err != nil {
			return status.Errorf(codes.Internal, "发布消息失败: %v", err)
		}
		
		messageCount++
		log.Printf("已发布消息到主题 %s: %s", req.Topic, req.Message)
	}
}

func (s *pubSubServer) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	pubsub := s.redisClient.Subscribe(stream.Context(), req.Topic)
	defer pubsub.Close()

	log.Printf("客户端已订阅主题: %s", req.Topic)

	ch := pubsub.Channel()
	for msg := range ch {
		err := stream.Send(&pb.SubscribeResponse{Message: msg.Payload})
		if err != nil {
			return status.Errorf(codes.Internal, "发送消息失败: %v", err)
		}
		log.Printf("已发送消息到订阅者 (主题 %s): %s", req.Topic, msg.Payload)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterPubSubServer(server, NewPubSubServer())

	log.Println("服务器运行在端口 1234")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("服务运行失败: %v", err)
	}
}