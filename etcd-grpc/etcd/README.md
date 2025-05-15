# etcd-grpc 服务注册与发现组件

etcd-grpc 是一个基于 etcd 构建的服务注册与发现组件，专为 gRPC 微服务架构设计。它提供了简单易用的 API 接口，使微服务能够自动注册到 etcd，并能够发现和连接其他服务。

---

## 组件核心功能

### 服务注册
- 将服务信息注册到 etcd
- 自动续约租约保持服务可用性
- 支持服务优雅下线

### 服务发现
- 基于 etcd 的实时服务发现
- 支持 gRPC 原生服务解析
- 内置负载均衡支持

### 高可靠性
- 支持 etcd 集群配置
- 错误处理和自动恢复
- 详细日志记录

---

## 架构设计

该组件主要包含以下核心模块：

- **配置管理 (Config)**
  - 提供可配置的 etcd 连接参数
  - 支持默认配置和自定义配置
- **客户端管理 (Client)**
  - 封装 etcd 客户端的创建和管理
  - 提供单例模式的默认客户端
- **服务注册 (ServiceRegistry)**
  - 处理服务注册和注销
  - 基于租约机制实现自动续约
- **服务发现 (ServiceDiscovery)**
  - 提供服务查询和连接建立
  - 返回支持负载均衡的 gRPC 连接
- **解析器 (Resolver)**
  - 实现 gRPC 解析器接口
  - 监听服务变化并更新连接状态

---

## 使用示例

### 初始化配置

```go
// 使用默认配置
err := etcd.InitDefaultClient(nil)

// 或使用自定义配置
config := &etcd.Config{
    Endpoints:   []string{"localhost:2379"},
    DialTimeout: 5 * time.Second,
}
err := etcd.InitDefaultClient(config)
```

### 服务注册

```go
registry, err := etcd.NewServiceRegistry(nil)
if err != nil {
    log.Fatalf("Failed to create registry: %v", err)
}

// 注册服务
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err = registry.Register(ctx, "user-service", "instance-1", "localhost:50051")
if err != nil {
    log.Fatalf("Failed to register service: %v", err)
}

// 程序结束时注销服务
defer registry.Deregister(context.Background(), "user-service", "instance-1")
```

### 服务发现和调用

```go
discovery, err := etcd.NewServiceDiscovery(nil)
if err != nil {
    log.Fatalf("Failed to create discovery: %v", err)
}

// 获取服务连接
conn, err := discovery.GetConnection(context.Background(), "user-service")
if err != nil {
    log.Fatalf("Failed to connect to service: %v", err)
}
defer conn.Close()

// 创建客户端并调用服务
client := pb.NewUserServiceClient(conn)
response, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: "123"})
```

---

## 最佳实践

- **优雅启动和关闭**
  - 确保在程序启动时初始化 etcd 客户端
  - 程序退出前释放资源并注销服务
- **错误处理**
  - 妥善处理连接和注册错误
  - 实现重试机制提高系统稳定性
- **监控和日志**
  - 关注服务注册和发现的日志
  - 考虑添加指标收集以监控服务健康状况
- **配置管理**
  - 将 etcd 地址等配置从环境变量或配置文件加载
  - 避免在代码中硬编码配置值

---

## 扩展点

- **认证支持**
  - 添加 TLS/SSL 支持以保护通信安全
  - 实现基于证书的身份验证
- **服务元数据**
  - 扩展服务注册以支持更多元数据
  - 支持基于标签或版本的服务筛选
- **高级负载均衡**
  - 实现基于权重的负载均衡
  - 支持自定义负载均衡策略

该组件为构建可靠、可扩展的 gRPC 微服务提供了坚实的基础，通过简单的 API 接口隐藏了服务注册和发现的复杂性。

---

## gRPC ClientConn 详解

### 什么是 ClientConn？

`grpc.ClientConn` 是 gRPC 客户端的核心连接对象。它并不是一个简单的 TCP 连接，而是一个连接管理器，负责：

- 维护与多个服务实例的 TCP 连接（连接池）
- 负载均衡：自动将请求分发到不同的服务实例
- 服务发现：动态感知 etcd 注册中心的服务上下线
- 连接状态管理与自动恢复

### 连接建立的时机

1. **初始化阶段**
   - 当你调用 `grpc.Dial` 或 `grpc.DialContext`（如 `GetConnection` 返回的 conn）时，gRPC 会创建一个 `ClientConn` 对象，并初始化服务发现解析器（如 etcd 解析器），获取所有可用服务实例的地址列表。
   - 此时**不会立即与所有服务实例建立 TCP 连接**。

2. **首次 RPC 调用时**
   - 当你第一次通过 `ClientConn` 发起 RPC 调用时，gRPC 会根据负载均衡策略（如 round_robin）选择一个服务实例，如果还没有与该实例建立连接，则会**按需建立 TCP 连接**。
   - 后续请求会复用已建立的连接，或在需要时为新实例建立新连接。

3. **服务变更时**
   - etcd 解析器监听 etcd 注册中心的服务上下线变化，动态更新可用服务实例列表。
   - 新增实例时，只有在有请求分发到该实例时才会建立连接。
   - 下线实例时，gRPC 会关闭对应的连接。

### 为什么每次调用可能访问不同的服务实例？

- 由于 gRPC 的负载均衡机制（如 round_robin），每次通过同一个 `ClientConn` 发起的 RPC 调用，可能会被分发到不同的服务实例。
- 这样可以实现请求的自动分流和高可用，无需客户端关心服务实例的具体地址和数量。

### 总结

- `ClientConn` 是 gRPC 客户端的连接管理器，支持多实例、负载均衡和服务发现。
- 连接是在**首次调用时按需建立**，不是在服务发现时全部建立。
- 负载均衡策略决定了每次请求实际访问的服务实例。