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

// AddComment  添加评论
func AddComment(c *gin.Context) {
	userid := c.GetInt("userid")
	var comment model.Comment
	c.BindJSON(&comment)
	comment.UserID = userid
	err := service.AddComment(&comment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = service.IncrCommentCount(comment.ArticleID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//通知有人评论
	userInfo, err := service.GetAUserInfoByUserId(comment.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	articleInfo, err := service.GetArticleById(comment.ArticleID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if userid != comment.ReplyCommentID {
		loc, _ := time.LoadLocation("Local")
		sendTime, _ := time.ParseInLocation("2006-01-02 15:04:05", comment.CommentTime, loc)
		var commentNotify = &model.CommentNotifyInfo{
			NotifyInfo: model.NotifyInfo{
				SendTime: strconv.FormatInt(sendTime.UnixNano(), 10),
				FromUser: userInfo,
				ToUserID: comment.ReplyUserID,
			},
			ArticleInfo:    articleInfo,
			Content:        comment.Content,
			ReplyCommentID: comment.ReplyCommentID,
		}
		err = rabbitmq.RMQ.CommentProducer(commentNotify)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// RemovePermission 是否有删除权限
func RemovePermission(c *gin.Context) {
	userid := c.GetInt("userid")
	strCommentUserId := c.Query("userid")
	commentUserId, _ := strconv.Atoi(strCommentUserId)
	strAuthorId := c.Query("author_id")
	authorId, _ := strconv.Atoi(strAuthorId)
	if userid != commentUserId && userid != authorId {
		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"isMyComment": "false",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"isMyComment": "true",
		})
	}
}

// DeleteComment  删除评论
func DeleteComment(c *gin.Context) {
	strCommentId := c.Query("comment_id")
	strArticleId := c.Query("article_id")
	strLevel := c.Query("level")
	commentId, _ := strconv.Atoi(strCommentId)
	articleId, _ := strconv.Atoi(strArticleId)
	level, _ := strconv.Atoi(strLevel)
	err := service.DeleteComment(commentId, level)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = service.DecrCommentCount(articleId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// IncrCommentGiveLike 评论点赞
func IncrCommentGiveLike(c *gin.Context) {
	strCommentId := c.Query("comment_id")
	userid := c.GetInt("userid")
	commentId, _ := strconv.Atoi(strCommentId)
	err := service.IncrCommentGiveLike(userid, commentId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//通知有人评论点赞
	comment, err := service.GetCommentById(commentId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := service.GetAUserInfoByUserId(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	articleInfo, err := service.GetArticleById(comment.ArticleID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if userid != comment.UserID {
		var commentNotify = &model.CommentNotifyInfo{
			NotifyInfo: model.NotifyInfo{
				SendTime: strconv.FormatInt(time.Now().UnixNano(), 10),
				FromUser: userInfo,
				ToUserID: comment.UserID,
			},
			ArticleInfo:    articleInfo,
			Content:        "[#@give_like@#]",
			ReplyCommentID: commentId,
		}
		err = rabbitmq.RMQ.CommentProducer(commentNotify)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DecrCommentGiveLike 取消评论点赞
func DecrCommentGiveLike(c *gin.Context) {
	strCommentId := c.Query("comment_id")
	userid := c.GetInt("userid")
	commentId, _ := strconv.Atoi(strCommentId)
	err := service.DecrCommentGiveLike(userid, commentId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// IsCommentGiveLike  判断评论是否点赞
func IsCommentGiveLike(c *gin.Context) {
	strCommentId := c.Query("comment_id")
	userid := c.GetInt("userid")
	commentId, _ := strconv.Atoi(strCommentId)
	isGiveLike, err := service.IsCommentGiveLike(userid, commentId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"isGiveLike": isGiveLike,
	})
}

// GetCommentByHeat 根据热度获一级评论
func GetCommentByHeat(c *gin.Context) {
	userid := c.GetInt("userid")
	strArticleId := c.Query("article_id")
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	articleId, _ := strconv.Atoi(strArticleId)
	commentList, err := service.GetCommentByHeat(userid, articleId, pageIndex, pageSize)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"commentList": commentList,
		})
	}
}

// GetReplyCommentByHeat 根据热度获二级评论
func GetReplyCommentByHeat(c *gin.Context) {
	userid := c.GetInt("userid")
	strArticleId := c.Query("article_id")
	articleId, _ := strconv.Atoi(strArticleId)
	strCommentId := c.Query("comment_id")
	commentId, _ := strconv.Atoi(strCommentId)
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	commentList, err := service.GetReplyCommentByHeat(userid, articleId, commentId, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"commentList": commentList,
		})
	}
}

// GetCommentById 根据commentId获取评论信息
func GetCommentById(c *gin.Context) {
	strCommentId := c.Query("comment_id")
	commentId, _ := strconv.Atoi(strCommentId)
	comment, err := service.GetCommentById(commentId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"comment": comment,
	})
}
