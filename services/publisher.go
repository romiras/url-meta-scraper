package services

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Publisher interface {
	Publish(queue string, payload []byte) error
}

type RedisPublisher struct {
	client *redis.Client
}

func NewRedisPublisher() *RedisPublisher {
	return &RedisPublisher{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			PoolSize: 10,
		}),
	}
}

func (p *RedisPublisher) Publish(queue string, payload []byte) error {
	fmt.Printf("Publish -> %s: %s\n", queue, string(payload[:50]))
	cmd := p.client.Publish(queue, payload)
	return cmd.Err()
}
