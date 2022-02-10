package models

import (
	"fmt"
	"go.uber.org/zap"
	"itlearn/api/logger"
	"itlearn/api/utils"
	"strconv"
	"strings"
)

type HomeBlockParam struct {
	Article *Article

	TagLinks      []*TagLink
	CreateTimeStr string
	//查看文章的地址
	Link string

	//修改文章的地址
	UpdateLink string
	DeleteLink string

	//记录是否登录
	IsLogin bool
}

type TagLink struct {
	TagName string
	TagUrl  string
}

// HomePagination 分页器
type HomePagination struct {
	HasPre   bool
	HasNext  bool
	ShowPage string
	PreLink  string
	NextLink string
}

//将tags字符串转化成首页模板所需要的数据结构
func createTagsLinks(tagStr string) []*TagLink {
	var tagLinks = make([]*TagLink, 0, strings.Count(tagStr, "&"))
	tagList := strings.Split(tagStr, "&")
	for _, tag := range tagList {
		tagLinks = append(tagLinks, &TagLink{tag, "/?tag=" + tag})
	}
	return tagLinks
}

// 生成home页面数据结构
func GenHomeBlocks(articleList []*Article, isLogin bool) (ret []*HomeBlockParam) {
	// 内存申请一次到位
	ret = make([]*HomeBlockParam, 0, len(articleList))
	for _, art := range articleList {
		// 将数据库model转换为首页模板所需要的model
		homeParam := HomeBlockParam{
			Article: art,
			IsLogin: isLogin,
		}
		homeParam.TagLinks = createTagsLinks(art.Tags)
		homeParam.CreateTimeStr = utils.SwitchTimeStampToStr(art.CreateTime)

		homeParam.Link = fmt.Sprintf("/article/show/%d", art.Id)
		homeParam.UpdateLink = fmt.Sprintf("/article/update?id=%d", art.Id)
		homeParam.DeleteLink = fmt.Sprintf("/article/delete?id=%d", art.Id)
		ret = append(ret, &homeParam) // 不再需要动态扩容
	}
	return
}

// 生成home页面分页数据结构
func GenHomePagination(page int) *HomePagination {
	pageObj := new(HomePagination)

	// 查询出总的条数
	num, _ := QueryArticleRowNum()

	// 从配置文件中读取每页显示的条数
	// 计算出总页数
	allPageNum := (num-1)/pageSize + 1

	pageObj.ShowPage = fmt.Sprintf("%d/%d", page, allPageNum)

	//当前页数小于等于1，那么上一页的按钮不能点击
	if page <= 1 {
		pageObj.HasPre = false
	} else {
		pageObj.HasPre = true
	}

	//当前页数大于等于总页数，那么下一页的按钮不能点击
	if page >= allPageNum {
		pageObj.HasNext = false
	} else {
		pageObj.HasNext = true
	}

	pageObj.PreLink = "/?page=" + strconv.Itoa(page-1)
	pageObj.NextLink = "/?page=" + strconv.Itoa(page+1)
	logger.Debug("GenHomePagination", zap.Any("pageObj", *pageObj))
	return pageObj
}
