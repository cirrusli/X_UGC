package router

import (
	"X_UGC/common/mw"
	"X_UGC/conf"
	"X_UGC/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func Init() *gin.Engine {
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
	r.GET("/sendCode", controller.SendCode)
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.GET("/check", controller.CheckAccount)
	r.POST("/checkCode", controller.CheckCode)

	//websocket聊天
	r.GET("/webSocket", controller.WsHandler)
	return r
}
