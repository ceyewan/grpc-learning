# GoChat: gRPC-based Chat Application

## 项目简介

GoChat 是一个基于 gRPC 的实时聊天应用程序，采用 Go 语言开发。该应用程序支持多用户聊天、消息广播和私人消息功能，展示了如何使用 gRPC 双向流 (bidirectional streaming) 实现实时通信应用。

## 功能特点

- 基于 gRPC 的双向流实时通信
- 支持多客户端同时连接
- 支持广播消息（发送给所有在线用户）
- 支持私聊消息（发送给特定用户）
- 自动用户名注册与注销
- 简洁的命令行界面

## 系统架构

该项目采用客户端-服务器架构：

1. **服务器 (Server)**：
   - 管理客户端连接
   - 处理用户注册和注销
   - 转发消息到目标客户端

2. **客户端 (Client)**：
   - 连接到服务器
   - 发送和接收消息
   - 提供用户界面以进行聊天操作

## 技术栈

- Go 语言 (1.24.x)
- gRPC 框架
- Protocol Buffers (protobuf) 作为接口定义语言
- 标准库的网络和并发原语

## 项目结构

```
gochat-grpc/
├── client/           # 客户端代码
│   └── client.go     # 客户端实现
├── server/           # 服务器代码
│   └── server.go     # 服务器实现
├── proto/            # 生成的 Protocol Buffers 代码
│   └── gochat/       # 生成的 Go 代码包
│       ├── chat.pb.go
│       └── chat_grpc.pb.go
├── gochat.proto      # Protocol Buffers 服务定义
├── go.mod            # Go 模块定义
├── go.sum            # 依赖版本锁定文件
└── README.md         # 项目文档（本文件）
```

## 安装与运行

### 前置条件

- 安装 Go (版本 1.24.x 或更高)
- 安装 Protocol Buffers 编译器 (protoc)
- 安装 Go 的 Protocol Buffers 插件

### 编译

1. 克隆仓库：

```bash
git clone <仓库地址>
cd gochat-grpc
```

2. 生成 Protocol Buffers 代码（如果需要重新生成）：

```bash
protoc --go_out=. --go-grpc_out=. gochat.proto
```

3. 编译服务器和客户端：

```bash
go build -o bin/server ./server
go build -o bin/client ./client
```

### 运行

1. 首先启动服务器：

```bash
./bin/server
```

服务器将在 1234 端口启动并等待客户端连接。

2. 启动客户端并指定用户名（可选）：

```bash
./bin/client -name YourUsername
```

如果不指定用户名，客户端会提示您输入一个用户名。

## 使用说明

### 客户端命令

启动客户端后，可以使用以下命令：

1. **发送广播消息**：
   直接输入消息然后按回车键，将发送给所有在线用户。

2. **发送私聊消息**：
   使用格式 `@用户名 消息内容` 可以向特定用户发送私聊消息。
   例如：`@Alice 你好！` 会将消息 "你好！" 只发送给用户 "Alice"。

3. **退出**：
   输入 `exit` 或按 Ctrl+C 可以退出聊天程序。

## 工作原理

### gRPC 双向流通信

GoChat 使用 gRPC 的双向流 (bidirectional streaming) RPC 功能实现实时通信：

1. 客户端和服务器之间建立长连接。
2. 客户端可以持续向服务器发送消息，服务器也可以持续向客户端推送消息。
3. 消息以 Protocol Buffers 格式序列化，高效传输。

### 消息路由

当服务器接收到消息时：

- 如果目标用户名为空，服务器将消息广播给所有其他连接的客户端。
- 如果指定了目标用户名，服务器仅将消息发送给该特定用户。

### 用户管理

- 用户在连接时自动注册（发送第一条消息时包含用户名）。
- 用户断开连接时自动注销。
- 服务器维护一个映射表，关联用户名和对应的 gRPC 流。
