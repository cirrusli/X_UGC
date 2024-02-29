package handler

import (
	"X_UGC/pkg/common/jwt"
	ws "X_UGC/pkg/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

// WsHandler 客户端连接
func WsHandler(c *gin.Context) {
	upGrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		},
		Subprotocols: []string{c.Request.Header.Get("Sec-WebSocket-Protocol")},
	}

	// 解析JWT令牌
	mc, err := jwt.ParseToken(upGrade.Subprotocols[0])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userId := strconv.Itoa(mc.UserID)

	// 创建WebSocket连接
	conn, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &ws.Client{
		Id:      userId,
		Socket:  conn,
		Message: make(chan []byte, 1024),
	}
	client.PullHistoryInformation()

	//注册
	ws.WsManager.RegisterClient(client)
	//起协程，实时接收和回复数据
	go client.Read()
	go client.Write()
}
