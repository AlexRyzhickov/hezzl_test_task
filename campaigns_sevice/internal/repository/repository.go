package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"hezzl_test_task/internal/models"
	"time"
)

type RedisCache struct {
	*redis.Client
}

func (c *RedisCache) Load(key string) ([]models.Item, error) {
	res := []models.Item{}
	val, err := c.Get(context.Background(), key).Result()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return res, err
	}
	return res, err
}

func (c *RedisCache) Store(ctx context.Context, key string, value []models.Item) error {
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
