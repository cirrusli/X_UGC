package router

import (
	"X_UGC/conf"
	"X_UGC/handler"
	"X_UGC/mw"
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
	r.GET("/sendCode", handler.SendCode)
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.GET("/check", handler.CheckAccount)
	r.POST("/checkCode", handler.CheckCode)

	//websocket聊天
	r.GET("/webSocket", handler.WsHandler)
	return r
}
