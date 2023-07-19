package service

import (
	"X_UGC/pkg/common/ffmpeg"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var chunkSize int64 = 1024 * 1024 * 5

// FileExist 判断文件或文件夹是否存在
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

// CreateFile 创建文件
func CreateFile(filePath string) (bool, error) {
	fileBool, err := FileExist(filePath)
	if fileBool && err == nil {
		return true, errors.New("文件以存在")
	} else {
		newFile, err := os.Create(filePath)
		defer newFile.Close()
		if err != nil {
			return false, errors.New("创建文件失败")
		}
	}
	return true, nil
}

// ArticleExist 检查文章是否存在       1、返回已存在的且满足条件的小分块的下标 2、若存在，返回完整路径
func ArticleExist(userid int, hashPath string) ([]string, string, error) {
	DirPath := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/video_" + hashPath
	FilePath := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/video_" + hashPath + "/article/"
	var chunkIndexList []string
	exist, _ := FileExist(DirPath)
	if !exist {
		return nil, "", nil
	} else {
		fileExist, _ := FileExist(FilePath)
		if fileExist {
			return nil, FilePath, nil
		} else {
			dir, _ := ioutil.ReadDir(DirPath)
			for _, fi := range dir {
				if fi.Size() != chunkSize {
					//首先文件名去除后缀，然后文件去除hash名和‘_’，得到不满足条件的chunk下标
					chunkIndex := strings.TrimPrefix(strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name())), hashPath+"_")
					chunkIndexList = append(chunkIndexList, chunkIndex)
				}
			}
			return chunkIndexList, "", nil
		}
	}
}

// ChunkUpload 分块上传文件
func ChunkUpload(chunkFile *multipart.FileHeader, hashPath string, chunkIndex string, userid int) error {
	DirPath := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/video_" + hashPath
	os.MkdirAll(DirPath, os.ModePerm)
	FilePath := DirPath + "/" + hashPath + "_" + chunkIndex
	fileBool, err := CreateFile(FilePath)
	if !fileBool {
		return err
	}
	// 获取现在文件大小
	fi, _ := os.Stat(FilePath)
	// 判断文件是否传输完成
	if fi.Size() == chunkFile.Size {
		return nil
	}
	start := fi.Size()
	// 进行断点上传
	// 打开之前上传文件
	newFile, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer newFile.Close()
	if err != nil {
		return errors.New("打开之前上传文件不存在")
	}
	upFile, _ := chunkFile.Open()
	err = UploadFile(upFile, start, newFile, start)
	return err
}

// UploadFile 上传文件
func UploadFile(upFile multipart.File, upSeek int64, newFile *os.File, fSeek int64) error {
	// 设置上传偏移量
	_, _ = upFile.Seek(upSeek, 0)
	// 设置文件偏移量
	_, _ = newFile.Seek(fSeek, 0)
	data := make([]byte, 1024, 1024)
	upFileReader := bufio.NewReader(upFile)
	newFileWriter := bufio.NewWriter(newFile)
	for {
		total, err := upFileReader.Read(data)
		if err == io.EOF {
			break
		}
		_, err = newFileWriter.Write(data[:total])
		if err != nil {
			return errors.New("文件上传失败")
		}
		if err = newFileWriter.Flush(); err != nil {
			return errors.New("文件上传失败")
		}
	}
	return nil
}

// MergeChunk 合并切片
func MergeChunk(chunkTotal int, fileSize int64, hashPath string, fileExt string, userid int) (string, string, error) {
	var lock sync.WaitGroup
	DirPath := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/video_" + hashPath
	dir, _ := ioutil.ReadDir(DirPath)
	var totalSize int64 = 0
	for _, fi := range dir {
		totalSize += fi.Size()
	}
	if len(dir) == chunkTotal && totalSize == fileSize {
		// 新文件创建
		FilePath := DirPath + "/article"
		FileName := FilePath + "/article" + fileExt
		_ = os.MkdirAll(FilePath, os.ModePerm)
		fileBool, err := CreateFile(FileName)
		if !fileBool {
			return "", "", err
		}
		// 读取文件片段 进行合并
		for i := 0; i < chunkTotal; i++ {
			lock.Add(1)
			go MergeFile(i, DirPath, hashPath, FileName, &lock)
		}
		lock.Wait()
		// 最终文件大小获取
		file, err := os.Open(FileName)
		if err != nil {
			return "", "", err
		}
		defer file.Close()
		md5Hash := md5.New()
		_, _ = io.Copy(md5Hash, file)
		fi, _ := file.Stat()
		//对比文件大小以及MD5Hash值是否一致
		if finalSize := fi.Size(); finalSize == finalSize && hashPath == hex.EncodeToString(md5Hash.Sum(nil)) {
			_ = os.MkdirAll(DirPath+"/cover", os.ModePerm)
			//生成封面
			CoverName := DirPath + "/cover/cover.jpg"
			ffmpeg.GetCover(FileName, CoverName)
			return FilePath + "/", CoverName, nil
		} else {
			return "", "", nil
		}
	} else {
		return "", "", nil
	}
}

// MergeFile 合并切片文件
func MergeFile(i int, DirPath, hashPath, FileName string, lock *sync.WaitGroup) {
	// 打开新文件文件
	file, err := os.OpenFile(FileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
	// 设置文件写入偏移量
	_, _ = file.Seek(chunkSize*int64(i), 0)
	iSize := strconv.Itoa(i)
	chunkFilePath := DirPath + "/" + hashPath + "_" + iSize
	chunkFile, err := os.Open(chunkFilePath)
	defer chunkFile.Close()
	if err != nil {
		log.Println(errors.New("打开分片文件失败"))
		return
	}
	// 写入数据
	data := make([]byte, 1024, 1024)
	chunkFileReader := bufio.NewReader(chunkFile)
	fileWriter := bufio.NewWriter(file)
	for {
		tal, err := chunkFileReader.Read(data)
		if err == io.EOF {
			// 删除文件 需要先关闭改文件
			chunkFile.Close()
			err := os.Remove(chunkFilePath)
			if err != nil {
				log.Println(errors.New("临时记录文件删除失败"))
				return
			}
			break
		}
		_, err = fileWriter.Write(data[:tal])
		if err != nil {
			log.Println(errors.New("文件合并失败"))
			return
		}
		if err = fileWriter.Flush(); err != nil {
			log.Println(errors.New("文件合并失败"))
			return
		}
	}
	lock.Done()
}

// CancelChunk 取消上传
func CancelChunk(userid int, hashPath string, fileExt string) error {
	DirPath := "./upload/article_resource/user_" + strconv.Itoa(userid) + "/video_" + hashPath
	FileName := DirPath + "/article/article" + fileExt
	exist, _ := FileExist(FileName)
	if exist {
		return nil
	}
	err := os.RemoveAll(DirPath)
	if err != nil {
		return err
	}
	return nil
}
