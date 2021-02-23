package anytool

import (
	"encoding/json"
	"net/http"
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

	JSON(w, code, &Response{
		Code:    code,
		Message: http.StatusText(code),
		Data:    value,
	})
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
