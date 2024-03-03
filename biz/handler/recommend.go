package handler

import (
	"X_UGC/biz/model"
	"X_UGC/biz/service"
	"X_UGC/pkg/common/bloomfilter"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// PushArticleFeed feed推荐获取article
func PushArticleFeed(c *gin.Context) {
	userid := c.GetInt("userid")
	articleIDs, err := service.PushArticleFeed(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	var articleInfoList []*model.ArticleInfo
	for _, articleId := range articleIDs {
		articleInfo, err := service.GetArticleById(articleId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		articleInfoList = append(articleInfoList, articleInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"articleInfoList": articleInfoList,
	})
}

// PushArticleByType 具体类型文章的推送
func PushArticleByType(c *gin.Context) {
	userid := c.GetInt("userid")
	strArticleTypeId := c.Query("article_type_id")
	articleTypeId, _ := strconv.Atoi(strArticleTypeId)
	articleIDs, err := service.PushArticleByTypeID(userid, articleTypeId, &bloomfilter.BloomFilter{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	var articleInfoList []*model.ArticleInfo
	for _, articleId := range articleIDs {
		articleInfo, err := service.GetArticleById(articleId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		articleInfoList = append(articleInfoList, articleInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"articleInfoList": articleInfoList,
	})
}
