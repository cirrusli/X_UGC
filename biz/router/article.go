package router

import (
	"X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"github.com/gin-gonic/gin"
)

func ArticlesGroups(r *gin.Engine) {
	ArticlesGroup := r.Group("/article", mw.Auth())
	{
		//发布文章
		ArticlesGroup.POST("/publishedArticle", handler.PublishedArticle)
		//获取用户所有文章
		ArticlesGroup.GET("/getAllArticle", handler.GetAllArticle)
		//分页获取用户所有文章
		ArticlesGroup.GET("/getAllArticleByPage", handler.GetAllArticleByPage)
		//根据ID获取文章信息
		ArticlesGroup.GET("/getArticleById", handler.GetArticleById)
		//文章点赞
		ArticlesGroup.PUT("/giveLike", handler.GiveLike)
		//取消文章赞
		ArticlesGroup.PUT("/delLike", handler.DelLike)
		//判断文章是否点赞
		ArticlesGroup.GET("/isGiveLike", handler.IsGiveLike)
		//获取点赞文章总数
		ArticlesGroup.GET("/getGiveLikeArticleCount", handler.GetGiveLikeArticleCount)
		//获取用户的点赞文章
		ArticlesGroup.GET("/getGiveLikeArticle", handler.GetGiveLikeArticle)
		//获取所有文章类型
		ArticlesGroup.GET("/getAllArticleType", handler.GetAllArticleType)

		//检查视频文章是否存在
		ArticlesGroup.GET("/articleExist", handler.ArticleExist)
		//分块上传文件
		ArticlesGroup.POST("/chunkUpload", handler.ChunkUpload)
		//合并切片
		ArticlesGroup.POST("/mergeChunk", handler.MergeChunk)
		//取消上传
		ArticlesGroup.POST("/cancelChunk", handler.CancelChunk)

		//feed流文章推送
		ArticlesGroup.GET("/pushArticleFeed", handler.PushArticleFeed)
		//具体类型文章推送
		ArticlesGroup.GET("/pushArticleByType", handler.PushArticleByType)
		//分页获取朋友的文章
		ArticlesGroup.GET("/getArticleFromFriend", handler.GetArticleFromFriend)
		//分页获取关注的用户的文章
		ArticlesGroup.GET("/getArticleFromFollow", handler.GetArticleFromFollow)

	}
}
