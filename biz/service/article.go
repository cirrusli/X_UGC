package service

import (
	"X_UGC/biz/dal/mysql"
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// GetArticleTypeById 根据类型id获取文章类型
func GetArticleTypeById(id int) (articleType *model.ArticleTypeDict) {
	articleType = new(model.ArticleTypeDict)
	if err := mysql.DB.Where("id = ?", id).First(articleType).Error; err != nil {
		return nil
	}
	return
}

// GetAllArticleType  获取所有文章类型
func GetAllArticleType() (articleTypeList []*model.ArticleTypeDict, err error) {
	if err = mysql.DB.Find(&articleTypeList).Error; err != nil {
		return nil, err
	}
	return
}

// UploadPhotoArticleResource 上传图片类型文章资源
func UploadPhotoArticleResource(c *gin.Context) (coverFilePath string, articleFileDir string, err error) {
	userid := c.GetInt("userid")
	//当前时间戳作为文件夹命名的一部分
	unixNano := time.Now().UnixNano()
	//获取封面文件
	coverFile, _ := c.FormFile("cover_file")
	//获取文章资源文件
	articleFiles := c.Request.MultipartForm.File["article_file"]
	if coverFile != nil {
		coverFileDir := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/photo_" + strconv.FormatInt(unixNano, 10) + "/cover/"
		_ = os.MkdirAll(coverFileDir, os.ModePerm)
		fileExt := filepath.Ext(coverFile.Filename)
		coverFilePath = coverFileDir + "cover" + fileExt
		if err = c.SaveUploadedFile(coverFile, coverFilePath); err != nil {
			return "", "", err
		}
	} else if coverFile == nil {
		coverFile = articleFiles[0]
		coverFileDir := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/photo_" + strconv.FormatInt(unixNano, 10) + "/cover/"
		os.MkdirAll(coverFileDir, os.ModePerm)
		fileExt := filepath.Ext(coverFile.Filename)
		coverFilePath = coverFileDir + "cover" + fileExt
		if err = c.SaveUploadedFile(coverFile, coverFilePath); err != nil {
			return "", "", err
		}
	}

	if articleFiles != nil {
		articleFileDir = "./upload/article_resource/user_" + strconv.Itoa(userid) + "/photo_" + strconv.FormatInt(unixNano, 10) + "/article/"
		os.MkdirAll(articleFileDir, os.ModePerm)
		for i, articleFile := range articleFiles {
			fileExt := filepath.Ext(articleFile.Filename)
			filePath := articleFileDir + "article_" + strconv.Itoa(i) + fileExt
			if err := c.SaveUploadedFile(articleFile, filePath); err != nil {
				return "", "", err
			}
		}
	}

	return coverFilePath, articleFileDir, nil
}

// AddArticle 创建一篇文章
func AddArticle(articleInfo *model.ArticleInfo) (err error) {
	err = mysql.DB.Create(articleInfo).Error
	return
}

// InsertToArticlePool 向不同类型文章池内插入文章id
func InsertToArticlePool(typeID int, articleId int) (err error) {
	err = redis.SAdd(model.ARTICLE_POOL+strconv.Itoa(typeID), strconv.Itoa(articleId))
	return
}

// InsertToFriendArticleList 向朋友的文章池里添加文章
func InsertToFriendArticleList(userid int, articleId int) error {
	userList, err := GetFocusOnEachOtherList(strconv.Itoa(userid))
	if err != nil {
		return err
	}
	for _, friendInfo := range userList {
		err = redis.LPush(model.FRIEND_ARTICLE_LIST+strconv.Itoa(friendInfo.UserInfo.UserID), articleId)
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertToFansArticleList 向粉丝的文章池里添加文章
func InsertToFansArticleList(userid int, articleId int) error {
	userList, err := GetFansList(strconv.Itoa(userid))
	if err != nil {
		return err
	}
	for _, fansInfo := range userList {
		err = redis.LPush(model.FOLLOW_ARTICLE_LIST+strconv.Itoa(fansInfo.UserInfo.UserID), articleId)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetArticleFromFriend 分页获取朋友的文章
func GetArticleFromFriend(userid int, pageIndex int64, pageSize int64) (articleInfoList []*model.ArticleInfo, err error) {
	articleIdList, err := redis.LRange(model.FRIEND_ARTICLE_LIST+strconv.Itoa(userid), (pageIndex-1)*pageSize, pageIndex*pageSize-1)
	if err != nil {
		return nil, err
	}
	for _, strArticleId := range articleIdList {
		articleId, _ := strconv.Atoi(strArticleId)
		articleInfo, err := GetArticleById(articleId)
		if err != nil {
			return nil, err
		}
		articleInfoList = append(articleInfoList, articleInfo)
	}
	return
}

// GetArticleFromFollow 分页获取关注的用户的文章
func GetArticleFromFollow(userid int, pageIndex int64, pageSize int64) (articleInfoList []*model.ArticleInfo, err error) {
	articleIdList, err := redis.LRange(model.FOLLOW_ARTICLE_LIST+strconv.Itoa(userid), (pageIndex-1)*pageSize, pageIndex*pageSize-1)
	if err != nil {
		return nil, err
	}
	for _, strArticleId := range articleIdList {
		articleId, _ := strconv.Atoi(strArticleId)
		articleInfo, err := GetArticleById(articleId)
		if err != nil {
			return nil, err
		}
		articleInfoList = append(articleInfoList, articleInfo)
	}
	return
}

// GetArticleById  根据article的id获取一篇文章
func GetArticleById(articleId int) (articleInfo *model.ArticleInfo, err error) {
	articleInfo = new(model.ArticleInfo)
	if err = mysql.DB.Where("id = ?", articleId).First(articleInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	articleInfo.AuthorInfo, _ = GetAUserInfoByUserId(articleInfo.AuthorID)
	dir, _ := os.ReadDir(articleInfo.ResourceDir)
	for _, fi := range dir {
		articleInfo.ResourceUrl = append(articleInfo.ResourceUrl, articleInfo.ResourceDir+fi.Name())
	}
	articleInfo.ArticleType = GetArticleTypeById(articleInfo.ArticleTypeID).ArticleType
	return
}

// GetAllArticle  获取用户所有文章
func GetAllArticle(userid int) (articleInfoList []*model.ArticleInfo, err error) {
	rows, err := mysql.DB.Model(&model.ArticleInfo{}).Where("author_id = ?", userid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var articleInfo = &model.ArticleInfo{}
		// ScanRows 方法用于将一行记录扫描至结构体
		mysql.DB.ScanRows(rows, articleInfo)
		// 业务逻辑
		articleInfo.AuthorInfo, _ = GetAUserInfoByUserId(userid)
		dir, _ := os.ReadDir(articleInfo.ResourceDir)
		for _, fi := range dir {
			articleInfo.ResourceUrl = append(articleInfo.ResourceUrl, articleInfo.ResourceDir+fi.Name())
		}
		articleInfo.ArticleType = GetArticleTypeById(articleInfo.ArticleTypeID).ArticleType
		articleInfoList = append(articleInfoList, articleInfo)
	}
	return articleInfoList, err
}

// GetAllArticleByPage 分页获取用户所有文章
func GetAllArticleByPage(userid int, pageIndex int, pageSize int) (articleInfoList []*model.ArticleInfo, err error) {
	rows, err := mysql.DB.Model(&model.ArticleInfo{}).Where("author_id = ?", userid).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var articleInfo = &model.ArticleInfo{}
		// ScanRows 方法用于将一行记录扫描至结构体
		mysql.DB.ScanRows(rows, articleInfo)
		// 业务逻辑
		articleInfo.AuthorInfo, _ = GetAUserInfoByUserId(userid)
		dir, _ := os.ReadDir(articleInfo.ResourceDir)
		for _, fi := range dir {
			articleInfo.ResourceUrl = append(articleInfo.ResourceUrl, articleInfo.ResourceDir+fi.Name())
		}
		articleInfo.ArticleType = GetArticleTypeById(articleInfo.ArticleTypeID).ArticleType
		articleInfoList = append(articleInfoList, articleInfo)
	}
	return articleInfoList, err
}

// GiveLikeByArticleId  给文章点赞
func GiveLikeByArticleId(userid string, strArticleId string, giveLikeTime int64) (err error) {
	//文章点赞总数加1
	if err = mysql.DB.Model(&model.ArticleInfo{}).
		Where("id = ?", strArticleId).Update("give_like_count", gorm.Expr("give_like_count + ?", 1)).Error; err != nil {
		return
	}

	//返回纳秒作为zest的score ，记录点赞时间和点赞文章
	if err = redis.ZAdd(model.GIVELIKE+userid, float64(giveLikeTime), strArticleId); err != nil {
		return
	}
	//用户文章类型权重加1
	articleId, _ := strconv.Atoi(strArticleId)
	articleTypeId, err := GetArticleTypeByArticleId(articleId)
	if err != nil {
		return err
	}
	if err = WeightIncr(userid, articleTypeId); err != nil {
		return
	}
	return
}

// DelLikeByArticleId  取消文章点赞
func DelLikeByArticleId(userid string, strArticleId string) (err error) {
	if err = mysql.DB.Model(&model.ArticleInfo{}).
		Where("id = ?", strArticleId).Update("give_like_count", gorm.Expr("give_like_count - ?", 1)).Error; err != nil {
		return
	}
	if err = redis.ZRem(model.GIVELIKE+userid, strArticleId); err != nil {
		return
	}
	//用户文章类型权重减1
	articleId, _ := strconv.Atoi(strArticleId)
	articleTypeId, err := GetArticleTypeByArticleId(articleId)
	if err != nil {
		return err
	}
	if err = WeightDecr(userid, articleTypeId); err != nil {
		return
	}
	return
}

// IsGiveLikeByArticleId 判断文章是否点赞
func IsGiveLikeByArticleId(userid string, articleId string) (int, error) {
	isGiveLike, err := redis.ZRank(model.GIVELIKE+userid, articleId)
	if err != nil {
		return -1, err
	}
	return isGiveLike, nil
}

// GetGiveLikeArticleCount   获取点赞文章总数
func GetGiveLikeArticleCount(userid string) (int, error) {
	Len, err := redis.ZCard(model.GIVELIKE + userid)
	if err != nil {
		return -1, err
	}
	return int(Len), nil
}

// GetGiveLikeArticleList  获取点赞文章列表
func GetGiveLikeArticleList(userid string, pageIndex int64, pageSize int64) ([]*model.ArticleInfo, error) {
	List, err := redis.ZRevRange(model.GIVELIKE+userid, (pageIndex-1)*pageSize, pageIndex*pageSize-1)
	if err != nil {
		return nil, err
	}
	var GiveLikeArticleList []*model.ArticleInfo
	for _, strArticleId := range List {
		articleId, _ := strconv.Atoi(strArticleId)
		articleInfo, err := GetArticleById(articleId)
		if err != nil {
			return nil, err
		}
		GiveLikeArticleList = append(GiveLikeArticleList, articleInfo)
	}
	return GiveLikeArticleList, nil
}

// IncrCommentCount 评论数 + 1
func IncrCommentCount(articleId int) (err error) {
	err = mysql.DB.Model(&model.ArticleInfo{}).Where("id=?", articleId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	return
}

// DecrCommentCount 评论数 - 1
func DecrCommentCount(articleId int) (err error) {
	err = mysql.DB.Model(&model.ArticleInfo{}).Where("id=?", articleId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	return
}

// GetArticleTypeByArticleId 根据文章id获取文章类型id
func GetArticleTypeByArticleId(articleId int) (articleTypeId int, err error) {
	var articleInfo = model.ArticleInfo{}
	if err = mysql.DB.Select("article_type_id").Where("id=?", articleId).First(&articleInfo).Error; err != nil {
		return -1, err
	}
	articleTypeId = articleInfo.ArticleTypeID
	return
}
