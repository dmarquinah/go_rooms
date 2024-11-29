package cache

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
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
