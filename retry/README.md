# gRPC HelloWorld 示例


这是一个使用 gRPC 框架实现的简单 HelloWorld 服务示例。本示例展示了 gRPC 的基本用法，包括服务定义、服务端和客户端实现。

> **增强说明**：
>
> 本项目的客户端实现了自动重连（retry）机制，并结合了 `context.WithTimeout`，确保在网络异常或服务端暂时不可用时，客户端能够自动重试连接，并在每次请求时设置超时，提升了健壮性和用户体验。

## 工作原理

gRPC 是一个高性能、开源的通用 RPC 框架，由 Google 开发：

1. 使用 Protocol Buffers (protobuf) 作为接口定义语言 (IDL)，定义服务接口和消息结构。
2. 服务端实现在 `.proto` 文件中定义的服务接口。
3. 使用 protoc 编译器生成客户端和服务端代码。
4. 服务端注册并运行 gRPC 服务。
5. 客户端创建 stub（存根）连接服务端并进行调用。


## 项目组件

- `helloworld.proto`：定义服务接口和消息格式
- `proto/`：存放生成的 protobuf 代码
- `server/`：gRPC 服务端实现
- `client/`：gRPC 客户端实现（包含自动重连与超时控制）

## 如何运行

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 生成 protobuf 代码（如果需要重新生成）

```bash
protoc --go_out=. --go-grpc_out=. helloworld.proto
```

### 3. 启动服务端

建议设置 gRPC 日志级别为 warning，确保重试的调试信息输出：

```bash
export GRPC_GO_LOG_SEVERITY_LEVEL=warning
```

然后启动服务端：

```bash
cd server
go run server.go
```

服务端将启动并在 1234 端口上监听。


### 4. 运行客户端

在另一个终端，同样建议设置日志级别：

```bash
export GRPC_GO_LOG_SEVERITY_LEVEL=warning
```

然后运行客户端：

```bash
cd client
go run client.go
```

客户端将自动尝试连接服务器并发送请求。
如果连接失败，客户端会自动重试，直到连接成功或超时。

## 输出示例

### 服务端
```
gRPC 服务器启动...
```

### 客户端
```
服务端返回: Hello, gRPC Client
```


## 主要特点

- 使用 Protocol Buffers 定义服务和消息
- 强类型的服务接口和消息
- 高效的二进制序列化
- 支持流式 RPC（本示例未演示）
- 跨语言支持
- 内置连接管理和安全机制
- **客户端支持自动重连与超时控制**：通过 gRPC 的重试机制和 `context.WithTimeout`，提升了客户端的健壮性和容错能力。

## 适用场景

gRPC 特别适合以下场景：

- 微服务架构中的服务间通信
- 需要高效通信的分布式系统
- 需要强类型 API 的系统
- 移动应用与后端服务通信
- 需要多语言支持的系统

## Proto 文件说明

`helloworld.proto` 文件定义了：

- `Greeter` 服务：包含 `SayHello` RPC 方法
- `HelloRequest` 消息：包含请求参数 `name`
- `HelloReply` 消息：包含响应消息 `message`

这是整个 gRPC 通信的基础，服务端和客户端代码都基于这个定义生成。