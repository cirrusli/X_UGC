package rabbitmq

import (
	ws "X_UGC/pkg/websocket"
	"log/slog"
)

// FollowConsumer 消费关注通知
func (mq *rabbitMQ) FollowConsumer() {
	delivery, err := mq.QueueConsume(FollowConsumerChannel, FollowQueue)
	if err != nil {
		slog.Info("消费错误:", err.Error())
	}
	for d := range delivery {
		ws.ProcessNotify(ws.FollowNotify, d.Body)
		err := d.Ack(false)
		if err != nil {
			slog.Info(err.Error())
		}
	}
}

// GiveLikeConsumer 消费点赞通知
func (mq *rabbitMQ) GiveLikeConsumer() {
	delivery, err := mq.QueueConsume(GiveLikeConsumerChannel, GiveLikeQueue)
	if err != nil {
		slog.Info("消费错误:", err.Error())
	}
	for d := range delivery {
		ws.ProcessNotify(ws.GiveLikeNotify, d.Body)
		err := d.Ack(false)
		if err != nil {
			slog.Info(err.Error())
		}
	}
}

// CommentConsumer 消费评论通知
func (mq *rabbitMQ) CommentConsumer() {
	delivery, err := mq.QueueConsume(CommentConsumerChannel, CommentQueue)
	if err != nil {
		slog.Info("消费错误:", err.Error())
	}
	for d := range delivery {
		ws.ProcessNotify(ws.CommentNotify, d.Body)
		err := d.Ack(false)
		if err != nil {
			slog.Info(err.Error())
		}
	}
}

func (mq *rabbitMQ) StartConsumers() {
	go mq.FollowConsumer()
	go mq.GiveLikeConsumer()
	go mq.CommentConsumer()
}
