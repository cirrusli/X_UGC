package websocket

import "X_UGC/model"

//定义消息类型
const (
	Logout          = iota //注销
	ToOneFriend            //一对一
	ToAllFriend            //给所有好友和粉丝推送通知
	ToAllOnlineUser        //所有在线用户
	ToSomeUser             //管理员给一些用户通知
	FollowNotify           //关注通知
	GiveLikeNotify         //点赞通知
	CommentNotify          //评论通知
)

//SendMsg 发送消息的类型
type SendMsg struct {
	Type      int      `json:"type"`       //消息类型
	Content   string   `json:"content"`    //消息内容
	SendTime  string   `json:"send_time"`  //发送消息的时间
	ClientIDs []string `json:"client_ids"` //要发送用户的列表
}

//ReplyMsg 回复的消息
type ReplyMsg struct {
	MsgType  int    `json:"msg_type"`
	FromUser string `json:"from_user"` //发送消息用户的所有信息
	SendTime string `json:"send_time"` //消息发送而来的时间
	Content  string `json:"content"`   //消息内容
}

//MessageData 发送数据信息
type MessageData struct {
	FromUser string   `json:"from_user"`
	Message  *SendMsg `json:"message"`
}

//HistoryMsg 发送历史数据
type HistoryMsg struct {
	UserInfo    *model.UserInfo `json:"user_info"`
	HistoryData []string        `json:"history_data"`
}
