package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed index.html
var indexHTML []byte

var (
	viewerCount int
	countMutex  sync.Mutex
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})
	http.HandleFunc("/events", sseHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SSE server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
