#!/bin/bash
set -e

# 启动 docker compose

echo "[INFO] 启动 docker compose..."
docker compose up -d

echo "[INFO] 生成 gRPC 代码..."
protoc --go_out=. --go-grpc_out=. proto/helloworld.proto

# 编译服务端
go build -o server_bin server/server.go

sleep 5 # 等待docker compose 完全启动

# 启动多个服务端实例
PORTS=(50051 50052 50053)
PIDS=()
echo "[INFO] 启动服务端实例..."
for port in "${PORTS[@]}"; do
    ./server_bin --port=$port &
    PIDS+=("$!")
    sleep 0.5 # 可根据实际情况调整，确保端口不冲突
    echo "[INFO] 服务端已启动: $port"
done

sleep 5 # 等待服务端完全启动
# 打印所有服务端进程的 PID
echo "[INFO] 服务端进程 PID: ${PIDS[@]}"

echo "[INFO] 启动客户端..."
go run client/client.go

echo "[INFO] 客户端已结束，关闭所有服务端进程..."
for pid in "${PIDS[@]}"; do
    if kill -0 $pid 2>/dev/null; then
        kill $pid
        echo "[INFO] 已关闭服务端进程: $pid"
    fi
done

# 删除服务端二进制文件
rm -f server_bin

sleep 2 # 等待服务端进程完全关闭

echo "[INFO] 关闭 docker compose..."
docker compose down