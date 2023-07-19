package router

import (
	"X_UGC/handler"
	"X_UGC/mw"
	"github.com/gin-gonic/gin"
)

func SearchGroups(r *gin.Engine) {
	SearchGroup := r.Group("/search", mw.Auth())
	{
		//根据用户名搜索用户
		//SearchGroup.GET("/searchUserByUserName", controller.SearchUserByUserName)
		//根据账号和用户名搜索用户
		SearchGroup.GET("/searchUser", handler.SearchUser)
		//根据标题或内容搜索文章
		SearchGroup.GET("/searchArticle", handler.SearchArticle)
		//根据文章类型以及标题或内容搜索文章
		SearchGroup.GET("/searchArticleByType", handler.SearchArticleByType)
		//删除搜索记录
		SearchGroup.DELETE("/delSearchRecord", handler.DelSearchRecord)
		//删除全部搜索记录
		SearchGroup.DELETE("/delAllSearchRecord", handler.DelAllSearchRecord)
		//获取所有搜索记录
		SearchGroup.GET("/getAllSearchRecord", handler.GetAllSearchRecord)
	}
}
