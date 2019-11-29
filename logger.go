package anytool

import (
	"log"
	"net/http"

	"github.com/thinkgos/memlog"
)

// LogsHTML 日志html页面
func LogsHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		html404(w, r)
		return
	}
	if err := logsTpl.Execute(w, nil); err != nil {
		log.Println("temple execute failed", err)
	}
}

// LogsInfos 日志信息
type LogsInfos struct {
	List []string `json:"list"`
}

// Logs 日志处理handler
func Logs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		responseOK(w, &LogsInfos{memlog.ReadLastMsgs()})
	} else if r.Method == http.MethodPost {
		memlog.Clear()
		response(w, http.StatusOK)
	} else {
		html404(w, r)
	}
}
