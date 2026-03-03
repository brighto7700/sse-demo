# 🚀 Go SSE: Blazing Fast Real-Time API

> Stop using WebSockets for everything. Here is a lightweight Server-Sent Events (SSE) implementation in standard library Go.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Live Demo](https://img.shields.io/badge/Live-Demo-brightgreen)](https://sse-demo.pxxl.click)

This repository is the companion code for the DEV.to article: **[Your Next Real-Time Feature Probably Doesn't Need WebSockets](https://dev.to/brighto7700/your-next-real-time-feature-probably-doesnt-need-websockets-go-sse-at-500-concurrent-connections-39ne)**.

It demonstrates how to build a unidirectional real-time data stream (a live viewer counter and server clock) using Go and Server-Sent Events. Zero dependencies on the server. Zero npm packages on the client.

**Try the Live Demo:** [https://sse-demo.pxxl.click](https://sse-demo.pxxl.click)

---

## ✨ Features

* **Standard HTTP:** No `ws://` protocol upgrades, completely firewall-friendly.
* **Auto-Reconnect:** Built right into the browser's native `EventSource` API.
* **Stupidly Efficient:** Holds hundreds of concurrent connections on a fraction of the memory a Node.js WebSocket server would use.
* **Proxy-Safe:** Includes the `X-Accel-Buffering: no` header to prevent reverse proxies (like Nginx) from holding your streams hostage.

## 📊 The Benchmark Flex

Because Go's concurrency model is incredibly lightweight, SSE in Go is a cheat code for cheap real-time features. 

Tested with 500 concurrent connections over 30 seconds (`hey -n 500 -c 500 -t 30 http://localhost:8080/events`) on a standard 2-core Linux VM:

| Metric | Result |
| :--- | :--- |
| **Concurrent connections** | 500 |
| **Memory usage (RSS)** | ~18 MB |
| **CPU at steady state** | ~2% |
| **Avg event latency** | < 2ms |
| **Connection drops** | 0 |

*(For comparison, a naive Node.js WebSocket server at this concurrency sits around 80–110MB RSS).*

## 🚀 Quick Start

Want to run this locally? It takes about 3 seconds.

1. Clone the repo:
   ```bash
   git clone [https://github.com/brighto7700/sse-demo.git](https://github.com/brighto7700/sse-demo.git)
   cd sse-demo
   ```

2. Run the Go server:
   ```bash
   go run main.go
   ```

3. Open your browser to `http://localhost:8080`. 
   *(Pro tip: Open an Incognito window next to it and watch the viewer count sync instantly!)*

## 🧪 Testing the Stream (Terminal)

You don't even need a browser to test SSE. You can watch the raw data stream right in your terminal:

```bash
curl -N http://localhost:8080/events
```
*(The `-N` flag tells curl not to buffer the output!)*

## 🚀 Deployment Notes

SSE requires a **long-lived server process**. Serverless functions (like Vercel or AWS Lambda) have hard execution timeouts and will violently kill your SSE connections after 10–30 seconds. 

Deploy this to a persistent container or VPS platform like [Pxxl App](https://pxxl.app), Fly.io, Railway, or Render. 
