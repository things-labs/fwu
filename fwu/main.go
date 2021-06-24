package main

import (
	"log"
	"net/http"

	"github.com/things-labs/fwu"
)

func main() {
	//html
	http.HandleFunc("/", fwu.ToolHTML)
	// api
	http.HandleFunc(fwu.URLAPIReboot, fwu.Reboot)
	http.HandleFunc(fwu.URLAPIConfig, fwu.UploadConfigFile)
	http.HandleFunc(fwu.URLAPIUpgrade, fwu.Upgrade)

	if err := http.ListenAndServe(":9527", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
