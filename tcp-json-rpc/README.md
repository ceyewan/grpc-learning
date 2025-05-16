# TCP JSON-RPC 示例

这是一个使用 JSON 序列化通过 TCP 传输的 Go RPC（远程过程调用）实现示例。该示例包含服务端和客户端两个组件。

## 工作原理

此实现使用 Go 标准库中的 `net/rpc/jsonrpc` 包通过 TCP 提供 JSON 编码的 RPC 服务：

1. 服务端定义了一个服务（`HelloService`）及其可远程调用的方法（`SayHello`）。
2. 服务端将该服务注册到 Go 的 RPC 注册表中，使其对客户端可用。
3. 服务端在 1234 端口监听 TCP 连接。
4. 当客户端连接时，服务端使用 JSON 编码处理 RPC 调用。
5. 客户端可以向服务端发起同步或异步调用。

与标准 TCP RPC 的主要区别在于使用 JSON 序列化替代了 Go 默认的 gob 编码。

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

- 使用 TCP 作为传输层
- 使用 JSON 进行序列化（替代 Go 默认的 gob 编码）
- 支持同步和异步 RPC 调用
- 通过信号处理实现优雅关闭
- 可与能生成 JSON RPC 调用的非 Go 客户端互操作

## 适用场景

这种 TCP JSON-RPC 实现适用于：
- 需要跨语言边界通信的服务
- 需要在传输中保持数据人类可读性的系统
- 调试 RPC 调用（因为 JSON 是人类可读的）
- 需要与非 Go 客户端互操作的简单微服务

请注意，虽然此实现比标准 TCP RPC 更具互操作性，但仍然不包括生产环境中必要的身份验证、加密或复杂的错误处理机制。