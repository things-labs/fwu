package main

import (
	"log"
	"net/http"

	_ "github.com/thinkgos/anytool"
	"github.com/thinkgos/memlog"
)

func main() {
	for i := 0; i < 100; i++ {
		memlog.Debug("hello world")
	}

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
