package anytool

import (
	"log"
	"net/http"

	"github.com/thinkgos/memlog"
)

func LogsHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		html404(w, r)
		return
	}
	if err := logsTpl.Execute(w, nil); err != nil {
		log.Println("temple execute failed", err)
	}
}

type LogsInfo struct {
	List []string `json:"list"`
}

func Logs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		response(w, CodeSuccess, &LogsInfo{memlog.ReadLastMsgs()})
	} else if r.Method == http.MethodPost {
		memlog.Clear()
		response(w, CodeSuccess)
	} else {
		html404(w, r)
	}
}
