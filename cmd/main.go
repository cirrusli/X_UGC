package main

import (
	"X_UGC/biz/dal"
	"X_UGC/biz/dal/rabbitmq"
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/router"
	"X_UGC/conf"
	ws "X_UGC/pkg/websocket"
	"fmt"
	"log"
	"os"
)

func main() {
	configPath := "./conf/config.yaml"
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	}

	if err := initConfig(configPath); err != nil {
		log.Fatalln("Failed to initialize config: ", err)
	}

	if err := startServer(); err != nil {
		log.Fatalln("Failed to start server: ", err)
	}
}

func initConfig(path string) error {
	if err := conf.Init(path); err != nil {
		return fmt.Errorf("load config file failed: %w", err)
	}

	if err := dal.InitMySQL(conf.C.MySQL); err != nil {
		return fmt.Errorf("init mysql failed: %w", err)
	}
	defer dal.CloseMySQL()
	dal.InitTables()

	if err := redis.Init(conf.C.Redis); err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}
	defer redis.Close()

	if err := rabbitmq.RMQ.InitRabbitMQ(conf.C.RabbitMQ); err != nil {
		return fmt.Errorf("init rabbitmq failed: %w", err)
	}
	go rabbitmq.RMQ.ListenConfirm()
	rabbitmq.RMQ.StartConsumers()
	defer rabbitmq.RMQ.Close()

	go ws.WsManager.Start()

	return nil
}

func startServer() error {
	r := router.Register()
	return r.Run(fmt.Sprintf(":%d", conf.C.Port))
}
