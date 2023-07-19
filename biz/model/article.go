package model

// ArticleInfo 用户上传文章结构
type ArticleInfo struct {
	ID            int       `json:"article_id"`                       //文章ID
	AuthorID      int       `json:"author_id" gorm:"not null"`        //作者ID
	AuthorInfo    *UserInfo `json:"author_info" gorm:"-"`             //作者信息
	ReleaseTime   string    `json:"release_time" gorm:"not null"`     //发布时间
	Title         string    `json:"title" gorm:"not null;index"`      //文章标题
	Content       string    `json:"content" gorm:"not null"`          //文章内容
	ArticleTypeID int       `json:"article_type_id" gorm:"not null"`  //文章类型id,外键绑定类型字典
	ArticleType   string    `json:"article_type" gorm:"-"`            //文章类型
	IsVideo       int       `json:"is_video" gorm:"not null"`         //0不是视频类型，1是视频类型
	ResourceDir   string    `json:"resource_dir"`                     //文章资源目录存储路径
	ResourceUrl   []string  `json:"resource_url"  gorm:"-"`           //文章资源路径
	CoverUrl      string    `json:"cover_url"`                        //封面地址
	GiveLikeCount int       `json:"give_like_count" gorm:"default:0"` //点赞总数
	CommentCount  int       `json:"comment_count" gorm:"default:0"`   //评论总数
}

// ArticleTypeDict 文章类型字典
//
//	1    美食
//	2	 旅行
//	3    新闻
//	4	 教育
//	5	 科普
//	6	 影视
type ArticleTypeDict struct {
	ID          int    `json:"article_type_id"`
	ArticleType string `json:"article_type" gorm:"not null;"`
}
