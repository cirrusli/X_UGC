package router

import (
	"X_UGC/common/mw"
	"X_UGC/handler"
	"github.com/gin-gonic/gin"
)

func UserInfoGroups(r *gin.Engine) {
	UserInfoGroup := r.Group("/userInfo", mw.Auth())
	{
		//更新个人账号
		UserInfoGroup.PUT("/updateUserAccount", handler.UpdateUserAccount)
		//更新个人密码
		UserInfoGroup.PUT("/updateUserPassword", handler.UpdateUserPassword)
		//修改个人信息
		UserInfoGroup.PUT("/updateUserInfo", handler.UpdateUserInfo)
		//查询个人信息
		UserInfoGroup.GET("/getUserInfo", handler.GetAUserInfoByUserId)
		//上传用户头像
		UserInfoGroup.POST("/headPhotoUpload", handler.HeadPhotoUpload)
		//上传背景图片
		UserInfoGroup.POST("/backPhotoUpload", handler.BackPhotoUpload)
	}
}
