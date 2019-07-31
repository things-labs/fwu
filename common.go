package anytool

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var Tpl404 = template.Must(template.New("logs").Parse(`<html><body>404 page not found</body></html>`))

func html404(w http.ResponseWriter, r *http.Request) {
	if err := Tpl404.Execute(w, nil); err != nil {
		log.Printf("Tpl404 template execute failed", err)
	}
}

// code 错误码
const (
	CodeSuccess                 = iota // 成功
	CodeSysException                   // 内部异常
	CodeSysFeatureNotSupport           // 不支持
	CodeSysOperationFailed             // 操作失败
	CodeSysInvalidArguments            // 无效参数
	CodeSysAuthorizationFailure        // 授权失败
	CodeSysExistYet                    // 已存在
	CodeSysInProcess                   // 正在处理中
	CodeSysFileCorrupted               // 文件已损坏
	CodeSysInvalidLicence              // 无效Licence
)

var errMessage = map[int]string{
	CodeSuccess:                 "Success",
	CodeSysException:            "Internal system exception",
	CodeSysFeatureNotSupport:    "feature Not support", // 不支持
	CodeSysOperationFailed:      "Operate failed",
	CodeSysInvalidArguments:     "Invalid arguments",
	CodeSysAuthorizationFailure: "Authorization failure",
	CodeSysExistYet:             "exist yet",
	CodeSysInProcess:            "system in process",
}

// codeErrMessage 根据code返回错误码信息
func codeErrMessage(code int) string {
	errMsg, ok := errMessage[code]
	if !ok {
		errMsg = errMessage[CodeSysFeatureNotSupport]
	}

	return errMsg
}

// Response 回复基本格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// response 信息回复
func response(w http.ResponseWriter, code int, py ...interface{}) {
	rsp := Response{Code: code}
	if len(py) > 0 {
		rsp.Payload = py[0]
	}
	if code != CodeSuccess {
		rsp.Message = codeErrMessage(code)
	}
	b, err := json.Marshal(rsp)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewReader(b))
}
