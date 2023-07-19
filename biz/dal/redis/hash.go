package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// HSet redis写入hash
func HSet(key string, val ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.HSet(ctx, key, val).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

// HGet	redis读取单个hash
func HGet(key string, field string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	val, err := RDB.HGet(ctx, key, field).Result()
	defer cancel()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// HGetAll redis读取所有hash
func HGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vals, err := RDB.HGetAll(ctx, key).Result()
	defer cancel()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return vals, nil
}

// HDel redis删除hash表中的元素
func HDel(key string, field string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.HDel(ctx, key, field).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

// HExists  redis判断元素是否存在
func HExists(key string, field string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := RDB.HExists(ctx, key, field).Result()
	defer cancel()
	if err != nil {
		return -1, err
	}
	if !result {
		return 0, nil
	}
	return 1, nil
}

// HLen   redis获取长度
func HLen(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	Len, err := RDB.HLen(ctx, key).Result()
	defer cancel()
	if err != nil {
		return -1, err
	}
	return Len, nil
}

// HIncrBy redis的hash增减值数字
func HIncrBy(key string, field string, incrNum int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.HIncrBy(ctx, key, field, incrNum).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}
