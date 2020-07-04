package services

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

type Subscriber interface {
	Subscribe(queue string, payload []byte) error
}

type RedisSubscriber struct {
	client *redis.Client
	PubSub *redis.PubSub
}

func NewRedisSubscriber(queue string) *RedisSubscriber {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	fmt.Printf("Subscribe -> %s\n", queue)
	pubSub := redisClient.Subscribe(queue)
	if pubSub == nil {
		panic(errors.New("subscribe-error"))
	}

	// // Wait for confirmation that subscription is created before publishing anything.
	// _, err := pubSub.Receive()
	// if err != nil {
	// 	panic(err)
	// }

	return &RedisSubscriber{
		client: redisClient,
		PubSub: pubSub,
	}
}

// func (s *RedisSubscriber) Subscribe(queue string) error {

// 	return nil
// }

func (s *RedisSubscriber) Close() error {
	return s.PubSub.Close()
}

// Consume messages.
func (s *RedisSubscriber) Consume(out chan string) {
	fmt.Println("Consume")
	for msg := range s.PubSub.Channel() {
		fmt.Println("\tGot from", msg.Channel, msg.Payload) // [:30]
		out <- msg.Payload
	}
}
