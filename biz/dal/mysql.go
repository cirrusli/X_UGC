package dal

import (
	"X_UGC/biz/model"
	"X_UGC/conf"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func InitMySQL(c *conf.MySQL) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.UserName, c.Password, c.Host, c.Port, c.DB)
	DB, err = gorm.Open("mysql", dsn)

	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func CloseMySQL() {
	err := DB.Close()
	if err != nil {
		return
	}
}

// InitTables 初始化表模型，包括外键等约束
func InitTables() {
	// 模型绑定
	DB.AutoMigrate(&model.User{}, &model.UserInfo{}, &model.ArticleTypeDict{}, &model.ArticleInfo{}, &model.FollowFansCount{}, &model.Weight{}, &model.Comment{})
	// 外键添加
	DB.Model(&model.UserInfo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade")
	DB.Model(&model.ArticleInfo{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade")
	DB.Model(&model.FollowFansCount{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade")
	DB.Model(&model.ArticleInfo{}).AddForeignKey("article_type_id", "article_type_dicts(id)", "cascade", "cascade")
	DB.Model(&model.Weight{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade")
	DB.Model(&model.Weight{}).AddForeignKey("article_type_id", "article_type_dicts(id)", "cascade", "cascade")
	DB.Model(&model.Comment{}).AddForeignKey("article_id", "article_infos(id)", "cascade", "cascade")
	DB.Model(&model.Comment{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade")
}
