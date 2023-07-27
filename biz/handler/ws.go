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
	upGrande := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		},
		Subprotocols: []string{c.Request.Header.Get("Sec-WebSocket-Protocol")},
	}

	mc, _ := jwt.ParseToken(upGrande.Subprotocols[0])
	userId := strconv.Itoa(mc.UserID)
	//创建连接
	conn, err := upGrande.Upgrade(c.Writer, c.Request, nil)
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
