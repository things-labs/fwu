//go:generate stringer -type=Code -linecomment
package anytool

import (
	"encoding/json"
	"net/http"
)

type Code int

// code 错误码
const (
	CodeFeatureNotSupport Code = 1000 + iota // feature Not support
	CodeOperationFailed                      // operate failed
	CodeInvalidArguments                     // Invalid arguments
	CodeExistYet                             // exist yet
)

// Response 回复基本格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func response(w http.ResponseWriter, code Code, data ...interface{}) {
	var value interface{}

	if len(data) > 0 {
		value = data[0]
	} else {
		value = "{}"
	}

	if code < 1000 {
		JSON(w, int(code), Response{
			Code:    int(code),
			Message: http.StatusText(int(code)),
			Data:    value,
		})
	} else {
		JSON(w, int(code), Response{
			Code:    int(code),
			Message: code.String(),
			Data:    value,
		})
	}
}

// JSON json传输
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	content, err := json.Marshal(payload)
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
