# TCP RPC 示例

这是一个基于 Go 内置 RPC（远程过程调用）库使用 TCP 传输的简单示例。该示例包含服务端和客户端两个组件。

## 工作原理

此实现使用 Go 标准库中的 `net/rpc` 包通过 TCP 提供基础的 RPC 服务：

1. 服务端定义了一个服务（`HelloService`）及其可远程调用的方法（`SayHello`）。
2. 服务端将该服务注册到 Go 的 RPC 注册表中，使其对客户端可用。
3. 服务端在 1234 端口监听 TCP 连接。
4. 当客户端连接时，服务端在新的 goroutine 中处理 RPC 调用。
5. 客户端可以向服务端发起同步或异步调用。

## 服务端组件

- `HelloService`：提供 `SayHello` 方法的服务。
- `Request`：包含要问候的名字的结构体。
- `Response`：包含问候消息的结构体。

## 如何运行

### 启动服务端

```bash
cd server
go run server.go
```

服务端将启动并在 1234 端口上监听。

### 运行客户端

在另一个终端：

```bash
cd client
go run client.go
```

客户端将连接到服务端并进行同步和异步的 RPC 调用。

## 输出示例

### 服务端
```
RPC server started, listening on port 1234...
```

### 客户端
```
Connected to RPC server successfully
Synchronous call succeeded, response: Hello, Sync World
Performing other tasks while waiting for async response...
Asynchronous call succeeded, response: Hello, Async World
```

## 主要特点

- 使用原生 TCP 作为传输层
- 支持同步和异步 RPC 调用
- 使用 Go 内置的 RPC 序列化（gob）
- 通过信号处理实现优雅关闭

## 适用场景

这种基本的 TCP RPC 实现适用于：
- 可信网络中的内部服务通信
- 学习 RPC 概念
- 不需要 HTTP 开销的简单微服务

请注意，这个实现不包括生产环境中必要的身份验证、加密或复杂的错误处理机制。