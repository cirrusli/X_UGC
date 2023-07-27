package handler

import (
	"X_UGC/biz/dal/rabbitmq"
	"X_UGC/biz/model"
	"X_UGC/biz/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetRelation 获取用户关系信息
func GetRelation(c *gin.Context) {
	userid := c.GetInt("userid")
	strUserid := strconv.Itoa(userid)
	objectId := c.Query("relation_id")
	relationInfo, err := service.GetRelationInfoByUserId(strUserid, objectId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"relationInfo": relationInfo,
		})
	}
}

// ChangeRelationStatus		更改好友之间状态
func ChangeRelationStatus(c *gin.Context) {
	userid := c.GetInt("userid")
	strUserid := strconv.Itoa(userid)
	objectId := c.Query("relation_id")
	status, err := service.GetRelationInfoStatus(strUserid, objectId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = service.ChangeRelationInfoStatus(strUserid, objectId, status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//通知有人关注
	if status == model.NoRelationship || status == model.Fans {

		fromUserId, _ := strconv.Atoi(strUserid)
		toUserId, _ := strconv.Atoi(objectId)
		userInfo, err := service.GetAUserInfoByUserId(fromUserId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		var followNotify = &model.FollowNotifyInfo{
			NotifyInfo: model.NotifyInfo{
				SendTime: strconv.FormatInt(time.Now().UnixNano(), 10),
				FromUser: userInfo,
				ToUserID: toUserId,
			},
		}
		err = rabbitmq.RMQ.FollowProducer(followNotify)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

// GetFollowList 获取关注用户的列表
func GetFollowList(c *gin.Context) {
	userid := c.GetInt("userid")
	strUserid := strconv.Itoa(userid)
	FollowList, err := service.GetFollowList(strUserid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"followList": FollowList,
		})
	}
}

// GetFansList 获取关注用户的列表
func GetFansList(c *gin.Context) {
	userid := c.GetInt("userid")
	strUserid := strconv.Itoa(userid)
	FansList, err := service.GetFansList(strUserid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"followList": FansList,
		})
	}
}

// GetFollowAndFansCount 获取关注和粉丝的总数
func GetFollowAndFansCount(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		userid = strconv.Itoa(c.GetInt("userid"))
	}
	followFansCount, err := service.GetFollowFansCount(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"follow_fans_count": followFansCount,
	})
}
