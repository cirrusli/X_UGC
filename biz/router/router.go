package router

import (
	"X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"X_UGC/conf"
	"github.com/gin-contrib/gzip"
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
	r.Use(gin.Logger(), gin.Recovery(), gzip.Gzip(gzip.DefaultCompression), mw.CORS())
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

	//用户信息路由组
	UserInfoGroups(r)

	//注册用户关系路由组
	RelationInfoGroups(r)

	//注册管理员操作路由组
	ManageGroups(r)

	//注册文章路由组
	ArticlesGroups(r)

	//注册评论路由组
	CommentGroups(r)

	//注册通知路由组
	NotifyGroups(r)

	//注册搜索路由组
	SearchGroups(r)
	return r
}
