package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"itlearn/api/logger"
	"itlearn/api/models"
	"itlearn/api/utils"
	"net/http"
	"time"
)

// 获取注册
func RegisterGet(c *gin.Context) {
	// 返回html
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "注册页"})
}

// 注册提交
func RegisterPost(c *gin.Context) {
	// 取出请求的数据
	// 判断注册是否重复  --> 拿着用户名去数据库查一下有没有
	// 写入数据库
	// 获取表单信息
	username := c.PostForm("username")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	logger.Debug(fmt.Sprintf("%s %s %s", username, password, repassword))

	// 注册之前先判断该用户名是否已经被注册，如果已经注册，返回错误
	id := models.QueryUserWithUsername(username)
	fmt.Println("id:", id)
	if id > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "用户名已经存在"})
		return
	}

	// 注册用户名和密码
	// 存储的密码是md5后的数据，那么在登录的验证的时候，也是需要将用户的密码md5之后和数据库里面的密码进行判断
	password = utils.MD5(password)
	logger.Debug(fmt.Sprintf("password after md5:%s", password))

	user := models.User{
		Username:   username,
		Password:   password,
		Status:     0,
		CreateTime: time.Now().Unix(),
	}
	_, err := models.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "注册失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "注册成功"})
	}
}
