package cmd

import (
	"X_UGC/conf"
	"X_UGC/dal"
	"X_UGC/dal/redis"
	"X_UGC/router"
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

	r := router.Init()
	err = r.Run(fmt.Sprintf(":%d", conf.C.Port))
	if err != nil {
		log.Fatalln("Server start failed: ", err)
	}
}
