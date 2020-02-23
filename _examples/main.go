package main

import (
	"log"
	"net/http"

	"github.com/thinkgos/anytool"
)

func main() {
	//html
	http.HandleFunc("/tools", anytool.ToolHTML)
	// api
	http.HandleFunc(anytool.URLAPIReboot, anytool.Reboot)
	http.HandleFunc(anytool.URLAPIConfig, anytool.UploadConfigFile)
	http.HandleFunc(anytool.URLAPIUpgrade, anytool.Upgrade)

	if err := http.ListenAndServe(":9527", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
