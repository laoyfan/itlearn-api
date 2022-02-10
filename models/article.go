package models

import (
	sql "github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"itlearn/api/dao"
	"itlearn/api/logger"
	"strings"
)

const (
	pageSize = 4
)

type Article struct {
	Id         int    `json:"id",form:"id"`
	Title      string `json:"title",form:"title"`
	Tags       string `json:"tags",form:"tags"`
	Short      string `json:"short",form:"short"`
	Content    string `json:"content",form:"content"`
	Author     string
	CreateTime int64 `db:"create_time"`
	Status     int   // Status=0为正常，1为删除，2为冻结
}

//-----------数据库操作---------------

// 增加文章
func AddArticle(article *Article) (int64, error) {
	return dao.ModifyDB("insert into article(title,tags,short,content,author,create_time,status) values(?,?,?,?,?,?,?)",
		article.Title, article.Tags, article.Short, article.Content, article.Author, article.CreateTime, article.Status)
}

// 更新文章
func UpdateArticle(article *Article) (int64, error) {
	sqlStr := "update article set title=?,tags=?,short=?,content=? where id=?"
	return dao.ModifyDB(sqlStr, article.Title, article.Tags, article.Short, article.Content, article.Id)
}

// 删除文章
func DeleteArticle(id string) (int64, error) {
	sqlStr := "delete from article where id=?"
	return dao.ModifyDB(sqlStr, id)
}

// 查询所有文章

/**
分页查询数据库
limit分页查询语句，
    语法：limit m,n

    m代表从多少位开始获取，与id值无关
    n代表获取多少条数据

	总共有10条数据，每页显示4条。  --> 总共需要(10-1)/4+1 页。
	问第2页数据是哪些？           --> 5,6,7,8  (2-1)*4,4

*/
// 查询数据库文章
func QueryAllArticle() ([]*Article, error) {
	sqlStr := "select id,title,tags,short,content,author,create_time from article"
	var articleList []*Article
	err := dao.QueryRows(&articleList, sqlStr)
	if err != nil {
		return nil, err
	}
	return articleList, nil
}

// 根据Page查询文章
func QueryCurrUserArticleWithPage(username string, pageNum int) (articleList []*Article, err error) {
	sqlStr := "select id,title,tags,short,content,author,create_time from article where author=? limit ?,?"

	articleList, err = queryArticleWithCon(pageNum, sqlStr, username)
	if err != nil {
		logger.Debug("queryArticleWithCon, ", zap.Any("error", err))
		return nil, err
	}

	logger.Debug("QueryCurrUserArticleWithPage,", zap.Any("articleList", articleList))
	return articleList, nil
}

// 根据Id查询文章
func QueryArticleWithId(id string) (article *Article, err error) {
	article = new(Article)
	sqlStr := "select id,title,tags,short,content,author,create_time from article where id=?"
	err = dao.QueryRowDB(article, sqlStr, id)
	return
}

// 根据查询条件查询指定页数有的文章
func queryArticleWithCon(pageNum int, sqlStr string, args ...interface{}) (articleList []*Article, err error) {
	pageNum--
	args = append(args, pageNum*pageSize, pageSize)
	logger.Debug("queryArticleWithCon", zap.Any("pageNum", pageNum), zap.Any("args", args))
	err = dao.QueryRows(&articleList, sqlStr, args...)
	logger.Debug("dao.QueryRows result", zap.Any("articleList", articleList))
	return
}

// 查询文章的总条数
func QueryArticleRowNum() (num int, err error) {
	err = dao.QueryRowDB(&num, "select count(id) from article")
	return
}

// 根据id查文章 按顺序
func QueryArticlesByIds(ids []int64, idStrs []string) ([]*Article, error) {
	// 让MySQL排序
	query, args, err := sql.In("select id, title from article where id in (?) order by FIND_IN_SET(id, ?)", ids, strings.Join(idStrs, ","))
	if err != nil {
		logger.Error("QueryArticlesByIds", zap.Any("error", err))
		return nil, err
	}

	var dest []*Article
	err = dao.QueryRows(&dest, query, args...)
	return dest, err
}
