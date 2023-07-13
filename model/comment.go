package model

type Comment struct {
	ID             int    `json:"comment_id"`
	ArticleID      int    `json:"article_id" gorm:"not null"`          //外键绑定文章id
	UserID         int    `json:"userid" gorm:"not null"`              //外键绑定用户id
	Level          int    `json:"level" gorm:"not null"`               //评论级别  1级评论   2级评论
	Content        string `json:"content" gorm:"not null"`             //评论内容
	CommentTime    string `json:"comment_time" gorm:"not null"`        //评论时间
	CommentGroup   int    `json:"comment_group" gorm:"not null;index"` //评论的组别,一级评论是所在文章的id，所属二级评论组别为一级评论id
	ReplyCommentID int    `json:"reply_comment_id"`                    //如果是二级评论，返回其指向的一级或二级评论id
	ReplyUserID    int    `json:"reply_user_id"`                       //回复的用户的id
	GiveLikeCount  int    `json:"give_like_count" gorm:"default:0"`    //点赞总数

	UserInfo            *UserInfo `json:"user_info" gorm:"-"`             //评论人的信息
	ReplyUserName       string    `json:"reply_user_name" gorm:"-"`       //回复的人的姓名
	SecondCommentsCount int       `json:"second_comments_count" gorm:"-"` //如果为一级，该一级评论下二级评论的总数
	IsGiveLike          bool      `json:"is_give_like" gorm:"-"`          //评论是否点赞
}
