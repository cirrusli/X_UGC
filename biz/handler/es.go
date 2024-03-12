package handler

import (
	"X_UGC/biz/dal/es"
	"X_UGC/biz/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// FindUser handles the route for searching a user by account and username
func FindUser(c *gin.Context) {
	userid := c.GetInt("userid")
	searchString := c.Query("searchString")
	pageIndex, err := strconv.Atoi(c.Query("pageIndex"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageIndex"})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageSize"})
		return
	}

	from := (pageIndex - 1) * pageSize // Calculate the starting point for the search results

	userInfoList, err := es.ES.Search(
		es.ES.Search.WithContext(context.Background()),
		es.ES.Search.WithIndex("user"),
		es.ES.Search.WithQuery(searchString),
		es.ES.Search.WithFrom(from),
		es.ES.Search.WithSize(pageSize),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = service.AddSearchRecord(userid, time.Now().UnixNano(), searchString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"userInfoList": userInfoList,
	})
}

// FindArticle handles the route for searching an article by title or content
func FindArticle(c *gin.Context) {
	// Implementation goes here
	c.JSON(http.StatusOK, gin.H{"message": "FindArticle handler"})
}

// FindArticleByType handles the route for searching an article by type, title or content
func FindArticleByType(c *gin.Context) {
	// Implementation goes here
	c.JSON(http.StatusOK, gin.H{"message": "FindArticleByType handler"})
}

// DelRecord handles the route for deleting a search record
func DelRecord(c *gin.Context) {
	// Implementation goes here
	c.JSON(http.StatusOK, gin.H{"message": "DelRecord handler"})
}

// DelAllRecords handles the route for deleting all search records
func DelAllRecords(c *gin.Context) {
	// Implementation goes here
	c.JSON(http.StatusOK, gin.H{"message": "DelAllRecords handler"})
}

// GetAllSearchRecords handles the route for getting all search records
func GetAllSearchRecords(c *gin.Context) {
	// Implementation goes here
	c.JSON(http.StatusOK, gin.H{"message": "GetAllSearchRecords handler"})
}
