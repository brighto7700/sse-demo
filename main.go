package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// viewerCount tracks connected clients.
// countMutex protects it from concurrent writes.
var (
	viewerCount int
	countMutex  sync.Mutex
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// SSE requires these three headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Flusher lets us push each event immediately instead of buffering.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported by this server", http.StatusInternalServerError)
		return
	}

	// Increment on connect, decrement on disconnect.
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
			// Client disconnected — the defer above handles cleanup.
			return
		case t := <-ticker.C:
			countMutex.Lock()
			current := viewerCount
			countMutex.Unlock()

			payload := fmt.Sprintf(`{"time": "%s", "viewers": %d}`, t.Format(time.Kitchen), current)
			// SSE format: "data: " prefix, double newline to signal end of event.
			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/events", sseHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SSE server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
