package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"itlearn/api/logger"
	"itlearn/api/models"
	"itlearn/api/utils"
	"net/http"
)

// get登陆页
func LoginGet(c *gin.Context) {
	//返回html
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "登录页"})
}

// 提交登陆
func LoginPost(c *gin.Context) {
	// 取出请求数据
	// 校验用户名密码是否正确
	// 返回响应
	username := c.PostForm("username")
	password := c.PostForm("password")
	logger.Debug("login", zap.String("username", username), zap.String("password", password))

	// 去数据库查,注意查找的时候，密码是MD5之后的密码查找
	id := models.QueryUserWithParam(username, utils.MD5(password))

	fmt.Println("id:", id)

	// 登陆成功
	if id > 0 {
		// 给响应种上Cookie
		session := sessions.Default(c)
		session.Set("login_user", username) // 在session中保存k-v,然后写入cookie
		session.Save()

		c.Redirect(http.StatusFound, "/home") // 浏览器收到这个就会跳转到我指定的页面
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "登录失败"})
	}
}

// 登出
func LogoutHandler(c *gin.Context) {
	//清除该用户登录状态的数据
	session := sessions.Default(c)
	session.Delete("login_user")
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
