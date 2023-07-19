package model

// Weight 权重表
type Weight struct {
	UserID        int `json:"userid" gorm:"primary_key;not null;auto_Increment:false;"`          //外键绑定用户id
	ArticleTypeID int `json:"article_type_id" gorm:"primary_key;not null;auto_Increment:false;"` //外键绑定文章类型id
	Weight        int `json:"weight" gorm:"not null"`
}
