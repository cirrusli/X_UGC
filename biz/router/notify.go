package router

import (
	"X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"github.com/gin-gonic/gin"
)

func NotifyGroups(r *gin.Engine) {
	NotifyGroup := r.Group("/notify", mw.Auth())
	{
		//获取通知历史消息
		NotifyGroup.GET("/getNotifyHistory", handler.GetNotifyHistory)
		//获取未读通知消息数量
		NotifyGroup.GET("/getUnreadNotify", handler.GetUnreadNotify)
		//重置未读通知消息数量
		NotifyGroup.PUT("/resetUnreadNotify", handler.ResetUnreadNotify)
		//获取未读聊天信息通知
		NotifyGroup.GET("/getUnreadChat", handler.GetUnreadChat)
		//重置未读聊天信息数目
		NotifyGroup.PUT("/resetUnreadChat", handler.ResetUnreadChat)
	}
}
