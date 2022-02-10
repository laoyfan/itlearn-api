package main

import (
	"fmt"
	"itlearn/api/config"
	"itlearn/api/dao"
	"itlearn/api/logger"
	"itlearn/api/routers"
)

func main() {
	// 用conf/conf.json初始化
	//if len(os.Args) < 2 {
	//	return
	//}
	//if err := config.Init(os.Args[1]); err != nil {
	//	fmt.Printf("config.Init failed, err:%v\n", err)
	//	return
	//}

	// 调试方便，先字符串初始化
	s := `{
"server": {
  "port": 8080
},
"mysql": {
  "host": "127.0.0.1",
  "port": 3306,
  "db": "gin_blog",
  "username": "root",
  "password": "wxlzs999"
},
"redis": {
  "host": "127.0.0.1",
  "port": 6379,
  "db": 0,
  "password": ""
},
"log":{
	"level": "debug",
  "filename": "log/gin_blog.log",
  "maxsize": 500,
  "max_age": 7,
  "max_backups": 10
}
}`

	// 初始化Config的全局变量
	if err := config.InitFromStr(s); err != nil {
		fmt.Printf("config.Init failed, err:%v\n", err)
		return
	}

	// 初始化日志模块
	if err := logger.InitLogger(config.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 初始化Mysql数据库
	if err := dao.InitMySQL(config.Conf.MySQLConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}

	// 初始化redis数据库
	if err := dao.InitRedis(config.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}

	// 初始化
	logger.Logger.Info("start project...")

	r := routers.SetupRouter() // 初始化路由
	r.Run()
}

/*func main() {
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
}*/
