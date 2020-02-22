package main

import (
	"log"
	"net/http"

	"github.com/thinkgos/anytool"
)

func main() {
	//html
	http.HandleFunc("/internal/tool", anytool.ToolHTML)
	http.HandleFunc("/internal/logs", anytool.LogsHTML)
	// api
	http.HandleFunc(anytool.URLAPIReboot, anytool.Reboot)
	http.HandleFunc(anytool.URLAPIConfig, anytool.UploadConfigFile)
	http.HandleFunc(anytool.URLAPIUpgrade, anytool.Upgrade)
	http.HandleFunc(anytool.URLAPILogs, anytool.Logs)

	if err := http.ListenAndServe(":9527", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
