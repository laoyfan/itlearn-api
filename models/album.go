package models

import "itlearn/api/dao"

type Album struct {
	Id         int
	Filepath   string
	Filename   string
	Status     int
	CreateTime int64 `db:"create_time"`
}

// 增加图片
func AddAlbum(album *Album) (int64, error) {
	return dao.ModifyDB("insert into album(filepath,filename,status,create_time)values(?,?,?,?)",
		album.Filepath, album.Filename, album.Status, album.CreateTime)
}

// 获取图片
func QueryAlbum() (dest []*Album, err error) {
	sqlStr := "select id,filepath,filename,status,create_time from album"
	err = dao.QueryRows(&dest, sqlStr)
	return
}
