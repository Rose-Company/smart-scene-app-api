package redis

import (
	"context"
	"fmt"
	"time"

	config "smart-scene-app-api/config"

	redis "github.com/redis/go-redis/v9"
)

const (
	ErrRecordNotFound = redis.Nil
)

type ClientI interface {
	Get(ctx context.Context, key string) (string, error)
	GetByte(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value string, expiry int64) error
	SetByte(ctx context.Context, key string, value []byte, expiry int64) error
	Delete(ctx context.Context, keys ...string) error
	Publish(ctx context.Context, channel string, message string) error
	Subscribe(ctx context.Context, channel string) *redis.PubSub
}

type redisClient struct {
	client redis.UniversalClient
}

func NewRedisClient() ClientI {
	configRedis := config.Config.Redis
	var addrs = []string{fmt.Sprintf("%v:%v", configRedis.Host, configRedis.Port)}
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: configRedis.Pass,
		DB:       configRedis.DB,
	})
	return &redisClient{
		client: client,
	}
}

func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisClient) Set(ctx context.Context, key string, value string, expiry int64) error {
	return c.client.Set(ctx, key, value, time.Duration(expiry)*time.Second).Err()
}

func (c *redisClient) GetByte(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *redisClient) SetByte(ctx context.Context, key string, value []byte, expiry int64) error {
	return c.client.Set(ctx, key, value, time.Duration(expiry)*time.Second).Err()
}

func (c *redisClient) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *redisClient) Publish(ctx context.Context, channel string, message string) error {
	return c.client.Publish(ctx, channel, message).Err()
}

func (c *redisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.client.Subscribe(ctx, channel)
}
