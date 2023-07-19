package service

import (
	"X_UGC/biz/dal"
	"X_UGC/biz/model"
)

// AddUser 添加一个用户
func AddUser(user *model.User) (err error) {
	err = dal.DB.Create(user).Error
	return
}

// GetAllUser 获取所有用户
func GetAllUser() (userList []*model.User, err error) {
	if err = dal.DB.Find(&userList).Error; err != nil {
		return nil, err
	}
	return
}

// GetAUserById 获取一个用户
func GetAUserById(id string) (user *model.User, err error) {
	user = new(model.User)
	if err = dal.DB.Where("id=?", id).First(user).Error; err != nil {
		return nil, err
	}
	return
}

// CheckAccount 账号检查
func CheckAccount(account string) (user *model.User, err error) {
	user = new(model.User)
	if err = dal.DB.Where("number=? or email=?", account, account).First(user).Error; err != nil {
		return nil, err
	}
	return
}

// CheckAccountAndPwd 账号和密码检查
func CheckAccountAndPwd(account string, pwd string) (user *model.User, err error) {
	user = new(model.User)
	if err = dal.DB.Where("(number=? or email=?) and password=?", account, account, pwd).First(user).Error; err != nil {
		return nil, err
	}
	return
}

// UpdateAUser 更新用户
func UpdateAUser(user *model.User) (err error) {
	err = dal.DB.Save(user).Error
	return
}

// UpdateUserAccount 更新用户账号
func UpdateUserAccount(userid int, accountType string, account string) (err error) {
	err = dal.DB.Model(&model.User{}).Where("id = ?", userid).Update(accountType, account).Error
	return
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userid int, newPassword string) (err error) {
	err = dal.DB.Model(&model.User{}).Where("id = ?", userid).Update("password", newPassword).Error
	return
}

// DeleteAUser 删除一个用户
func DeleteAUser(id string) (err error) {
	err = dal.DB.Where("id=?", id).Delete(&model.User{}).Error
	return
}

// AddUserInfo 创建个人信息
func AddUserInfo(userInfo *model.UserInfo) (err error) {
	err = dal.DB.Create(userInfo).Error
	return
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(userInfo *model.UserInfo) (err error) {
	err = dal.DB.Save(userInfo).Error
	return
}

// GetAllUserInfo 获取所有用户信息
func GetAllUserInfo() (userInfoList []*model.UserInfo, err error) {
	if err = dal.DB.Select("*").
		Joins("inner join follow_fans_counts on user_infos.user_id=follow_fans_counts.user_id").
		Find(&userInfoList).Error; err != nil {
		return nil, err
	}
	return
}

// GetAUserInfoByUserId 根据用户id获取一个用户信息
func GetAUserInfoByUserId(userid int) (userInfo *model.UserInfo, err error) {
	userInfo = new(model.UserInfo)
	if err = dal.DB.Select("*").Where("user_infos.user_id=?", userid).
		Joins("inner join follow_fans_counts on user_infos.user_id=follow_fans_counts.user_id").
		First(userInfo).Error; err != nil {
		return nil, err
	}
	return
}

// GetAUserNameByUserId 根据用户id获取一个用户名
func GetAUserNameByUserId(userid int) (userName string, err error) {
	var userinfo = model.UserInfo{}
	if err = dal.DB.Select("user_name").Where("user_id=?", userid).First(&userinfo).Error; err != nil {
		return "", err
	}
	userName = userinfo.UserName
	return
}

// UpdateHeadPhoto  更新头像信息
func UpdateHeadPhoto(filePath string, id int) (err error) {
	err = dal.DB.Model(&model.UserInfo{}).Where("user_id=?", id).Update("head_photo", filePath).Error
	return
}

// UpdateBackPhoto  更新背景图片信息
func UpdateBackPhoto(filePath string, id int) (err error) {
	err = dal.DB.Model(&model.UserInfo{}).Where("user_id=?", id).Update("back_photo", filePath).Error
	return
}
