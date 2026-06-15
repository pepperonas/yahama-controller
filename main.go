// Yamaha Receiver Control — a tiny Go replacement for the old Node/Express+PM2
// app. Serves the PWA, exposes /api/health and /api/set-receiver-ip, and reverse-
// proxies /api/receiver/* to the (runtime-configurable) receiver IP. Drop-in on
// port 5001 — the nginx /app/yamaha/ and /proxy/yamaha/api/ routes are unchanged.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	mu         sync.RWMutex
	receiverIP string
	baseDir    string
	configFile string
	ipRe       = regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
)

type config struct {
	ReceiverIP string `json:"receiverIP"`
}

func loadConfig() {
	if b, err := os.ReadFile(configFile); err == nil {
		var c config
		if json.Unmarshal(b, &c) == nil {
			mu.Lock()
			receiverIP = c.ReceiverIP
			mu.Unlock()
		}
	}
}

func saveConfig(ip string) {
	mu.Lock()
	receiverIP = ip
	mu.Unlock()
	b, _ := json.MarshalIndent(config{ReceiverIP: ip}, "", "  ")
	_ = os.WriteFile(configFile, b, 0o644)
}

func getIP() string {
	mu.RLock()
	defer mu.RUnlock()
	return receiverIP
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"status":     "OK",
		"receiverIP": getIP(),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

func setIP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IP string `json:"ip"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil || body.IP == "" {
		writeJSON(w, 400, map[string]string{"error": "IP address is required"})
		return
	}
	if !ipRe.MatchString(body.IP) {
		writeJSON(w, 400, map[string]string{"error": "Invalid IP address format"})
		return
	}
	saveConfig(body.IP)
	log.Printf("Receiver IP set to: %s", body.IP)
	writeJSON(w, 200, map[string]any{"success": true, "ip": body.IP})
}

// reverse-proxy /api/receiver/* → http://<receiverIP>/* (prefix stripped)
func receiverProxy(w http.ResponseWriter, r *http.Request) {
	ip := getIP()
	if ip == "" {
		writeJSON(w, 500, map[string]string{"error": "Failed to connect to receiver", "message": "no receiver IP set"})
		return
	}
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = ip
			req.Host = ip
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/receiver")
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			writeJSON(w, 500, map[string]string{"error": "Failed to connect to receiver", "message": err.Error()})
		},
	}
	proxy.ServeHTTP(w, r)
}

// serve PWA: try public/<path>, then <root>/<path>; "/" → index.html
func static(w http.ResponseWriter, r *http.Request) {
	clean := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if clean == "." || clean == "" {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, filepath.Join(baseDir, "index.html"))
		return
	}
	if strings.HasPrefix(clean, "..") { // path-traversal guard
		http.NotFound(w, r)
		return
	}
	for _, c := range []string{filepath.Join(baseDir, "public", clean), filepath.Join(baseDir, clean)} {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			if strings.HasSuffix(c, ".html") || strings.HasSuffix(c, "manifest.json") || strings.HasSuffix(c, "sw.js") {
				w.Header().Set("Cache-Control", "no-cache")
			}
			http.ServeFile(w, r, c)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}
	baseDir = os.Getenv("YAMAHA_DIR")
	if baseDir == "" {
		baseDir, _ = os.Getwd()
	}
	configFile = filepath.Join(baseDir, "receiver-config.json")
	loadConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", health)
	mux.HandleFunc("/api/set-receiver-ip", setIP)
	mux.HandleFunc("/api/receiver/", receiverProxy)
	mux.HandleFunc("/api/receiver", receiverProxy)
	mux.HandleFunc("/", static)

	log.Printf("Yamaha Receiver Control (Go) on :%s  receiverIP=%q  dir=%s", port, getIP(), baseDir)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, cors(mux)))
}
