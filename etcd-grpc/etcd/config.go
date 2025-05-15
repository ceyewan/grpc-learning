package etcd

import "time"

// Config 定义etcd客户端配置
type Config struct {
	Endpoints   []string      // etcd服务器地址列表
	DialTimeout time.Duration // 连接超时时间
	LogLevel    string        // 日志级别
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Endpoints:   []string{"localhost:23791", "localhost:23792", "localhost:23793"},
		DialTimeout: 5 * time.Second,
		LogLevel:    "info",
	}
}
