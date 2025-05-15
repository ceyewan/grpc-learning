package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	pb "gochat/proto/gochat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	username = flag.String("name", "", "Username for chat")
)

func main() {
	flag.Parse()

	// 检查用户名，如果命令行没有提供，则提示输入
	if *username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your username: ")
		name, _ := reader.ReadString('\n')
		*username = strings.TrimSpace(name)
		if *username == "" {
			// 如果用户没有输入，使用默认用户名
			*username = fmt.Sprintf("User%d", time.Now().Unix()%1000)
			fmt.Printf("Using default username: %s\n", *username)
		}
	}

	// 设置连接到服务器
	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	// 创建聊天流
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Failed to start chat: %v", err)
	}

	// 设置信号处理，优雅退出
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup

	// 发送首条消息，用于注册用户身份
	firstMsg := &pb.ChatMessage{
		Username: *username,
		Content:  "has joined the chat", // 可选的初始消息
	}
	if err := stream.Send(firstMsg); err != nil {
		log.Fatalf("Failed to send first message: %v", err)
	}

	// 开始接收消息的 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Chat ended by server")
				cancel()
				return
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				cancel()
				return
			}
			// 显示接收到的消息
			if msg.Targetname == "" {
				// 广播消息
				fmt.Printf("[Broadcast] %s: %s\n", msg.Username, msg.Content)
			} else {
				// 私聊消息
				fmt.Printf("[Private from %s]: %s\n", msg.Username, msg.Content)
			}
		}
	}()

	// 开始发送消息的 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("=== Chat started ===")
		fmt.Println("- To send a broadcast message: just type your message")
		fmt.Println("- To send a private message: @username message")
		fmt.Println("- To exit: type 'exit'")

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if scanner.Scan() {
					input := scanner.Text()
					if input == "exit" {
						fmt.Println("Exiting chat...")
						cancel()
						return
					}

					var targetName string
					var content string

					// 检查是否是私聊消息（以@开头）
					if strings.HasPrefix(input, "@") {
						parts := strings.SplitN(input, " ", 2)
						if len(parts) >= 2 {
							targetName = strings.TrimPrefix(parts[0], "@")
							content = parts[1]
						} else {
							fmt.Println("Invalid format for private message. Use: @username message")
							continue
						}
					} else {
						// 广播消息
						content = input
					}

					msg := &pb.ChatMessage{
						Username:   *username,
						Targetname: targetName,
						Content:    content,
					}

					if err := stream.Send(msg); err != nil {
						log.Printf("Failed to send message: %v", err)
						cancel()
						return
					}
				} else {
					if scanner.Err() != nil {
						log.Printf("Error reading input: %v", scanner.Err())
					}
					cancel()
					return
				}
			}
		}
	}()

	// 等待信号或上下文取消
	select {
	case <-signalChan:
		fmt.Println("\nReceived interrupt signal, exiting...")
		cancel()
	case <-ctx.Done():
		// 其他goroutine已经取消了上下文
	}

	// 等待goroutine完成
	wg.Wait()
	fmt.Println("Chat client closed")
}
