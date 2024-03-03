package service

import (
	"X_UGC/biz/dal/mysql"
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"github.com/jinzhu/gorm"
	"strconv"
)

// AddComment  添加评论
func AddComment(comment *model.Comment) (err error) {
	err = mysql.DB.Create(comment).Error
	return
}

// DeleteComment  删除评论
func DeleteComment(commentId int, level int) (err error) {
	err = mysql.DB.Where("id=?", commentId).Delete(&model.Comment{}).Error
	if err != nil {
		return err
	}
	if level == 1 {
		err = mysql.DB.Where("level=2 and comment_group=?", commentId).Delete(&model.Comment{}).Error
	} else if level == 2 {
		err = mysql.DB.Where("level=2 and reply_comment_id=?", commentId).Delete(&model.Comment{}).Error
	}
	return
}

// IncrCommentGiveLike 评论点赞
func IncrCommentGiveLike(userid int, commentId int) (err error) {
	err = mysql.DB.Model(&model.Comment{}).Where("id=?", commentId).Update("give_like_count", gorm.Expr("give_like_count + ?", 1)).Error
	if err != nil {
		return err
	}
	err = redis.SAdd(model.COMMENT_GIVELIKE+strconv.Itoa(userid), strconv.Itoa(commentId))
	return
}

// DecrCommentGiveLike 取消评论点赞
func DecrCommentGiveLike(userid int, commentId int) (err error) {
	err = mysql.DB.Model(&model.Comment{}).Where("id=?", commentId).Update("give_like_count", gorm.Expr("give_like_count - ?", 1)).Error
	if err != nil {
		return err
	}
	err = redis.SRem(model.COMMENT_GIVELIKE+strconv.Itoa(userid), strconv.Itoa(commentId))
	return
}

// IsCommentGiveLike  判断评论是否点赞
func IsCommentGiveLike(userid int, commentId int) (exist bool, err error) {
	exist, err = redis.SIsMember(model.COMMENT_GIVELIKE+strconv.Itoa(userid), strconv.Itoa(commentId))
	return
}

// GetCommentById 根据commentId获取评论信息
func GetCommentById(commentId int) (comment *model.Comment, err error) {
	comment = new(model.Comment)
	if err = mysql.DB.Where("id=?", commentId).First(comment).Error; err != nil {
		return nil, err
	}
	return
}

// GetCommentByHeat 根据热度获一级评论
func GetCommentByHeat(userid int, articleId int, pageIndex int, pageSize int) (commentList []*model.Comment, err error) {
	rows, err := mysql.DB.Model(&model.Comment{}).Where("article_id=? and level=1", articleId).Order("give_like_count desc,id desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment = &model.Comment{}
		// ScanRows 方法用于将一行记录扫描至结构体
		mysql.DB.ScanRows(rows, comment)
		//业务逻辑
		userInfo, err := GetAUserInfoByUserId(comment.UserID)
		if err != nil {
			return nil, err
		}
		mysql.DB.Model(&model.Comment{}).Where("reply_comment_id =?", comment.ID).Count(&comment.SecondCommentsCount)
		comment.UserInfo = userInfo
		comment.IsGiveLike, _ = IsCommentGiveLike(userid, comment.ID)
		commentList = append(commentList, comment)
	}
	return commentList, nil
}

// GetReplyCommentByHeat 根据热度获二级评论
func GetReplyCommentByHeat(userid int, articleId int, commentId int, pageIndex int, pageSize int) (commentList []*model.Comment, err error) {
	rows, err := mysql.DB.Model(&model.Comment{}).Where("article_id=? and level=2 and comment_group=?", articleId, commentId).Order("give_like_count desc,id desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment = &model.Comment{}
		// ScanRows 方法用于将一行记录扫描至结构体
		mysql.DB.ScanRows(rows, comment)
		//业务逻辑
		replyUserName, err := GetAUserNameByUserId(comment.ReplyUserID)
		if err != nil {
			return nil, err
		}
		comment.ReplyUserName = replyUserName
		userInfo, err := GetAUserInfoByUserId(comment.UserID)
		if err != nil {
			return nil, err
		}
		comment.UserInfo = userInfo
		comment.IsGiveLike, _ = IsCommentGiveLike(userid, comment.ID)
		commentList = append(commentList, comment)
	}
	return commentList, nil
}
