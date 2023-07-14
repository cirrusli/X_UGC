package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// Set redis写入string
func Set(key string, val string, expTime int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.SetEX(ctx, key, val, time.Duration(expTime)*time.Second).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

// Get redis读取string
func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	val, err := RDB.Get(ctx, key).Result()
	defer cancel()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}
