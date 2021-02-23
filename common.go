package anytool

import (
	"net/http"

	"github.com/thinkgos/render"
)

// Response 回复基本格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func response(w http.ResponseWriter, code int, data ...interface{}) {
	var value interface{}

	if len(data) > 0 {
		value = data[0]
	} else {
		value = "{}"
	}

	render.JSON(w, code, Response{
		Code:    code,
		Message: http.StatusText(code),
		Data:    value,
	})
}
