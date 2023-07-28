package router

import (
	"X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"github.com/gin-gonic/gin"
)

func CommentGroups(r *gin.Engine) {
	CommentGroup := r.Group("/comment", mw.Auth())
	{
		//添加评论
		CommentGroup.POST("/addComment", handler.AddComment)
		//是否有删除权限
		CommentGroup.GET("/removePermission", handler.RemovePermission)
		//删除评论
		CommentGroup.DELETE("/deleteComment", handler.DeleteComment)
		//评论点赞
		CommentGroup.PUT("/incrCommentGiveLike", handler.IncrCommentGiveLike)
		//取消评论点赞
		CommentGroup.PUT("/decrCommentGiveLike", handler.DecrCommentGiveLike)
		//判断评论是否点赞
		CommentGroup.GET("/isCommentGiveLike", handler.IsCommentGiveLike)
		//根据热度获一级评论
		CommentGroup.GET("/getCommentByHeat", handler.GetCommentByHeat)
		//根据热度获二级评论
		CommentGroup.GET("/getReplyCommentByHeat", handler.GetReplyCommentByHeat)
		//根据commentId获取评论信息
		CommentGroup.GET("/getCommentById", handler.GetCommentById)
	}
}
