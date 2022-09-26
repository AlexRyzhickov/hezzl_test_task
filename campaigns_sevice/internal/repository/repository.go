package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"hezzl_test_task/campaigns_sevice/internal/models"
	"time"
)

type RedisCache struct {
	*redis.Client
}

func (c *RedisCache) Load(key string) (models.Item, bool) {
	item := models.Item{}
	val, err := c.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return item, false
	}
	if err != nil {
		return item, false
	}
	err = json.Unmarshal([]byte(val), &item)
	if err != nil {
		return item, false
	}
	return item, true
}

func (c *RedisCache) Store(ctx context.Context, key string, value models.Item) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.Set(ctx, key, jsonVal, time.Second*10).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
