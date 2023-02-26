package cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

type Cache struct {
	Client *redis.Client
}

func NewCache() (*Cache, error) {
	host := os.Getenv("REDISHOST")
	port := os.Getenv("REDISPORT")
	password := os.Getenv("REDISPASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("could not ping the cache: %v", err)
	}
	return &Cache{Client: client}, nil
}

type KeyNotFound struct {
	Key string
}

func (k KeyNotFound) Error() string {
	return fmt.Sprintf("value not found for given key: %s", k.Key)
}

func (c *Cache) Get(key string) (value string, found bool, err error) {
	val, err := c.Client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to get value for key %s. reason: %v", key, err)
	}
	return val, true, nil
}

func (c *Cache) Set(key, value string) error {
	err := c.Client.Set(key, value, 5*time.Minute).Err()
	return fmt.Errorf("failed to set value for key: %s and value %s. reason: %v", key, value, err)
}
