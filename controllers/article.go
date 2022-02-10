package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"itlearn/api/logger"
	"itlearn/api/logic"
	"itlearn/api/models"
	"net/http"
	"strconv"
	"time"
)

/*
当访问/add路径的时候回触发AddArticleGet方法
响应的页面是通过HTML
*/

// 获取写文章
func AddArticleGet(c *gin.Context) {
	//获取session
	isLogin := c.MustGet("is_login").(bool)
	c.HTML(http.StatusOK, "write_article.html", gin.H{"isLogin": isLogin})
}

// 提交写好的文章
func AddArticlePost(c *gin.Context) {
	//获取浏览器传输的数据，通过表单的name属性获取值
	//获取表单信息
	title := c.PostForm("title")
	tags := c.PostForm("tags")
	short := c.PostForm("short")
	content := c.PostForm("content")
	currentUser := c.MustGet("login_user").(string)

	logger.Debug("AddArticlePost", zap.String("title", title), zap.String("tags", tags))

	//实例化model，将它出入到数据库中
	art := &models.Article{
		Title:      title,
		Tags:       tags,
		Short:      short,
		Content:    content,
		Author:     currentUser,
		CreateTime: time.Now().Unix(),
	}
	_, err := models.AddArticle(art)

	//返回数据给浏览器
	response := gin.H{}

	if err == nil {
		//无误
		response = gin.H{"code": 1, "message": "ok"}
	} else {
		logger.Error("AddArticlePost failed", zap.Any("error", err))
		response = gin.H{"code": 0, "message": "error"}
	}
	c.JSON(http.StatusOK, response)
}

// 展示文章
func ShowArticleGet(c *gin.Context) {
	isLogin := c.MustGet("is_login")
	idStr := c.Param("id")

	// 查询文章
	article, err := models.QueryArticleWithId(idStr)
	if err != nil {
		logger.Error("QueryArticleWithId failed", zap.Any("error", err))
		c.String(http.StatusOK, "bad id")
		return
	}
	if article == nil {
		c.String(http.StatusOK, "bad id")
		return
	}

	// 增加文章的阅读数
	err = logic.IncArticleReadCount(idStr)
	if err != nil {
		logger.Error("ArticleReadCountIncr failed", zap.Any("error", err))
	}

	c.HTML(http.StatusOK, "show_article.html", gin.H{"isLogin": isLogin, "Title": article.Title, "Content": article.Content})
}

// UpdateArticleGet 更新文章
func UpdateArticleGet(c *gin.Context) {
	isLogin := c.MustGet("is_login")
	idStr := c.Query("id")

	//获取id所对应的文章信息
	article, err := models.QueryArticleWithId(idStr)
	if err != nil {
		logger.Error("QueryArticleWithId failed", zap.Any("error", err))
		c.String(http.StatusOK, "bad id")
		return
	}
	if article == nil {
		c.String(http.StatusOK, "bad id")
		return
	}

	c.HTML(http.StatusOK, "write_article.html", gin.H{"isLogin": isLogin, "article": article})
}

// 更新文章
func UpdateArticlePost(c *gin.Context) {
	// 获取浏览器传输的数据，通过表单的name属性获取值
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, "bad id")
	}
	title := c.PostForm("title")
	tags := c.PostForm("tags")
	short := c.PostForm("short")
	content := c.PostForm("content")

	// 实例化model，修改数据库
	art := &models.Article{
		Id:      id,
		Title:   title,
		Tags:    tags,
		Short:   short,
		Content: content,
	}

	logger.Debug("UpdateArticlePost", zap.Any("article", *art))

	_, err = models.UpdateArticle(art)

	//返回数据给浏览器
	response := gin.H{}

	if err == nil {
		//无误
		response = gin.H{"code": 1, "message": "更新成功"}
	} else {
		response = gin.H{"code": 0, "message": "更新失败"}
	}

	c.JSON(http.StatusOK, response)
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	idStr := c.Query("id")
	_, err := models.DeleteArticle(idStr)
	if err != nil {
		logger.Error("DeleteArticle failed", zap.Any("error", err))
	}
	c.Redirect(http.StatusFound, "/home")
}

// 按照阅读数排行返回前n篇文章的id和title
func ArticleTopN(c *gin.Context) {
	nStr := c.Param("n")
	n, err := strconv.ParseInt(nStr, 0, 16)
	if err != nil {
		logger.Error("ArticleTopN", zap.Any("error", err))
		n = 5
	}

	// 调用业务逻辑层 获取返回数据结果
	articleList := logic.GetArticleReadCountTopN(n)

	// 3. 返回
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": articleList,
	})
	return
}
