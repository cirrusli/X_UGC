package handler

import (
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"X_UGC/biz/service"
	"X_UGC/pkg/common/email"
	"X_UGC/pkg/common/jwt"
	"X_UGC/pkg/common/pwd"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

func SendCode(c *gin.Context) {
	e := c.Query("email")
	code, err := email.SendCode(e)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = redis.Set(e, code, 60)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

// CheckCode 判断验证码是否正确
func CheckCode(c *gin.Context) {
	jsonData := make(map[string]string)
	c.BindJSON(&jsonData)
	val, err := redis.Get(jsonData["email"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if val == jsonData["code"] {
		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"check_result": true,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"check_result": false,
	})
}

// Register 注册
func Register(c *gin.Context) {
	// 前端页面填写用户信息 点击提交 会发请求到这里
	// 1. 从请求中把数据拿出来
	jsonData := map[string]string{}
	c.BindJSON(&jsonData)
	var user = model.User{
		Number:   jsonData["number"],
		Email:    jsonData["email"],
		Password: pwd.HashPassword(jsonData["password"]),
	}
	//2.判定邮箱验证码
	val, err := redis.Get(user.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if val == jsonData["code"] {
		// 3. 存入数据库
		err = service.AddUser(&user)
		var userInfo = model.UserInfo{
			UserName: "用户" + user.Number[5:],
			UserID:   user.ID,
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		if err = service.AddUserInfo(&userInfo); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		if err = service.AddFollowCount(user.ID); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		if err = service.AddWeightInfo(user.ID); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})

	} else {
		c.JSON(http.StatusOK, "code false")
	}
}

// Login 登录
func Login(c *gin.Context) {
	account := c.PostForm("account")
	p := c.PostForm("password")
	user, err := service.CheckAccount(account)
	isPwd := pwd.CheckPasswordHash(user.Password, p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		if isPwd {
			tokenString, _ := jwt.GenToken(user.ID)
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"token":  tokenString,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
			})
		}
	}
}

// CheckAccount 检查账号
func CheckAccount(c *gin.Context) {
	account := c.Query("account")
	user, err := service.CheckAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		if user.ID > 0 {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
			})
		}
	}
}

// GetAllUser 查询所有用户
func GetAllUser(c *gin.Context) {
	// 查询users这个表里的所有数据
	userList, err := service.GetAllUser()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   userList,
		})
	}
}

// GetAUserById 根据id获取用户
func GetAUserById(c *gin.Context) {
	userid := c.Query("userid")

	user, err := service.GetAUserById(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   user,
		})
	}
}

// UpdateUserAccount 更新账号
func UpdateUserAccount(c *gin.Context) {
	userid := c.GetInt("userid")
	jsonData := make(map[string]string)
	c.BindJSON(&jsonData)

	//2.判定邮箱验证码
	val, err := redis.Get(jsonData["email"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if val == jsonData["code"] {
		if err = service.UpdateUserAccount(userid, jsonData["account_type"], jsonData["account"]); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, "code false")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// UpdateUserPassword 更新密码
func UpdateUserPassword(c *gin.Context) {
	userid := c.GetInt("userid")
	jsonData := make(map[string]string)
	_ = c.BindJSON(&jsonData)
	user, err := service.GetAUserById(strconv.Itoa(userid))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//2.判定邮箱验证码
	val, err := redis.Get(user.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if val == jsonData["code"] {
		if pwd.CheckPasswordHash(user.Password, jsonData["old_pwd"]) {
			err := service.UpdateUserPassword(userid, pwd.HashPassword(jsonData["new_pwd"]))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusOK, "old_pwd false")
			return
		}
	} else {
		c.JSON(http.StatusOK, "code false")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DeleteAUser 删除用户
func DeleteAUser(c *gin.Context) {
	userid := c.Query("userid")
	if err := service.DeleteAUser(userid); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userid := c.GetInt("userid")
	userInfo, err := service.GetAUserInfoByUserId(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.BindJSON(userInfo)
	if err = service.UpdateUserInfo(userInfo); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

// GetAUserInfoByUserId 获取用户信息
func GetAUserInfoByUserId(c *gin.Context) {
	userid := c.GetInt("userid")
	userInfo, err := service.GetAUserInfoByUserId(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	user, err := service.GetAUserById(strconv.Itoa(userid))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   map[string]interface{}{"user": map[string]string{"number": user.Number, "email": user.Email}, "userinfo": userInfo},
	})

}

// GetAllUserInfo 获取所有用户信息
func GetAllUserInfo(c *gin.Context) {
	// 查询user_infos这个表里的所有数据
	userInfoList, err := service.GetAllUserInfo()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   userInfoList,
		})
	}
}

// HeadPhotoUpload 头像上传
func HeadPhotoUpload(c *gin.Context) {
	userid := c.GetInt("userid")
	//获取表单file的数据
	fileHeader, err := c.FormFile("head_photo_file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	fileExt := filepath.Ext(fileHeader.Filename)
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".jpeg" || fileExt == ".webp" || fileExt == ".bmp" {
		fileExt = ".jpeg"
	} else {
		c.JSON(http.StatusOK, gin.H{"error": "图片格式错误，请重新上传！"})
		return
	}
	fileDir := "./upload/head_photo/"
	filePath := fileDir + "head_photo_user_" + strconv.Itoa(userid) + fileExt
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		//持久化头像路径
		err = service.UpdateHeadPhoto(filePath, userid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"head_photo": filePath,
		})
	}
}

// BackPhotoUpload 背景图片上传
func BackPhotoUpload(c *gin.Context) {
	userid := c.GetInt("userid")
	//获取表单file的数据
	fileHeader, err := c.FormFile("back_photo_file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	fileExt := filepath.Ext(fileHeader.Filename)
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".jpeg" || fileExt == ".webp" || fileExt == ".bmp" {
		fileExt = ".jpeg"
	} else {
		c.JSON(http.StatusOK, gin.H{"error": "图片格式错误，请重新上传！"})
		return
	}
	fileDir := "./upload/back_photo/"
	filePath := fileDir + "back_photo_user_" + strconv.Itoa(userid) + fileExt
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		//持久化背景图片路径
		err = service.UpdateBackPhoto(filePath, userid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"back_photo": filePath,
		})
	}

}
