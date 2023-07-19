package handler

import (
	"X_UGC/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// SearchUser 根据账号和用户名搜索用户
func SearchUser(c *gin.Context) {
	userid := c.GetInt("userid")
	searchString := c.Query("searchString")
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	userInfoList, err := service.SearchUser(searchString, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = service.AddSearchRecord(userid, time.Now().UnixNano(), searchString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"userInfoList": userInfoList,
	})
}

// SearchArticle  根据标题或内容搜索文章
func SearchArticle(c *gin.Context) {
	userid := c.GetInt("userid")
	searchString := c.Query("searchString")
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	articleInfoList, err := service.SearchArticle(searchString, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = service.AddSearchRecord(userid, time.Now().UnixNano(), searchString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"articleInfoList": articleInfoList,
	})
}

// SearchArticleByType 根据文章类型以及标题或内容搜索文章
func SearchArticleByType(c *gin.Context) {
	userid := c.GetInt("userid")
	searchString := c.Query("searchString")
	articleTypeId, _ := strconv.Atoi(c.Query("article_type_id"))
	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	articleInfoList, err := service.SearchArticleByType(searchString, articleTypeId, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = service.AddSearchRecord(userid, time.Now().UnixNano(), searchString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"articleInfoList": articleInfoList,
	})
}

// DelSearchRecord 删除搜索记录
func DelSearchRecord(c *gin.Context) {
	userid := c.GetInt("userid")
	searchString := c.Query("searchString")
	err := service.DelSearchRecord(userid, searchString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DelAllSearchRecord 删除所有搜索记录
func DelAllSearchRecord(c *gin.Context) {
	userid := c.GetInt("userid")
	err := service.DelAllSearchRecord(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// GetAllSearchRecord   获取所有搜索记录
func GetAllSearchRecord(c *gin.Context) {
	userid := c.GetInt("userid")
	pageIndex, _ := strconv.ParseInt(c.Query("pageIndex"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("pageSize"), 10, 64)
	searchRecordList, err := service.GetAllSearchRecord(userid, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"searchRecordList": searchRecordList,
	})
}
