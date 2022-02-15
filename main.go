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
  "password": "root"
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
