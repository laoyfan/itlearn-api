package routers

import (
	"github.com/gin-contrib/sessions"       // session包 定义了一套session操作的接口 类似于 database/sql
	"github.com/gin-contrib/sessions/redis" // session具体存储的介质
	"github.com/gin-gonic/gin"
	"html/template"
	"itlearn/api/controllers"
	"itlearn/api/logger"
	"itlearn/api/middlewares"
	"time"
)

// 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 接管Gin框架的Logger模块  和 Recovery模块
	r.Use(logger.GinLogger(logger.Logger), logger.GinRecovery(logger.Logger, true))

	// 设置时间格式
	r.SetFuncMap(template.FuncMap{
		"timeStr": func(timestamp int64) string {
			return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		},
	})

	// 配置静态文件
	r.Static("/static", "static")

	// 配置模板
	r.LoadHTMLGlob("views/*")

	// 设置session   和中间件middleware
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 登录注册 无需认证
	{
		r.GET("/register", controllers.RegisterGet)
		r.POST("/register", controllers.RegisterPost)

		r.GET("/login", controllers.LoginGet)
		r.POST("/login", controllers.LoginPost)

		// 获取阅读排行榜前几
		r.GET("/article/top/:n", controllers.ArticleTopN)
	}

	// 需要认证的一些路由
	{
		// 路由组注册中间件
		basicAuthGroup := r.Group("/", middlewares.BasicAuth())
		basicAuthGroup.GET("/home", controllers.HomeGet)
		basicAuthGroup.GET("/", controllers.IndexGet)
		basicAuthGroup.GET("/logout", controllers.LogoutHandler)

		//路由组
		article := basicAuthGroup.Group("/article")
		{
			// 写文章
			article.GET("/add", controllers.AddArticleGet)
			article.POST("/add", controllers.AddArticlePost)

			// 文章详情
			article.GET("/show/:id", controllers.ShowArticleGet)

			// 更新文章
			article.GET("/update", controllers.UpdateArticleGet)
			article.POST("/update", controllers.UpdateArticlePost)

			// 删除文章
			article.GET("/delete", controllers.DeleteArticle)

		}

		// 相册
		basicAuthGroup.GET("/album", controllers.AlbumGet)

		// 文件上传
		basicAuthGroup.POST("/upload", controllers.UploadPost)
	}

	return r
}
