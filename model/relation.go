package model

const (
	NoRelationship   = iota //没关系
	Follow                  //是本人关注的人
	Fans                    //是本人的粉丝
	FocusOnEachOther        //互相关注
)

// RelationInfo 用户关系信息结构
type RelationInfo struct {
	Status   int       `json:"status"`
	UserInfo *UserInfo `json:"user_info"` //所查看用户的个人信息
}

// FollowFansCount 关注、粉丝总数
type FollowFansCount struct {
	UserID      int `json:"userid" gorm:"primary_key;unique;not null;auto_Increment:false;"` //外键绑定用户id
	FollowCount int `json:"follow_count" gorm:"default:0"`                                   //关注总数
	FansCount   int `json:"fans_count" gorm:"default:0"`                                     //粉丝总数
}
