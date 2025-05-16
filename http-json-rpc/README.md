# HTTP JSON-RPC 示例

这是一个基于 HTTP 传输的 JSON-RPC 服务实现示例。该示例包含 Go 服务端和 Python 客户端，展示了跨语言互操作性。

## 工作原理

此实现结合了 Go 的 `net/rpc/jsonrpc` 包和标准 HTTP 服务器，通过 HTTP 提供 RPC 服务：

1. 服务端定义了一个服务（`HelloService`）及其可远程调用的方法（`SayHello`）。
2. 服务端将该服务注册到 Go 的 RPC 注册表中。
3. 服务端在 `/jsonrpc` 端点创建 HTTP 处理器。
4. 当 HTTP POST 请求到达此端点时，服务端将其处理为 JSON-RPC 调用。
5. 客户端发送 JSON 格式的 POST 请求来调用远程过程。

这种方法的主要优势是它通过标准 HTTP 工作，使其可以跨不同平台和语言访问。

## 服务端组件

- `HelloService`：提供 `SayHello` 方法的服务。
- `Request`：包含要问候的名字的结构体。
- `Response`：包含问候消息的结构体。
- HTTP 处理器：使用自定义 `io.ReadWriteCloser` 适配器将 HTTP 请求转换为 RPC 调用。

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
python client.py
```

Python 客户端将连接到服务端并进行 RPC 调用。

## 输出示例

### 服务端
```
HTTP JSON-RPC server started, listening on port 1234...
```

### 客户端
```
Message: Hello, Python HTTP Client
```

## 主要特点

- 使用 HTTP 作为传输层
- 使用 JSON 进行序列化
- 跨语言支持（Go 服务端，Python 客户端）
- 简单的 HTTP 端点（`/jsonrpc`）
- 通过信号处理实现优雅关闭

## 适用场景

这种 HTTP JSON-RPC 实现适用于：
- 需要 RPC 风格通信的 Web 应用
- 需要被多语言客户端访问的服务
- 需要简单 JSON 调用的公共 API
- HTTP 是通用协议的遗留系统集成

## 客户端实现细节

Python 客户端演示了如何构造 JSON-RPC 请求：
1. 准备包含 `method`、`params` 和 `id` 字段的 JSON 载荷
2. 设置正确的 Content-Type 头
3. 向服务器的 JSON-RPC 端点发送 POST 请求
4. 解析 JSON 响应以提取结果

这展示了使用 HTTP 上的 JSON 进行 RPC 通信的互操作性优势。