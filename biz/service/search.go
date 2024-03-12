package service

import (
	"X_UGC/biz/dal/mysql"
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"os"
	"strconv"
)

// AddSearchRecord  添加搜索记录
func AddSearchRecord(userid int, searchTime int64, searchString string) (err error) {
	err = redis.ZAdd(model.SEARCH_RECORD+strconv.Itoa(userid), float64(searchTime), searchString)
	return
}

// DelSearchRecord 删除搜索记录
func DelSearchRecord(userid int, searchString string) (err error) {
	err = redis.ZRem(model.SEARCH_RECORD+strconv.Itoa(userid), searchString)
	return
}

// DelAllSearchRecord 删除搜索记录
func DelAllSearchRecord(userid int) (err error) {
	err = redis.Del(model.SEARCH_RECORD + strconv.Itoa(userid))
	return
}

// GetAllSearchRecord  查找搜索记录
func GetAllSearchRecord(userid int, pageIndex int64, pageSize int64) (searchRecordList []string, err error) {
	searchRecordList, err = redis.ZRevRange(model.SEARCH_RECORD+strconv.Itoa(userid), (pageIndex-1)*pageSize, pageIndex*pageSize-1)
	return
}

// SearchUserByUserName 根据用户名搜索用户
func SearchUserByUserName(searchString string, pageIndex int, pageSize int) (userInfoList []*model.UserInfo, err error) {
	err = mysql.DB.Select("*").Where("instr(user_name,?)", searchString).
		Joins("inner join follow_fans_counts on user_infos.user_id=follow_fans_counts.user_id").
		Order("fans_count desc").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&userInfoList).Error
	return
}

// SearchUser 根据账号和用户名搜索用户
func SearchUser(searchString string, pageIndex int, pageSize int) (userInfoList []*model.UserInfo, err error) {
	err = mysql.DB.Select("*").Where("instr(number,?) or instr(email,?) or instr(user_name,?)", searchString, searchString, searchString).
		Joins("inner join users on user_infos.user_id=users.id").
		Joins("inner join follow_fans_counts on user_infos.user_id=follow_fans_counts.user_id").
		Order("fans_count desc").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).
		Find(&userInfoList).Error
	return
}

// SearchArticle  根据标题或内容搜索文章
func SearchArticle(searchString string, pageIndex int, pageSize int) (articleInfoList []*model.ArticleInfo, err error) {
	rows, err := mysql.DB.Model(&model.ArticleInfo{}).
		Where("instr(title,?) or instr(content,?)", searchString, searchString).
		Order("give_like_count desc,comment_count desc,id desc").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var articleInfo = &model.ArticleInfo{}
		// ScanRows 方法用于将一行记录扫描至结构体
		mysql.DB.ScanRows(rows, articleInfo)
		// 业务逻辑
		articleInfo.AuthorInfo, _ = GetAUserInfoByUserId(articleInfo.AuthorID)
		dir, _ := os.ReadDir(articleInfo.ResourceDir)
		for _, fi := range dir {
			articleInfo.ResourceUrl = append(articleInfo.ResourceUrl, articleInfo.ResourceDir+fi.Name())
		}
		articleInfo.ArticleType = GetArticleTypeById(articleInfo.ArticleTypeID).ArticleType
		articleInfoList = append(articleInfoList, articleInfo)
	}
	return
}

// SearchArticleByType 根据文章类型以及标题或内容搜索文章
func SearchArticleByType(searchString string, articleTypeId int, pageIndex int, pageSize int) (articleInfoList []*model.ArticleInfo, err error) {
	rows, err := mysql.DB.Model(&model.ArticleInfo{}).
		Where("(instr(title,?) or instr(content,?)) and article_type_id=?", searchString, searchString, articleTypeId).
		Order("give_like_count desc,comment_count desc,id desc").
		Offset((pageIndex - 1) * pageSize).Limit(pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var articleInfo = &model.ArticleInfo{}
		// ScanRows 方法用于将一行记录扫描至结构体
		mysql.DB.ScanRows(rows, articleInfo)
		// 业务逻辑
		articleInfo.AuthorInfo, _ = GetAUserInfoByUserId(articleInfo.AuthorID)
		dir, _ := os.ReadDir(articleInfo.ResourceDir)
		for _, fi := range dir {
			articleInfo.ResourceUrl = append(articleInfo.ResourceUrl, articleInfo.ResourceDir+fi.Name())
		}
		articleInfo.ArticleType = GetArticleTypeById(articleInfo.ArticleTypeID).ArticleType
		articleInfoList = append(articleInfoList, articleInfo)
	}
	return
}
