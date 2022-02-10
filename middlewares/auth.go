package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 最基础的认证校验 只要cookie中带了login_user标识就认为是登录用户
func BasicAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		// c代表了请求相关的所有内容，获取当前请求对应的session数据
		session := sessions.Default(c)

		loginUser := session.Get("login_user")
		// 请求对应的session中找不到我想要的数据，说明不是登录的用户
		if loginUser == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort() // 终止当前请求的处理函数调用链
			return    // 终止当前处理函数
		}

		// 根据loginUser 去数据库里用户对象取出来  gob是go语言里面二进制的数据格式
		// 如果是一个登录的用户，我就在c上设置两个自定义的键值对！！！
		c.Set("is_login", true)
		c.Set("login_user", loginUser)
		c.Next()
	}
}
