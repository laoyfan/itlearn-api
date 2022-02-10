package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

const (
	secret = "你猜不到的东西"
)

//传入的数据不一样，那么MD5后的32位长度的数据肯定会不一样
func MD5(str string) string {
	md5str := fmt.Sprintf("%x", md5.Sum(append([]byte(str), []byte(secret)...)))
	return md5str
}

//将传入的时间戳转为时间
func SwitchTimeStampToStr(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
