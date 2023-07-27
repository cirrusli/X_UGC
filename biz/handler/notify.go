package handler

import (
	"X_UGC/biz/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetNotifyHistory 获取通知历史消息
func GetNotifyHistory(c *gin.Context) {
	userid := c.GetInt("userid")
	notifyType := c.Query("notifyType")
	strPageIndex := c.Query("pageIndex")
	strPageSize := c.Query("pageSize")
	pageIndex, _ := strconv.ParseInt(strPageIndex, 10, 64)
	pageSize, _ := strconv.ParseInt(strPageSize, 10, 64)
	notifyHistory, err := service.GetNotifyHistory(userid, notifyType, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":                      "success",
		"notifyHistory_" + notifyType: notifyHistory,
	})
}

// GetUnreadNotify  获取未读通知消息数量
func GetUnreadNotify(c *gin.Context) {
	userid := c.GetInt("userid")
	unreadNotify, err := service.GetUnreadNotify(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"unreadNotify": unreadNotify,
	})
}

// ResetUnreadNotify    重置未读通知消息数量
func ResetUnreadNotify(c *gin.Context) {
	userid := c.GetInt("userid")
	notifyType := c.Query("notifyType")
	err := service.ResetUnreadNotify(userid, notifyType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// GetUnreadChat 获取未读聊天信息通知
func GetUnreadChat(c *gin.Context) {
	userid := c.GetInt("userid")
	unreadChatList, err := service.GetUnreadChat(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"unreadChatList": unreadChatList,
	})
}

// ResetUnreadChat 重置未读聊天信息数目
func ResetUnreadChat(c *gin.Context) {
	userid := c.GetInt("userid")
	chatUserid := c.Query("chat_userid")
	err := service.ResetUnreadChat(strconv.Itoa(userid), chatUserid, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
