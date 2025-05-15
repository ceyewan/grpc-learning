package etcd

import (
	"context"
	"fmt"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	Register(ctx context.Context, serviceName, instanceID, addr string) error
	Deregister(ctx context.Context, serviceName, instanceID string) error
}

// EtcdRegistry 实现基于etcd的服务注册
type EtcdRegistry struct {
	client *Client
}

// NewServiceRegistry 创建服务注册实例
func NewServiceRegistry(client *Client) (*EtcdRegistry, error) {
	if client == nil {
		var err error
		client, err = GetDefaultClient()
		if err != nil {
			return nil, err
		}
	}

	return &EtcdRegistry{
		client: client,
	}, nil
}

// Register 注册服务
func (r *EtcdRegistry) Register(ctx context.Context, serviceName, instanceID, addr string) error {
	// 创建租约
	lease, err := r.client.client.Grant(ctx, 5)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	// 设置保持活动
	keepAliveCh, err := r.client.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return fmt.Errorf("failed to setup keepalive: %w", err)
	}

	// 处理keepalive响应
	go func() {
		for {
			select {
			case resp, ok := <-keepAliveCh:
				if !ok {
					log.Printf("Keepalive channel closed for service %s/%s", serviceName, instanceID)
					return
				}
				log.Printf("Lease renewed for service %s/%s, TTL: %d", serviceName, instanceID, resp.TTL)
			case <-ctx.Done():
				log.Printf("Service registry context canceled for %s/%s", serviceName, instanceID)
				return
			}
		}
	}()

	// 写入服务信息
	key := fmt.Sprintf("/services/%s/%s", serviceName, instanceID)
	_, err = r.client.client.Put(ctx, key, addr, clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	log.Printf("Service registered successfully: %s/%s at %s", serviceName, instanceID, addr)
	return nil
}

// Deregister 注销服务
func (r *EtcdRegistry) Deregister(ctx context.Context, serviceName, instanceID string) error {
	key := fmt.Sprintf("/services/%s/%s", serviceName, instanceID)
	_, err := r.client.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	log.Printf("Service deregistered: %s/%s", serviceName, instanceID)
	return nil
}
