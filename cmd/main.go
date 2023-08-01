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
	if len(os.Args) < 2 {
		log.Fatalln("Missing parameter ./conf/config.yaml")
	}
	err := conf.Init(os.Args[1])
	if err != nil {
		log.Fatalln("Load config file failed: ", err)
	}
	err = dal.InitMySQL(conf.C.MySQL)
	if err != nil {
		log.Fatalln("Init mysql failed: ", err)
	}
	defer dal.CloseMySQL()
	dal.InitTables()

	err = redis.Init(conf.C.Redis)
	if err != nil {
		log.Fatalln("Init redis failed: ", err)
	}
	defer redis.Close()

	err = rabbitmq.RMQ.InitRabbitMQ(conf.C.RabbitMQ)
	if err != nil {
		fmt.Printf("init rabbitmq failed, err:%v\n", err)
		return
	}
	//启动协程监听confirm发布确认
	go rabbitmq.RMQ.ListenConfirm()

	rabbitmq.RMQ.StartConsumers()

	defer rabbitmq.RMQ.Close()

	//启动ws服务
	go ws.WsManager.Start()

	r := router.Register()
	err = r.Run(fmt.Sprintf(":%d", conf.C.Port))
	if err != nil {
		log.Fatalln("Server start failed: ", err)
	}
}
