package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

var C = new(App)

// App 应用程序配置
type App struct {
	Release   bool `yaml:"release"`
	Port      int  `yaml:"port"`
	*MySQL    `yaml:"mysql"`
	*Redis    `yaml:"redis"`
	*RabbitMQ `yaml:"rabbitmq"`
	*ES       `yaml:"es"`
}

// MySQL 数据库配置
type MySQL struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`   //数据库名
	Host     string `yaml:"host"` //地址
	Port     int    `yaml:"port"` //端口
}

// Redis redis配置
type Redis struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`     //地址
	Port     int    `yaml:"port"`     //端口号
	DB       int    `yaml:"db"`       //redis数据库下标
	PoolSize int    `yaml:"poolSize"` //连接池大小
}

// RabbitMQ  rabbitmq配置
type RabbitMQ struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"` //地址
	Port     int    `yaml:"port"` //端口号
}
type ES struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"` //地址
	Port     int    `yaml:"port"` //端口号
}

func Init(file string) error {
	f, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, C)
}
