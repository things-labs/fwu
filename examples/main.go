package main

import (
	"log"
	"net/http"

	_ "github.com/thinkgos/anytool"
)

func main() {
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
