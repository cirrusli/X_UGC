package handler

import (
	"X_UGC/biz/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ArticleExist 检查视频文章是否存在
func ArticleExist(c *gin.Context) {
	userid := c.GetInt("userid")
	hashPath := c.Query("hashPath")
	chunkIndexList, FilePath, err := service.ArticleExist(userid, hashPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":         "success",
			"articleFileDir": FilePath,
			"chunkIndexList": chunkIndexList,
		})
	}
}

// ChunkUpload 分块上传文件
func ChunkUpload(c *gin.Context) {
	userid := c.GetInt("userid")
	hashPath := c.PostForm("hashPath")
	chunkIndex := c.PostForm("chunkIndex")
	chunkFile, err := c.FormFile("chunkFile")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error:": err.Error()})
		return
	}
	err = service.ChunkUpload(chunkFile, hashPath, chunkIndex, userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

// MergeChunk 合并切片
func MergeChunk(c *gin.Context) {
	userid := c.GetInt("userid")
	strChunkTotal := c.PostForm("chunkTotal")
	chunkTotal, _ := strconv.Atoi(strChunkTotal)
	strFileSize := c.PostForm("fileSize")
	fileSize, _ := strconv.ParseInt(strFileSize, 10, 64)
	hashPath := c.PostForm("hashPath")
	fileExt := c.PostForm("fileExt")
	FilePath, CoverName, err := service.MergeChunk(chunkTotal, fileSize, hashPath, fileExt, userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":         "success",
			"articleFileDir": FilePath,
			"coverFilePath":  CoverName,
		})
	}
}

// CancelChunk 取消上传
func CancelChunk(c *gin.Context) {
	userid := c.GetInt("userid")
	hashPath := c.PostForm("hashPath")
	fileExt := c.PostForm("fileExt")
	err := service.CancelChunk(userid, hashPath, fileExt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error:": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
