package service

import (
	"X_UGC/biz/dal/mysql"
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"github.com/jinzhu/gorm"
	"strconv"
)

// AddFollowCount 初始化个人关注粉丝列
func AddFollowCount(userid int) (err error) {
	err = mysql.DB.Create(&model.FollowFansCount{
		UserID: userid,
	}).Error
	return
}

// IncrFollowCount 关注总数+1
func IncrFollowCount(userid string) (err error) {
	err = mysql.DB.Model(&model.FollowFansCount{}).Where("user_id=?", userid).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
	return
}

// IncrFansCount 粉丝总数+1
func IncrFansCount(userid string) (err error) {
	err = mysql.DB.Model(&model.FollowFansCount{}).Where("user_id=?", userid).Update("fans_count", gorm.Expr("fans_count + ?", 1)).Error
	return
}

// DecrFollowCount 关注总数-1
func DecrFollowCount(userid string) (err error) {
	err = mysql.DB.Model(&model.FollowFansCount{}).Where("user_id=?", userid).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
	return
}

// DecrFansCount 粉丝总数-1
func DecrFansCount(userid string) (err error) {
	err = mysql.DB.Model(&model.FollowFansCount{}).Where("user_id=?", userid).Update("fans_count", gorm.Expr("fans_count - ?", 1)).Error
	return
}

// GetFollowFansCount 获取关注和粉丝的总数
func GetFollowFansCount(userid string) (FollowFansCount *model.FollowFansCount, err error) {
	FollowFansCount = new(model.FollowFansCount)
	if err = mysql.DB.Where("user_id=?", userid).First(FollowFansCount).Error; err != nil {
		return nil, err
	}
	return
}

// GetRelationInfoByUserId  获取好友信息
func GetRelationInfoByUserId(key string, field string) (*model.RelationInfo, error) {
	val, err := redis.HGet(model.RELATION+key, field)
	if err != nil {
		return nil, err
	}
	if val == "" {
		val = "0"
	}
	status, _ := strconv.Atoi(val)
	userId, _ := strconv.Atoi(field)
	userInfo, _ := GetAUserInfoByUserId(userId)
	friendInfo := &model.RelationInfo{
		Status:   status,
		UserInfo: userInfo,
	}
	return friendInfo, nil
}

// GetRelationInfoStatus   获取和好友之间状态
func GetRelationInfoStatus(key string, field string) (int, error) {
	val, err := redis.HGet(model.RELATION+key, field)
	if err != nil {
		return -1, err
	}
	if val == "" {
		val = "0"
	}
	status, _ := strconv.Atoi(val)
	return status, nil
}

// ChangeRelationInfoStatus  改变和好友之间的状态状态
func ChangeRelationInfoStatus(key string, field string, status int) (err error) {
	switch status {
	case model.NoRelationship:
		if err = redis.HSet(model.RELATION+key, field, "1"); err != nil {
			return
		}
		if err = IncrFollowCount(key); err != nil {
			return
		}
		if err = redis.HSet(model.RELATION+field, key, "2"); err != nil {
			return
		}
		if err = IncrFansCount(field); err != nil {
			return
		}
	case model.Follow:
		if err = redis.HSet(model.RELATION+key, field, "0"); err != nil {
			return
		}
		if err = DecrFollowCount(key); err != nil {
			return
		}
		if err = redis.HSet(model.RELATION+field, key, "0"); err != nil {
			return
		}
		if err = DecrFansCount(field); err != nil {
			return
		}
	case model.Fans:
		if err = redis.HSet(model.RELATION+key, field, "3"); err != nil {
			return
		}
		if err = IncrFollowCount(key); err != nil {
			return
		}
		if err = redis.HSet(model.RELATION+field, key, "3"); err != nil {
			return
		}
		if err = IncrFansCount(field); err != nil {
			return
		}
	case model.FocusOnEachOther:
		if err = redis.HSet(model.RELATION+key, field, "2"); err != nil {
			return
		}
		if err = DecrFollowCount(key); err != nil {
			return
		}
		if err = redis.HSet(model.RELATION+field, key, "1"); err != nil {
			return
		}
		if err = DecrFansCount(field); err != nil {
			return
		}
	}
	return nil
}

// GetFollowList 	获取关注的人的列表
func GetFollowList(userid string) ([]*model.RelationInfo, error) {
	vals, err := redis.HGetAll(model.RELATION + userid)
	if err != nil {
		return nil, err
	}
	var FollowList []*model.RelationInfo
	for field, val := range vals {
		status, _ := strconv.Atoi(val)
		if status == model.Follow || status == model.FocusOnEachOther {
			userId, _ := strconv.Atoi(field)
			userInfo, _ := GetAUserInfoByUserId(userId)
			RelationInfo := &model.RelationInfo{
				Status:   status,
				UserInfo: userInfo,
			}
			FollowList = append(FollowList, RelationInfo)
		}
	}
	return FollowList, nil
}

// GetFansList 	获取个人粉丝的列表
func GetFansList(userid string) ([]*model.RelationInfo, error) {
	vals, err := redis.HGetAll(model.RELATION + userid)
	if err != nil {
		return nil, err
	}
	var FansList []*model.RelationInfo
	for field, val := range vals {
		status, _ := strconv.Atoi(val)
		if status == model.Fans || status == model.FocusOnEachOther {
			userid, _ := strconv.Atoi(field)
			userInfo, _ := GetAUserInfoByUserId(userid)
			RelationInfo := &model.RelationInfo{
				Status:   status,
				UserInfo: userInfo,
			}
			FansList = append(FansList, RelationInfo)
		}
	}
	return FansList, nil
}

// GetFocusOnEachOtherList     获取互相关注的列表、
func GetFocusOnEachOtherList(userid string) ([]*model.RelationInfo, error) {
	vals, err := redis.HGetAll(model.RELATION + userid)
	if err != nil {
		return nil, err
	}
	var FansList []*model.RelationInfo
	for field, val := range vals {
		status, _ := strconv.Atoi(val)
		if status == model.FocusOnEachOther {
			userId, _ := strconv.Atoi(field)
			userInfo, _ := GetAUserInfoByUserId(userId)
			RelationInfo := &model.RelationInfo{
				Status:   status,
				UserInfo: userInfo,
			}
			FansList = append(FansList, RelationInfo)
		}
	}
	return FansList, nil
}
