package dal

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

// InitRedis 初始化连接
func InitRedis(cfg *conf.Redis) (err error) {
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

func CloseRedis() {
	err := RDB.Close()
	if err != nil {
		return
	}
}
