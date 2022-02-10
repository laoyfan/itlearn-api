package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"itlearn/api/logger"
	"itlearn/api/models"
	"net/http"
	"os"
	"path"
	"time"
)

// 文件上传
func UploadPost(c *gin.Context) {
	fh, err := c.FormFile("upload")
	if err != nil {
		logger.Warn("UploadPost", zap.Any("error", err))
		c.JSON(http.StatusOK, gin.H{"msg": "无效的参数"})
		return
	}
	logger.Debug("UploadPost", zap.String("filename", fh.Filename), zap.Int64("fileSize", fh.Size))

	now := time.Now()
	fileType := "other"
	// 判断后缀为图片的文件，如果是图片我们才存入到数据库中
	fileExt := path.Ext(fh.Filename)
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".gif" || fileExt == ".jpeg" {
		fileType = "img"
	}

	// 准备好要创建的文件夹路径
	fileDir := fmt.Sprintf("static/upload/%s/%d/%d/%d", fileType, now.Year(), now.Month(), now.Day())

	// ModePerm是0777，这样拥有该文件夹路径的执行权限
	// 创建文件夹
	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "服务繁忙，稍后再试。"})
		return
	}

	// 文件路径
	timeStamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d-%s", timeStamp, fh.Filename)

	// 文件路径+文件名  拼接好
	filePathStr := path.Join(fileDir, fileName)

	// 将浏览器客户端上传的文件拷贝到本地路径的文件里面，此处也可以使用io操作
	c.SaveUploadedFile(fh, filePathStr)

	if fileType == "img" {
		album := &models.Album{Filepath: filePathStr, Filename: fileName, CreateTime: timeStamp}
		models.AddAlbum(album)
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "message": "上传成功"})
}

// 获取所有的图片
func AlbumGet(c *gin.Context) {
	isLogin := c.GetBool("is_login")
	albums, err := models.QueryAlbum()
	if err != nil {
		logger.Error("AlbumGet", zap.Any("error", err))
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "服务繁忙，请稍后再试。"})
		return
	}

	c.HTML(http.StatusOK, "album.html", gin.H{"isLogin": isLogin, "albums": albums})
}
