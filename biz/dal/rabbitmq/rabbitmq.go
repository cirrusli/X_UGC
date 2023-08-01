package rabbitmq

import (
	"X_UGC/conf"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"log/slog"
)

var (
	RMQ = new(rabbitMQ)
)

type rabbitMQ struct {
	channel       *amqp.Channel
	channelPool   map[string]*amqp.Channel //channel连接池
	mqConn        *amqp.Connection
	notifyConfirm chan amqp.Confirmation //确认发送到mq的channel
}

const (
	NotifyExchange = "notify_exchange" //通知交换机

	FollowProducerChannel   = "follow_producer_channel"   //关注通知队列
	GiveLikeProducerChannel = "giveLike_producer_channel" //点赞通知队列
	CommentProducerChannel  = "comment_producer_channel"  //评论通知队列
	FollowConsumerChannel   = "follow_consumer_channel"   //关注通知队列
	GiveLikeConsumerChannel = "giveLike_consumer_channel" //点赞通知队列
	CommentConsumerChannel  = "comment_consumer_channel"  //评论通知队列

	FollowQueue   = "follow_queue"   //关注通知队列
	GiveLikeQueue = "giveLike_queue" //点赞通知队列
	CommentQueue  = "comment_queue"  //评论通知队列

	FollowKey   = "follow_key"   //关注通知的routingKey
	GiveLikeKey = "giveLike_key" //点赞通知的routingKey
	CommentKey  = "comment_key"  //评论通知的routingKey
)

// InitRabbitMQ 初始胡MQ
func (mq *rabbitMQ) InitRabbitMQ(cfg *conf.RabbitMQ) (err error) {
	err = mq.ConnRabbitMQ(cfg)
	if err != nil {
		return err
	}

	//创建交换机，处理 关注，点赞，评论的通知交换机
	if err = mq.PrepareExchange(NotifyExchange, amqp.ExchangeDirect); err != nil {
		return err
	}

	//初始化三个队列，分别是处理 关注，点赞，评论的通知队列
	if _, err = mq.PrepareQueue(FollowQueue); err != nil {
		return err
	}
	if _, err = mq.PrepareQueue(GiveLikeQueue); err != nil {
		return err
	}
	if _, err = mq.PrepareQueue(CommentQueue); err != nil {
		return err
	}

	//交换机和三个队列之间bind
	if err = mq.QueueBindExchange(FollowQueue, FollowKey, NotifyExchange); err != nil {
		return err
	}
	if err = mq.QueueBindExchange(GiveLikeQueue, GiveLikeKey, NotifyExchange); err != nil {
		return err
	}
	if err = mq.QueueBindExchange(CommentQueue, CommentKey, NotifyExchange); err != nil {
		return err
	}

	return
}

// ConnRabbitMQ 建立rabbitmq连接
func (mq *rabbitMQ) ConnRabbitMQ(cfg *conf.RabbitMQ) (err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.UserName, cfg.Password, cfg.Host, cfg.Port)
	mq.mqConn, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	mq.channel, err = mq.mqConn.Channel()
	if err != nil {
		return err
	}
	if err = mq.channel.Qos(1, 0, true); err != nil {
		log.Println("Queue Consume: ", err.Error())
		return err
	}
	mq.channelPool = make(map[string]*amqp.Channel)
	//初始化confirm信道
	mq.notifyConfirm = make(chan amqp.Confirmation, 128)
	var channels = [6]string{FollowProducerChannel, GiveLikeProducerChannel, CommentProducerChannel, FollowConsumerChannel, GiveLikeConsumerChannel, CommentConsumerChannel}
	for _, channelName := range channels {
		mq.channelPool[channelName], err = mq.mqConn.Channel()
		if err != nil {
			return err
		}
		//设置不公平分发
		if err := mq.channelPool[channelName].Qos(1, 0, true); err != nil {
			log.Println("Qos error: ", err.Error())
			return err
		}
		//confirm开启发布确认
		if err := mq.channelPool[channelName].Confirm(false); err != nil {
			log.Println("Confirm error: ", err.Error())
			return err
		}
		//异步发布确认
		mq.channelPool[channelName].NotifyPublish(mq.notifyConfirm)
	}

	return
}

// PrepareExchange 准备rabbitmq的Exchange
func (mq *rabbitMQ) PrepareExchange(exchangeName, exchangeType string) error {
	err := mq.channel.ExchangeDeclare(
		exchangeName, // exchange
		exchangeType, // type
		true,         // durable 是否持久化，默认持久需要根据情况选择
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	)

	if nil != err {
		return err
	}

	return nil
}

// PrepareQueue 直接初始化队列
func (mq *rabbitMQ) PrepareQueue(queueName string) (queue amqp.Queue, err error) {
	queue, err = mq.channel.QueueDeclare(
		queueName, //name
		true,      //durable，是否持久化，默认持久需要根据情况选择
		false,     //delete when unused
		false,     //exclusive
		false,     //no-wait
		nil,       //arguments
	)
	return
}

// QueueBindExchange 队列绑定exchange
func (mq *rabbitMQ) QueueBindExchange(queueName, routingKey, exchangeName string) error {
	return mq.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

// Close 关闭连接
func (mq *rabbitMQ) Close() {
	err := mq.channel.Close()
	if err != nil {
		slog.Error("channel close error: ", err.Error())
		return
	}
	for _, channel := range mq.channelPool {
		err = channel.Close()
		if err != nil {
			slog.Error("channel close error: ", err.Error())
			return
		}
	}
	err = mq.mqConn.Close()
	if err != nil {
		slog.Error("mqConn close error: ", err.Error())
		return
	}
}

// ExchangeSend 通过exchange发送消息
func (mq *rabbitMQ) ExchangeSend(channelName, exchangeName, routingKey string, publishing amqp.Publishing) error {
	return mq.channelPool[channelName].Publish(
		exchangeName, //exchangeName
		routingKey,   //routing key
		false,        //mandatory
		false,        //immediate
		publishing,
	)
}

// ListenConfirm 异步发布确认
func (mq *rabbitMQ) ListenConfirm() {
	for {
		confirm := <-mq.notifyConfirm
		if !confirm.Ack {
			//这里表示消息发送到mq失败,可以处理失败流程
			log.Printf("message sending failure,it's number is : %d  ", confirm.DeliveryTag)
		}
	}
}

// QueueSend 通过队列发送消息
func (mq *rabbitMQ) QueueSend(channelName, queueName string, publishing amqp.Publishing) error {
	return mq.channelPool[channelName].Publish(
		"",        //exchangeName
		queueName, //queue name
		false,     //mandatory
		false,     //immediate
		publishing,
	)
}

// QueueConsume 消费队列,内部方法会阻塞,使用时需要单独启用一个线程处理，常驻后台执行
func (mq *rabbitMQ) QueueConsume(channelName, queueName string) (delivery <-chan amqp.Delivery, err error) {
	//后期可调整参数
	delivery, err = mq.channelPool[channelName].Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Println("Queue Consume: ", err.Error())
		return nil, err
	}
	return delivery, nil
}
