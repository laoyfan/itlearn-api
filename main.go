package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/api/auth/register", func(c *gin.Context) {
		//获取参数
		name := c.PostForm("name")
		phone := c.PostForm("phone")
		password := c.PostForm("password")

		if len(phone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}

		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}

		if len(name) == 0 {
			name = RandomString(10)
		}

		log.Println(name, phone, password)

		c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	resule := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range resule {
		resule[i] = letters[rand.Intn(len(letters))]
	}
	return string(resule)
}
