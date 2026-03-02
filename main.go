package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	viewerCount int
	countMutex  sync.Mutex
)

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go SSE Demo</title>
    <style>
        body { font-family: system-ui, sans-serif; padding: 2rem; background: #0f172a; color: white; }
        .box { padding: 1rem; background: #1e293b; border-radius: 8px; border: 1px solid #334155; margin-bottom: 1rem; }
        .highlight { font-family: monospace; color: #10b981; font-size: 1.5rem; }
        .badge { background: #ef4444; color: white; padding: 0.2rem 0.8rem; border-radius: 9999px; font-weight: bold; }
    </style>
</head>
<body>
    <h1>🚀 Go SSE — Live Viewer Counter</h1>
    <div class="box">
        <p>Live Viewers: <span id="viewers-count" class="badge">0</span> 👀</p>
    </div>
    <div class="box">
        <p>Server Time:</p>
        <div id="live-time" class="highlight">Connecting...</div>
    </div>
    <script>
        const evtSource = new EventSource('/events');
        evtSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            document.getElementById('live-time').innerText = data.time;
            document.getElementById('viewers-count').innerText = data.viewers;
        };
        evtSource.onerror = function() {
            document.getElementById('live-time').innerText = "Reconnecting...";
        };
    </script>
</body>
</html>`

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported by this server", http.StatusInternalServerError)
		return
	}

	countMutex.Lock()
	viewerCount++
	countMutex.Unlock()

	defer func() {
		countMutex.Lock()
		viewerCount--
		countMutex.Unlock()
		log.Printf("Client disconnected. Live viewers: %d", viewerCount)
	}()

	log.Printf("New client connected. Live viewers: %d", viewerCount)

	ctx := r.Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			countMutex.Lock()
			current := viewerCount
			countMutex.Unlock()

			payload := fmt.Sprintf(`{"time": "%s", "viewers": %d}`, t.Format(time.Kitchen), current)
			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, indexHTML)
	})
	http.HandleFunc("/events", sseHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SSE server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
