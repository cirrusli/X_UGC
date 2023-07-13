package model

// User 用户结构体
type User struct {
	ID       int    `json:"userid"`                        //用户id
	Number   string `json:"number" gorm:"unique;not null"` //电话号
	Email    string `json:"email" gorm:"unique;not null"`  //电子邮箱
	Password string `json:"password"`                      //密码
}

// UserInfo 个人信息
type UserInfo struct {
	UserID    int    `json:"userid" gorm:"primary_key;unique;not null;auto_Increment:false;"`   //外键绑定用户id
	UserName  string `json:"username" gorm:"unique;"`                                           //用户名
	Sex       string `json:"sex"`                                                               //性别
	Birthday  string `json:"birthday"`                                                          //生日
	Place     string `json:"place"`                                                             //住址
	Info      string `json:"info"`                                                              //简介
	HeadPhoto string `json:"head_photo" gorm:"default:'./upload/head_photo/default_head.jpg';"` //头像
	BackPhoto string `json:"back_photo" gorm:"default:'./upload/back_photo/default_back.jpg';"` //背景图片

	FollowCount int `json:"follow_count" gorm:"-"` //关注总数
	FansCount   int `json:"fans_count" gorm:"-"`   //粉丝总数
}
