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

var DEFAULT_ID = "Redis"
var DEFAULT_ADDR = "localhost:6379"
var DEFAULT_PASSWORD = ""

func CreateRedisInstance(id string) (*redis.Client, error) {
	addr := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	if id == "" {
		id = DEFAULT_ID
	}

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
		return nil, err
	}

	log.Printf("Connected to %s successfully", id)
	return rdb, nil
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

// Add a single value into the set defined by a "key" name
func AddToSet(key string, value string) error {
	if err := rdb.SAdd(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

// Remove existing value from the set defined by a "key" name
func RemoveFromSet(key string, value string) error {
	if err := rdb.SRem(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func GetSetMemberValues(key string) ([]string, error) {
	result, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
