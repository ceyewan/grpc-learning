# gRPC PubSub 系统

这是一个基于 gRPC 和 Redis 构建的发布/订阅系统学习项目。本项目实现了客户端流式的发布者和服务端流式的订阅者，用于展示 gRPC 流式通信的特点。

## 项目特性

- **客户端流式发布**：使用 gRPC 客户端流式 API 实现消息发布
- **服务端流式订阅**：使用 gRPC 服务端流式 API 实现消息订阅
- **Redis 中间件**：使用 Redis 作为消息存储和分发的中间件
- **有意义的消息**：发布者向不同主题发布有类型化的消息
- **灵活订阅**：订阅者可以通过命令行选择订阅的主题

## 快速开始

### 1. 启动 Redis Docker 服务

```bash
docker-compose up -d
```

这将启动 Redis 服务器

### 2. 生成 gRPC 代码

如果修改了 proto 文件，需要重新生成代码：

```bash
protoc --go_out=. --go-grpc_out=. pubsub.proto
```

### 3. 启动服务器

在一个终端中运行：

```bash
go run server/server.go
```

### 4. 启动订阅者

在另一个终端中运行（可以指定要订阅的主题）：

```bash
go run subscriber/subscriber.go -topic=topic1
```

可以在多个终端中订阅不同主题：`topic1`、`topic2` 或 `topic3`

### 5. 启动发布者

在第三个终端中运行：

```bash
go run publisher/publisher.go
```

发布者将向三个主题各发送 10 条有意义的消息，然后自动关闭。

### 6. 停止 Redis 服务

完成后，可以停止 Redis 服务：

```bash
docker-compose down
```

## 项目结构

```
pubsub-grpc/
├── docker-compose.yml    # Redis Docker 配置
├── proto/                # 生成的 gRPC 代码
├── pubsub.proto          # 协议定义
├── publisher/            # 发布者客户端
├── server/               # gRPC 服务器
└── subscriber/           # 订阅者客户端
```

## 工作原理

1. **服务器**：启动 gRPC 服务器，连接到 Redis 并提供发布/订阅接口
2. **发布者**：使用客户端流式接口，向三个主题（topic1、topic2、topic3）各发送 10 条有意义的消息
3. **订阅者**：使用服务端流式接口，订阅特定主题并实时接收消息

## 消息主题

本项目中的三个主题发布不同类型的消息：
- `topic1`: 技术新闻相关消息
- `topic2`: 市场动态相关消息
- `topic3`: 系统通知相关消息

## 前置要求

- Go 1.16+
- Protocol Buffers 编译器 (protoc)
- Docker 和 Docker Compose（用于运行 Redis）

## API 定义

在 `pubsub.proto` 文件中定义了两个核心 RPC 方法：

```proto
// 发布消息 - 客户端流式
rpc Publish (stream PublishRequest) returns (PublishResponse);

// 订阅消息 - 服务端流式
rpc Subscribe (SubscribeRequest) returns (stream SubscribeResponse);
```

## 开发说明

如果需要修改项目：

1. 编辑 `pubsub.proto` 文件，定义新的消息或服务
2. 运行 `protoc` 命令重新生成 gRPC 代码
3. 修改服务器、发布者或订阅者的实现代码
4. 通过上述命令测试您的更改

## 查看 Docker 日志

如果需要查看 Redis 的日志：

```bash
docker-compose logs
```

## 问题排查

- 如果无法连接到 Redis，请确保 Redis 服务已启动：`docker-compose up -d`
- 如果遇到 gRPC 相关错误，尝试重新生成代码：`protoc --go_out=. --go-grpc_out=. pubsub.proto`
- 如果发布者或订阅者报错，确保服务器正在运行

## 学习资源

- [gRPC 官方文档](https://grpc.io/docs/)
- [Protocol Buffers 文档](https://developers.google.com/protocol-buffers)
- [Go Redis 客户端文档](https://redis.uptrace.dev/guide/)
