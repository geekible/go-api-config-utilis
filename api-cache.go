package goapiconfigutilis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type ApiCache struct {
	hostname    string
	port        int
	expiration  time.Duration
	redisClient *redis.Client
}

func InitApiCache(hostname string, port int, expiration time.Duration) *ApiCache {
	return &ApiCache{
		hostname:   hostname,
		port:       port,
		expiration: expiration,
	}
}

func (c *ApiCache) Build(password string) {
	opts := redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.hostname, c.port),
		DB:   0,
	}

	if password != "" {
		opts.Password = password
	}

	c.redisClient = redis.NewClient(&opts)
}

func (c *ApiCache) SetValue(key string, value interface{}) error {
	return c.redisClient.Set(context.Background(), key, value, c.expiration).Err()
}

func (c *ApiCache) SetStructValue(key string, value any) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisClient.Set(context.Background(), key, json, c.expiration).Err()
}

func (c *ApiCache) GetValue(key string) (interface{}, error) {
	val, err := c.redisClient.Get(context.Background(), "name").Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (c *ApiCache) GetStructValue(key string, data any) error {
	val, err := c.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return err
	}

	return nil
}
