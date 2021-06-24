package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/things-labs/fwu"
)

var BuildTime = "unknown"
var Version = "unknown"

func main() {
	fmt.Printf("Version: %s\r\n", Version)
	fmt.Printf("BuildTime: %s\r\n", BuildTime)

	// html
	http.HandleFunc("/", fwu.IndexHTML)
	// api
	http.HandleFunc("/api/fwu/reboot", fwu.Reboot)
	http.HandleFunc("/api/fwu/config", fwu.UploadConfigFile)
	http.HandleFunc("/api/fwu/upgrade", fwu.Upgrade)

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
