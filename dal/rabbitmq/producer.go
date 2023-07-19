package rabbitmq

import (
	"X_UGC/model"
	"encoding/json"
	"github.com/streadway/amqp"
)

// FollowProducer 生产关注通知
func (mq *rabbitMQ) FollowProducer(followNotify *model.FollowNotifyInfo) (err error) {
	message, _ := json.Marshal(followNotify)
	err = mq.ExchangeSend(FollowProducerChannel, NotifyExchange, FollowKey, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	return
}

// GiveLikeProducer 生产点赞通知
func (mq *rabbitMQ) GiveLikeProducer(giveLikeNotify *model.GiveLikeNotifyInfo) (err error) {
	message, _ := json.Marshal(giveLikeNotify)
	err = mq.ExchangeSend(GiveLikeProducerChannel, NotifyExchange, GiveLikeKey, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	return
}

// CommentProducer 生产评论通知
func (mq *rabbitMQ) CommentProducer(commentNotify *model.CommentNotifyInfo) (err error) {
	message, _ := json.Marshal(commentNotify)
	err = mq.ExchangeSend(CommentProducerChannel, NotifyExchange, CommentKey, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	return
}
