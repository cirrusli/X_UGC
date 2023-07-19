package router

import (
	handler2 "X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"X_UGC/conf"
	"X_UGC/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func Register() *gin.Engine {
	if conf.C.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), mw.CORS())
	r.Static("/static/", "./view/static/")
	r.StaticFS("/upload", gin.Dir("./upload/", true))
	r.LoadHTMLGlob("./view/*.html")

	r.GET("/", Index)
	//发送邮箱验证码
	r.GET("/sendCode", handler2.SendCode)
	r.POST("/register", handler2.Register)
	r.POST("/login", handler2.Login)
	r.GET("/check", handler2.CheckAccount)
	r.POST("/checkCode", handler2.CheckCode)

	//websocket聊天
	r.GET("/webSocket", handler.WsHandler)
	return r
}
