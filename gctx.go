package fwu

import (
	"net/http"
	"strconv"
	"time"

	"github.com/things-go/render"
)

const CustomCode = 999999
const CustomStatus = 499

type Error struct {
	Code    int
	Message string
	Detail  string
}

func (e *Error) Error() string {
	if e.Detail == "" {
		return "code: " + strconv.Itoa(e.Code) + ", message: " + e.Message
	}
	return "code: " + strconv.Itoa(e.Code) + ", message: " + e.Message + ", detail: " + e.Detail
}

func New(code int, args ...string) error {
	err := &Error{Code: code}
	if len(args) >= 1 {
		err.Message = args[0]
	}
	if len(args) >= 2 {
		err.Detail = args[1]
	}
	return err
}

func NewCustomError(args ...string) error {
	return New(CustomCode, args...)
}

func Parse(err error) *Error {
	e, ok := err.(*Error)
	if !ok {
		return &Error{
			-1,
			"",
			err.Error(),
		}
	}
	if e.Code == 0 && e.Message == "" {
		e.Code = -1
		e.Detail = err.Error()
	}
	return e
}

// ResponseBody 回复基本格式
type ResponseBody struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Ts     time.Time   `json:"ts"`
	Data   interface{} `json:"data"`
}

func Response(w http.ResponseWriter, obj interface{}) {
	mp := render.H{
		"code": 0,
		"msg":  "ok",
		"ts":   time.Now().UnixNano(),
	}

	if obj != nil {
		mp["data"] = obj
	} else {
		mp["data"] = render.H{}
	}
	render.JSON(w, http.StatusOK, mp)
}

func ResponseOK(w http.ResponseWriter) {
	Response(w, nil)
}

func AbortErrWithStatus(w http.ResponseWriter, status int, err error, data ...interface{}) {
	e := Parse(err)
	mp := map[string]interface{}{
		"code": e.Code,
		"ts":   time.Now().UnixNano(),
	}
	if e.Message != "" {
		mp["msg"] = e.Message
	}
	if e.Detail != "" {
		mp["detail"] = e.Detail
	}
	if len(data) > 0 && data[0] != nil {
		mp["data"] = data[0]
	} else {
		mp["data"] = render.H{}
	}
	if e.Code == -1 {
		status = http.StatusInternalServerError
	}
	render.JSON(w, status, mp)
}

func AbortError(w http.ResponseWriter, err error) {
	AbortErrWithStatus(w, CustomStatus, err)
}

func AbortErrBadRequest(w http.ResponseWriter, err error) {
	AbortErrWithStatus(w, http.StatusBadRequest, err)
}
