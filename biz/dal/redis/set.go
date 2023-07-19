package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// SAdd redis添加set集合元素
func SAdd(key string, member string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.SAdd(ctx, key, member).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

// SRem redis删除set集合元素
func SRem(key string, member string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.SRem(ctx, key, member).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

// SMembers redis获取set集合所有元素
func SMembers(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	results, err := RDB.SMembers(ctx, key).Result()
	defer cancel()
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SIsMember   redis判断元素是否在set中
func SIsMember(key string, member string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	exists, err := RDB.SIsMember(ctx, key, member).Result()
	defer cancel()
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SRandMember redis随机取出set中1个值
func SRandMember(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := RDB.SRandMember(ctx, key).Result()
	defer cancel()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

// SRandMemberN redis随机取出set中n个值
func SRandMemberN(key string, count int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := RDB.SRandMemberN(ctx, key, count).Result()
	defer cancel()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
