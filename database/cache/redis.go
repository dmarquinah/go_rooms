package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

var DEFAULT_ADDR = "localhost:6379"
var DEFAULT_PASSWORD = ""

func InitRedis() error {
	addr := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	if addr == "" {
		addr = DEFAULT_ADDR
	}

	if password == "" {
		password = DEFAULT_PASSWORD
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		return err
	}

	log.Println("Connected to Redis successfully")
	return nil
}

// Publish message to Redis channel
func PublishMessage(channel string, message []byte) error {
	return rdb.Publish(ctx, channel, message).Err()
}

// Subscribe to a Redis channel
func SubscribeToChannel(channel string, handler func(message map[string]string)) {
	sub := rdb.Subscribe(ctx, channel)

	// Process messages
	go func() {
		for msg := range sub.Channel() {
			var data map[string]string
			if err := json.Unmarshal([]byte(msg.Payload), &data); err == nil {
				handler(data)
			} else {
				log.Printf("Failed to deserialize Redis message: %v", err)
				return
			}
		}
	}()
}
