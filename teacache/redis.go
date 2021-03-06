package teacache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/iwind/TeaGo/maps"
	"time"
)

// 内存缓存管理器
type RedisManager struct {
	Capacity float64       // 容量
	Life     time.Duration // 有效期

	Network  string
	Host     string
	Port     int
	Password string
	Sock     string

	client *redis.Client
}

func NewRedisManager() *RedisManager {
	m := &RedisManager{}
	return m
}

func (this *RedisManager) SetOptions(options map[string]interface{}) {
	m := maps.NewMap(options)
	this.Network = m.GetString("network")
	this.Host = m.GetString("host")
	this.Port = m.GetInt("port")
	this.Password = m.GetString("password")
	this.Sock = m.GetString("sock")

	addr := ""
	if this.Network == "tcp" {
		if this.Port > 0 {
			addr = fmt.Sprintf("%s:%d", this.Host, this.Port)
		} else {
			addr = this.Host + ":6379"
		}
	} else if this.Network == "sock" {
		addr = this.Sock
	}

	this.client = redis.NewClient(&redis.Options{
		Network:      this.Network,
		Addr:         addr,
		Password:     this.Password,
		DialTimeout:  10 * time.Second, // TODO 换成可配置
		ReadTimeout:  10 * time.Second, // TODO 换成可配置
		WriteTimeout: 10 * time.Second, // TODO 换成可配置
		TLSConfig:    nil,              // TODO 支持TLS
	})
}

func (this *RedisManager) Write(key string, data []byte) error {
	cmd := this.client.Set("TEA_CACHE_"+key, string(data), this.Life)
	return cmd.Err()
}

func (this *RedisManager) Read(key string) (data []byte, err error) {
	cmd := this.client.Get("TEA_CACHE_" + key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return []byte(cmd.Val()), nil
}
