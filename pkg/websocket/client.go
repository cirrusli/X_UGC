package websocket

import (
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"X_UGC/biz/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Client 单个websocket
type Client struct {
	Id      string
	Socket  *websocket.Conn
	Message chan []byte
}

// 从websocket中直接读取数据
func (c *Client) Read() {
	defer func() {
		//客户端关闭
		if err := c.Socket.Close(); err != nil {
			WsManager.UnRegisterClient(c)
			return
		}
		//关闭后直接注销客户端
		WsManager.UnRegisterClient(c)
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		//读取数据失败
		if err != nil || messageType == websocket.CloseMessage {
			WsManager.UnRegisterClient(c)
			_ = c.Socket.Close()
			return
		}

		fmt.Println(string(message))
		//解析发送过来的参数
		var data SendMsg

		err = json.Unmarshal(message, &data)
		if err != nil {
			continue
		}

		if data.Type == Logout {
			//关闭后直接注销客户端
			if err := c.Socket.Close(); err != nil {
				WsManager.UnRegisterClient(c)
				return
			}
			WsManager.UnRegisterClient(c)
			return
		}

		//消息处理
		c.ProcessData(&data)

	}
}

func (c *Client) ProcessData(data *SendMsg) {

	//类型为1，好友私聊
	if data.Type == ToOneFriend {
		WsManager.Message <- &MessageData{
			FromUser: c.Id,
			Message:  data,
		}
		//保存到redis数据库
		//自己保存发送的消息
		var sendBuilder strings.Builder
		sendBuilder.WriteString("send+")
		sendBuilder.WriteString(data.ClientIDs[0])
		sendBuilder.WriteString("+")
		sendBuilder.WriteString(data.SendTime)
		err := redis.HSet(model.MESSAGE_HISTORY+c.Id, sendBuilder.String(), data.Content)
		if err != nil {
			_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
		}
		err = service.ResetUnreadChat(c.Id, data.ClientIDs[0], data.SendTime)
		if err != nil {
			_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
		}
		//接收者保存接收到的消息
		var receiveBuilder strings.Builder
		receiveBuilder.WriteString("receive+")
		receiveBuilder.WriteString(c.Id)
		receiveBuilder.WriteString("+")
		receiveBuilder.WriteString(data.SendTime)
		err = redis.HSet(model.MESSAGE_HISTORY+data.ClientIDs[0], receiveBuilder.String(), data.Content)
		if err != nil {
			_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
		}
		res, _ := redis.HGet(model.UNREAD_CHAT+data.ClientIDs[0], c.Id)
		if res == "" {
			res = data.SendTime + "+0"
		}
		unreadNum, _ := strconv.Atoi(strings.Split(res, "+")[1])
		err = service.IncrUnreadChat(data.ClientIDs[0], c.Id, data.SendTime, unreadNum+1)
		if err != nil {
			_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
		}

	} else if data.Type == ToAllFriend { //类型为2，所有好友和粉丝推送
		//获取好友列表
		FansList, _ := service.GetFansList("relation_" + c.Id)
		WsManager.Message <- &MessageData{
			FromUser: c.Id,
			Message:  data,
		}
		for _, val := range FansList {
			//保存到redis数据库
			//自己保存发送的消息
			var sendBuilder strings.Builder
			sendBuilder.WriteString("send+")
			toSendID := strconv.Itoa(val.UserInfo.UserID)
			sendBuilder.WriteString(toSendID)
			sendBuilder.WriteString("+")
			sendBuilder.WriteString(data.SendTime)

			err := redis.HSet(model.MESSAGE_HISTORY+c.Id, sendBuilder.String(), data.Content)
			if err != nil {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
			}
			err = service.ResetUnreadChat(c.Id, data.ClientIDs[0], data.SendTime)
			if err != nil {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
			}
			//接收者保存接收到的消息
			var receiveBuilder strings.Builder
			receiveBuilder.Reset()
			receiveBuilder.WriteString("receive+")
			receiveBuilder.WriteString(c.Id)
			receiveBuilder.WriteString("+")
			receiveBuilder.WriteString(data.SendTime)

			err = redis.HSet(model.MESSAGE_HISTORY+toSendID, receiveBuilder.String(), data.Content)
			if err != nil {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
			}
			res, _ := redis.HGet(model.UNREAD_CHAT+data.ClientIDs[0], c.Id)
			if res == "" {
				res = data.SendTime + "+0"
			}
			unreadNum, _ := strconv.Atoi(strings.Split(res, "+")[1])
			err = service.IncrUnreadChat(data.ClientIDs[0], c.Id, data.SendTime, unreadNum+1)
			if err != nil {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
			}
		}

	} else if data.Type == ToAllOnlineUser { //类型为3，该平台所有在线用户
		WsManager.Message <- &MessageData{
			FromUser: c.Id,
			Message:  data,
		}

	} else if data.Type == ToSomeUser {
		WsManager.Message <- &MessageData{
			FromUser: c.Id,
			Message:  data,
		}
		//保存到redis数据库
		var MsgBuilder strings.Builder
		MsgBuilder.WriteString(c.Id)
		MsgBuilder.WriteString("+")
		MsgBuilder.WriteString(data.SendTime)
		var err error
		for _, v := range data.ClientIDs {
			err = redis.HSet(model.MESSAGE_HISTORY+v, MsgBuilder.String(), data.Content)
			if err != nil {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("消息发送失败!"))
			}
		}
	}
}

// 写入数据到websocket中
func (c *Client) Write() {
	ticker := time.NewTicker(time.Second * 30)
	defer func() {
		ticker.Stop()
		//客户端关闭
		if err := c.Socket.Close(); err != nil {
			WsManager.UnRegisterClient(c)
			return
		}
		//关闭后直接注销客户端
		WsManager.UnRegisterClient(c)
	}()
	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				//关闭后直接注销客户端
				if err := c.Socket.Close(); err != nil {
					WsManager.UnRegisterClient(c)
					return
				}
				WsManager.UnRegisterClient(c)
				//消息通道关闭后直接注销客户端
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				//关闭后直接注销客户端
				if err := c.Socket.Close(); err != nil {
					WsManager.UnRegisterClient(c)
					return
				}
				WsManager.UnRegisterClient(c)
				return
			}
		// 服务端心跳:每20秒ping一次客户端，查看其是否在线
		case <-ticker.C:
			_ = c.Socket.SetWriteDeadline(time.Now().Add(time.Second * 35))

			if err := c.Socket.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {

				//关闭后直接注销客户端
				if err = c.Socket.Close(); err != nil {
					WsManager.UnRegisterClient(c)
					return
				}
				WsManager.UnRegisterClient(c)
				return
			}
		}
	}
}

// PullHistoryInformation 拉取历史信息
func (c *Client) PullHistoryInformation() {
	HistoryInformation, err := redis.HGetAll(model.MESSAGE_HISTORY + c.Id)
	if err != nil {
		_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("历史消息获取失败!"))
	}
	var HistoryData = make(map[string][]string)
	for key, val := range HistoryInformation {
		arr := strings.Split(key, "+")
		status := arr[0]
		UserId := arr[1]
		sendTime := arr[2]
		HistoryData[UserId] = append(HistoryData[UserId], sendTime+"+"+status+"+"+val)
	}

	replyHistoryData := make(map[string]*HistoryMsg)
	for key := range HistoryData {
		sort.Strings(HistoryData[key])
		id, _ := strconv.Atoi(key)
		UserInfo, _ := service.GetAUserInfoByUserId(id)
		replyHistoryData[key] = &HistoryMsg{
			UserInfo:    UserInfo,
			HistoryData: HistoryData[key],
		}
	}

	byteReplyHistoryData, _ := json.Marshal(struct {
		Type string                 `json:"type"`
		Data map[string]*HistoryMsg `json:"data"`
	}{"history", replyHistoryData})
	err = c.Socket.WriteMessage(websocket.TextMessage, byteReplyHistoryData)
	if err != nil {
		//关闭后直接注销客户端
		if err := c.Socket.Close(); err != nil {
			WsManager.UnRegisterClient(c)
			return
		}
		WsManager.UnRegisterClient(c)
		return
	}
}

// ProcessNotify 处理通知
func ProcessNotify(notifyType int, notifyInfo []byte) {
	switch notifyType {
	case FollowNotify:
		var followNotifyInfo model.FollowNotifyInfo
		_ = json.Unmarshal(notifyInfo, &followNotifyInfo)
		var content = make(map[string]interface{})
		content["form_user_info"] = followNotifyInfo.FromUser
		contentJson, _ := json.Marshal(content)
		var data = &SendMsg{
			Type:      FollowNotify,
			Content:   fmt.Sprintf("%s", contentJson),
			SendTime:  followNotifyInfo.SendTime,
			ClientIDs: []string{strconv.Itoa(followNotifyInfo.ToUserID)},
		}
		//实时发送通知信息
		WsManager.Message <- &MessageData{
			FromUser: strconv.Itoa(followNotifyInfo.FromUser.UserID),
			Message:  data,
		}
		score, _ := strconv.ParseFloat(data.SendTime, 64)
		//持久化通知信息
		err := redis.ZAdd(model.NOTIFY_HISTORY+"follow:"+strconv.Itoa(followNotifyInfo.ToUserID), score, data.Content)
		if err != nil {
			log.Println("关注通知持久化失败!", err)
		}
		err = redis.HIncrBy(model.UNREAD_NOTIFY+"follow", strconv.Itoa(followNotifyInfo.ToUserID), 1)
		if err != nil {
			log.Println("关注未读通知更新失败", err)
		}
	case GiveLikeNotify:
		var giveLikeNotifyInfo model.GiveLikeNotifyInfo
		json.Unmarshal(notifyInfo, &giveLikeNotifyInfo)
		var content = make(map[string]interface{})
		content["form_user_info"] = giveLikeNotifyInfo.FromUser
		content["article_info"] = giveLikeNotifyInfo.ArticleInfo
		contentJson, _ := json.Marshal(content)
		var data = &SendMsg{
			Type:      GiveLikeNotify,
			Content:   fmt.Sprintf("%s", contentJson),
			SendTime:  giveLikeNotifyInfo.SendTime,
			ClientIDs: []string{strconv.Itoa(giveLikeNotifyInfo.ToUserID)},
		}
		//实时发送通知信息
		WsManager.Message <- &MessageData{
			FromUser: strconv.Itoa(giveLikeNotifyInfo.FromUser.UserID),
			Message:  data,
		}
		score, _ := strconv.ParseFloat(data.SendTime, 64)
		//持久化通知信息
		err := redis.ZAdd(model.NOTIFY_HISTORY+"giveLike:"+strconv.Itoa(giveLikeNotifyInfo.ToUserID), score, data.Content)
		if err != nil {
			log.Println("点赞通知持久化失败!", err)
		}
		err = redis.HIncrBy(model.UNREAD_NOTIFY+"giveLike", strconv.Itoa(giveLikeNotifyInfo.ToUserID), 1)
		if err != nil {
			log.Println("点赞未读通知更新失败", err)
		}
	case CommentNotify:
		var commentNotifyInfo model.CommentNotifyInfo
		_ = json.Unmarshal(notifyInfo, &commentNotifyInfo)
		var content = make(map[string]interface{})
		content["form_user_info"] = commentNotifyInfo.FromUser
		content["article_info"] = commentNotifyInfo.ArticleInfo
		content["content"] = commentNotifyInfo.Content
		content["reply_comment_id"] = commentNotifyInfo.ReplyCommentID
		contentJson, _ := json.Marshal(content)
		var data = &SendMsg{
			Type:      CommentNotify,
			Content:   fmt.Sprintf("%s", contentJson),
			SendTime:  commentNotifyInfo.SendTime,
			ClientIDs: []string{strconv.Itoa(commentNotifyInfo.ToUserID)},
		}
		//实时发送通知信息
		WsManager.Message <- &MessageData{
			FromUser: strconv.Itoa(commentNotifyInfo.FromUser.UserID),
			Message:  data,
		}
		score, _ := strconv.ParseFloat(data.SendTime, 64)
		//持久化通知信息
		err := redis.ZAdd(model.NOTIFY_HISTORY+"comment:"+strconv.Itoa(commentNotifyInfo.ToUserID), score, data.Content)
		if err != nil {
			log.Println("评论通知持久化失败!", err)
		}
		err = redis.HIncrBy(model.UNREAD_NOTIFY+"comment", strconv.Itoa(commentNotifyInfo.ToUserID), 1)
		if err != nil {
			log.Println("评论未读通知更新失败", err)
		}
	default:
		log.Println("未知消息无法处理")
	}
}
