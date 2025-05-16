# etcd-gRPC 服务发现与负载均衡

这个项目演示了如何使用 etcd 实现 gRPC 服务的服务注册、服务发现和负载均衡功能。

## 项目概述

本项目实现了一个简单的 gRPC 服务，并使用 etcd 作为服务注册中心。项目包含以下核心功能：

- 基于 etcd 的服务注册
- 基于 etcd 的服务发现
- gRPC 客户端负载均衡
- 多实例服务部署

## 系统架构

```
┌─────────────┐          ┌─────────────┐
│   Client    │◄─────────┤  etcd集群   │
└─────┬───────┘          └─────▲───────┘
      │                        │
      │                        │ 注册服务
      │ 服务发现               │
      │ 负载均衡               │
      ▼                        │
┌─────────────┐          ┌─────┴───────┐
│   Service   ├──────────► Service实例1 │
└─────────────┘          ├─────────────┤
                         │ Service实例2 │
                         ├─────────────┤
                         │ Service实例3 │
                         └─────────────┘
```

## 技术栈

- Go 1.24+
- gRPC
- etcd (服务注册与发现)
- Protocol Buffers
- Docker & Docker Compose

## 快速开始

### 前置条件

- 安装 [Go 1.24+](https://golang.org/dl/)
- 安装 [Docker](https://docs.docker.com/get-docker/) 和 [Docker Compose](https://docs.docker.com/compose/install/)
- 安装 [Protocol Buffers 编译器](https://github.com/protocolbuffers/protobuf/releases)
- 安装 Go 的 Protocol Buffers 插件：

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 运行项目

使用提供的启动脚本一键启动项目：

```bash
chmod +x start.sh
./start.sh
```

这将：
1. 启动一个由三个节点组成的 etcd 集群
2. 编译 Protocol Buffers 定义
3. 构建并启动三个 gRPC 服务实例
4. 运行客户端，测试服务发现和负载均衡
5. 完成测试后关闭所有组件

### 手动运行

如果需要手动运行各组件，可以按以下步骤操作：

1. 启动 etcd 集群：

```bash
docker compose up -d
```

2. 编译 Protocol Buffers 定义：

```bash
protoc --go_out=. --go-grpc_out=. proto/helloworld.proto
```

3. 启动多个服务实例：

```bash
go run server/server.go --port=50051
go run server/server.go --port=50052
go run server/server.go --port=50053
```

4. 运行客户端：

```bash
go run client/client.go
```

## 工作原理

### 服务注册

服务启动时，会自动向 etcd 注册自己的地址信息：

1. 创建带有 TTL 的 etcd 租约
2. 在 etcd 中使用键值对存储服务信息
3. 通过 KeepAlive 机制保持租约有效
4. 服务退出时自动注销

### 服务发现

客户端通过 etcd 发现可用的服务实例：

1. 客户端创建自定义的 gRPC 解析器
2. 解析器从 etcd 中查询可用服务地址
3. 监听 etcd 中的服务变更，动态更新地址列表

### 负载均衡

客户端使用 gRPC 内置的负载均衡功能：

1. 客户端使用 Round Robin 策略在多个服务实例间分发请求
2. 当服务实例变化时，客户端会自动更新实例列表

## 项目结构

```
etcd-grpc/
├── client/           # gRPC 客户端实现
├── etcd/             # etcd 工具库
│   ├── client.go     # etcd 客户端
│   ├── config.go     # 配置项
│   ├── discovery.go  # 服务发现
│   ├── registry.go   # 服务注册
│   ├── README.md     # etcd 服务注册发现实现详解
│   └── resolver.go   # gRPC 解析器
├── proto/            # Protocol Buffers 定义
├── server/           # gRPC 服务实现
├── docker-compose.yml # etcd 集群配置
└── start.sh          # 启动脚本
```

## 注意事项

- 本项目为演示目的而设计，生产环境部署可能需要进一步优化
- etcd 集群在生产环境中应当进行适当的安全配置
- 服务实例应根据实际需求进行扩展和管理

## 进阶应用

- 添加安全认证（TLS/SSL）
- 实现服务健康检查
- 增加监控和指标收集
- 优化负载均衡策略

## 许可证

[MIT License](LICENSE)
