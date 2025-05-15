package main

import (
	"io"
	"log"
	"net"
	"sync"

	pb "gochat/proto/gochat"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[string]pb.ChatService_ChatServer // 用户名到流的映射
}

func NewServer() *server {
	return &server{
		clients: make(map[string]pb.ChatService_ChatServer),
	}
}

// Chat 实现
func (s *server) Chat(stream pb.ChatService_ChatServer) error {
	// 第一条消息用于识别用户和关联流
	firstMsg, err := stream.Recv()
	if err != nil {
		log.Printf("Error receiving first message: %v", err)
		return err
	}

	username := firstMsg.Username
	// 用户连接 - 自动注册
	s.mu.Lock()
	s.clients[username] = stream
	s.mu.Unlock()
	log.Printf("%s joined the chat", username)

	// 处理首条消息（如果有实际内容）
	if firstMsg.Content != "" {
		s.handleMessage(firstMsg, username)
	}

	// 处理后续消息
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}
		s.handleMessage(msg, username)
	}

	// 用户断开连接 - 自动注销
	s.mu.Lock()
	delete(s.clients, username)
	s.mu.Unlock()
	log.Printf("%s left the chat", username)
	return nil
}

// 处理消息发送
func (s *server) handleMessage(msg *pb.ChatMessage, senderUsername string) {
	if msg.Targetname == "" {
		// 广播
		s.mu.Lock()
		for uname, clientStream := range s.clients {
			if uname != senderUsername {
				_ = clientStream.Send(msg)
			}
		}
		s.mu.Unlock()
	} else {
		// 单发
		s.mu.Lock()
		targetStream, ok := s.clients[msg.Targetname]
		s.mu.Unlock()
		if ok {
			_ = targetStream.Send(msg)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, NewServer())
	log.Println("gRPC Chat Server started on port 1234")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
