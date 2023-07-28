package router

import (
	"X_UGC/biz/handler"
	"X_UGC/biz/mw"
	"github.com/gin-gonic/gin"
)

func ManageGroups(r *gin.Engine) {
	ManageGroup := r.Group("/manage", mw.Auth())
	{
		// 查看所有用户
		ManageGroup.GET("/getAllUser", handler.GetAllUser)
		//查询一个用户
		ManageGroup.GET("/getUser", handler.GetAUserById)
		// 修改某一个用户
		//ManageGroup.PUT("/updateUserAccount", handler.UpdateAUserAccount)
		// 删除某一个用户
		ManageGroup.DELETE("/deleteUser", handler.DeleteAUser)
		//查询所有用户个人信息
		ManageGroup.GET("/getAllUserInfo", handler.GetAllUserInfo)
	}

}
