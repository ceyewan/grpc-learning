package etcd

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	GetConnection(ctx context.Context, serviceName string) (*grpc.ClientConn, error)
}

// EtcdDiscovery 实现基于etcd的服务发现
type EtcdDiscovery struct {
	client *Client
}

// NewServiceDiscovery 创建服务发现实例
func NewServiceDiscovery(client *Client) (*EtcdDiscovery, error) {
	if client == nil {
		var err error
		client, err = GetDefaultClient()
		if err != nil {
			return nil, err
		}
	}

	return &EtcdDiscovery{
		client: client,
	}, nil
}

// GetConnection 获取服务连接
func (d *EtcdDiscovery) GetConnection(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	// 注册解析器
	builder := &EtcdResolverBuilder{
		client:      d.client.client,
		serviceName: serviceName,
	}
	resolver.Register(builder)

	// 创建连接
	conn, err := grpc.NewClient(
		fmt.Sprintf("etcd:///%s", serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return conn, nil
}
