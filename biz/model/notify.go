package model

// NotifyInfo 通知信息的基类
type NotifyInfo struct {
	SendTime string    `json:"send_time"`  //发送时间
	FromUser *UserInfo `json:"from_user"`  //来自谁
	ToUserID int       `json:"to_user_id"` //给谁发
}

// FollowNotifyInfo  关注通知
type FollowNotifyInfo struct {
	NotifyInfo
}

// GiveLikeNotifyInfo 点赞通知
type GiveLikeNotifyInfo struct {
	NotifyInfo               //继承基类
	ArticleInfo *ArticleInfo `json:"article_info"` //文章信息
}

// CommentNotifyInfo 评论通知
type CommentNotifyInfo struct {
	NotifyInfo
	ArticleInfo    *ArticleInfo `json:"article_info"`     //文章信息
	Content        string       `json:"content"`          //评论内容
	ReplyCommentID int          `json:"reply_comment_id"` //回复评论的id
}

// NotifyHistory 历史通知
type NotifyHistory struct {
	SendTime string                 `json:"send_time"`
	Notify   map[string]interface{} `json:"notify"`
}

// UnreadChat 未读聊天信息通知
type UnreadChat struct {
	SendTime   string `json:"send_time"`
	UnreadNum  string `json:"unread_num"`
	ChatUserID string `json:"chat_userid"`
}
