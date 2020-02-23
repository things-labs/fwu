package anytool

import (
	"encoding/json"
	"net/http"
)

// code 错误码
const (
	CodeSysFeatureNotSupport = 1000 + iota // 不支持
	CodeSysOperationFailed                 // 操作失败
	CodeSysInvalidArguments                // 无效参数
	CodeSysExistYet                        // 已存在
)

var errMessage = map[int]string{
	CodeSysFeatureNotSupport: "feature Not support", // 不支持
	CodeSysOperationFailed:   "Operate failed",
	CodeSysInvalidArguments:  "Invalid arguments",
	CodeSysExistYet:          "exist yet",
}

// codeMessage 根据code返回错误码信息
func codeMessage(code int) string {
	var msg string

	if code < 1000 {
		msg = http.StatusText(code)
	} else {
		msg = errMessage[code]
	}

	if msg == "" {
		msg = errMessage[CodeSysFeatureNotSupport]
	}

	return msg
}

// Response 回复基本格式
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func response(w http.ResponseWriter, code int) {
	rsp := Response{Code: code}
	if code != http.StatusOK {
		rsp.Message = codeMessage(code)
	}

	if code < 1000 {
		JSON(w, code, rsp)
	} else {
		JSON(w, http.StatusBadRequest, rsp)
	}
}

func responseOK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// JSON json传输
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(content)
	if err != nil {
		panic(err)
	}
}
