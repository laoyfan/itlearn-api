package logic

import (
	"fmt"
	"go.uber.org/zap"
	"itlearn/api/dao"
	"itlearn/api/logger"
	"itlearn/api/models"
	"strconv"
	"time"
)

// 点击文章，阅读数加 1
// 每次请求`/article/show/:id`URL的时候 执行redis命令 zincrby code_language 1 golang

// 给指定文章的阅读数+1
func IncArticleReadCount(articleId string) error {
	// zincrby code_language 1 golang
	todayStr := time.Now().Format("20060102")
	key := fmt.Sprintf(dao.KeyArticleCount, todayStr)

	return dao.Client.ZIncrBy(key, 1, articleId).Err()
}

// 获取阅读排行榜排名前N的文章
func GetArticleReadCountTopN(n int64) []*models.Article {
	// 1. zrevrange Key 0 n-1 从redis取出前n位的文章id
	todayStr := time.Now().Format("20060102")
	key := fmt.Sprintf(dao.KeyArticleCount, todayStr)
	idStrs, err := dao.Client.ZRevRange(key, 0, n-1).Result()
	if err != nil {
		logger.Error("ZRevRange", zap.Any("error", err))
	}

	// 2. 根据上一步获取的文章id查询数据库取文章标题  ["3" "1" "5"]
	// select id, title from article where id in (3, 1, 5);  // 文章的顺序对吗？ 不对
	// 		1. 让MySQL排序
	// select id, title from article where id in (3, 1, 5) order by FIND_IN_SET(id, (3, 1, 5));
	// 		2. 查询出来自己排序
	// 先准备好要查询的ID Slice
	var ids = make([]int64, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 0, 16)
		if err != nil {
			logger.Warn("ArticleTopN:strconv.ParseInt failed", zap.Any("error", err))
			continue
		}

		ids = append(ids, id)
	}

	articleList, err := models.QueryArticlesByIds(ids, idStrs)
	if err != nil {
		logger.Error("queryArticlesByIds", zap.Any("error", err))
	}
	return articleList
}
