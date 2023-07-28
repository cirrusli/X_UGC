package router

import (
	"X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"github.com/gin-gonic/gin"
)

func RelationInfoGroups(r *gin.Engine) {
	RelationInfoGroup := r.Group("/relationInfo", mw.Auth())
	{
		//获取用户关系信息
		RelationInfoGroup.GET("/getRelation", handler.GetRelation)
		//更改用户关系信息
		RelationInfoGroup.PUT("/changeRelationStatus", handler.ChangeRelationStatus)
		//获取关注的人的列表
		RelationInfoGroup.GET("/getFollowList", handler.GetFollowList)
		//获取粉丝列表
		RelationInfoGroup.GET("/getFansList", handler.GetFansList)
		//获取关注和粉丝的总数
		RelationInfoGroup.GET("/getFollowAndFansCount", handler.GetFollowAndFansCount)
	}
}
