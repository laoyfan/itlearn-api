package controllers

const (
	CodeSuccess      = 2000
	CodeBadRequest   = 2001
	CodeInvalidParam = 2002
	CodeFailed       = 5000
	CodeError        = 5001
)

var MsgMap = map[int]string{
	CodeSuccess:      "success",
	CodeBadRequest:   "bad request",
	CodeInvalidParam: "无效的参数",
	CodeFailed:       "请求失败",
	CodeError:        "啊哦，服务器走丢了",
}

func ShowMsg(code int) string {
	v, ok := MsgMap[code]
	if !ok {
		return MsgMap[CodeError]
	}
	return v
}
