package main

import (
	"context"
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

func (s *pubSubServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	err := s.redisClient.Publish(ctx, req.Topic, req.Message).Err()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}
	return &pb.PublishResponse{Success: true}, nil
}

func (s *pubSubServer) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	pubsub := s.redisClient.Subscribe(stream.Context(), req.Topic)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		err := stream.Send(&pb.SubscribeResponse{Message: msg.Payload})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send message: %v", err)
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterPubSubServer(server, NewPubSubServer())

	log.Println("Server is running on port 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
