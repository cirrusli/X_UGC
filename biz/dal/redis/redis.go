package redis

import (
	"X_UGC/conf"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	RDB *redis.Client
)

// Init 初始化连接
func Init(cfg *conf.Redis) (err error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.UserName,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = RDB.Ping(ctx).Result()
	return err
}

func Close() {
	err := RDB.Close()
	if err != nil {
		return
	}
}

// Del   删除key
func Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.Del(ctx, key).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

// Expire 设置过期时间
func Expire(key string, expTime int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := RDB.Expire(ctx, key, time.Duration(expTime)*time.Second).Err()
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}
