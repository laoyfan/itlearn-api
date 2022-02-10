package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"itlearn/api/logger"
	"itlearn/api/models"
	"net/http"
	"strconv"
)

// 获取首页
func HomeGet(c *gin.Context) {
	//获取session，判断用户是否登录
	isLogin := c.MustGet("is_login").(bool)
	username := c.MustGet("login_user").(string)
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	logger.Debug("HomeGet", zap.Int("page", page))

	articleList, err := models.QueryCurrUserArticleWithPage(username, page)
	if err != nil {
		logger.Error("models.QueryCurrUserArticleWithPage failed", zap.Any("error", err))
	}

	logger.Debug("models.QueryCurrUserArticleWithPage", zap.Any("articleList", articleList))

	data := models.GenHomeBlocks(articleList, isLogin)
	pageData := models.GenHomePagination(page)

	logger.Debug("models.GenHomeBlocks", zap.Any("data", data))

	c.HTML(http.StatusOK, "home.html", gin.H{"isLogin": isLogin, "data": data, "pageData": pageData})
}
