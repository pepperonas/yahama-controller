package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// IP regex validation
// ---------------------------------------------------------------------------

func TestIPRegex_ValidAddresses(t *testing.T) {
	valid := []string{
		"192.168.1.1",
		"192.168.178.31",
		"10.0.0.1",
		"172.16.254.1",
		"255.255.255.255",
		"0.0.0.0",
		"1.2.3.4",
	}
	for _, ip := range valid {
		if !ipRe.MatchString(ip) {
			t.Errorf("expected %q to be a valid IP", ip)
		}
	}
}

func TestIPRegex_InvalidAddresses(t *testing.T) {
	invalid := []string{
		"",
		"256.0.0.1",
		"192.168.1",
		"192.168.1.1.1",
		"abc.def.ghi.jkl",
		"192.168.1.300",
		"not-an-ip",
		"192.168. 1.1",
		"::1",               // IPv6
		"192.168.1.1:8080",  // with port
	}
	for _, ip := range invalid {
		if ipRe.MatchString(ip) {
			t.Errorf("expected %q to be an invalid IP", ip)
		}
	}
}

// ---------------------------------------------------------------------------
// Config load / save round-trip
// ---------------------------------------------------------------------------

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	origFile := configFile
	origIP := receiverIP
	defer func() {
		configFile = origFile
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()

	configFile = filepath.Join(dir, "test-config.json")

	saveConfig("10.0.0.42")
	if got := getIP(); got != "10.0.0.42" {
		t.Fatalf("getIP() = %q, want %q", got, "10.0.0.42")
	}

	// Reset in-memory and reload from disk
	mu.Lock()
	receiverIP = ""
	mu.Unlock()
	loadConfig()
	if got := getIP(); got != "10.0.0.42" {
		t.Fatalf("after reload getIP() = %q, want %q", got, "10.0.0.42")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()
	origFile := configFile
	origIP := receiverIP
	defer func() {
		configFile = origFile
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()

	configFile = filepath.Join(dir, "nonexistent.json")
	mu.Lock()
	receiverIP = "previous"
	mu.Unlock()

	loadConfig() // must not panic or change receiverIP
	if got := getIP(); got != "previous" {
		t.Fatalf("loadConfig on missing file must not mutate receiverIP, got %q", got)
	}
}

func TestLoadConfig_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	origFile := configFile
	origIP := receiverIP
	defer func() {
		configFile = origFile
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()

	configFile = filepath.Join(dir, "bad.json")
	_ = os.WriteFile(configFile, []byte("not json!!!"), 0o644)

	mu.Lock()
	receiverIP = "stable"
	mu.Unlock()

	loadConfig() // must not panic or change receiverIP
	if got := getIP(); got != "stable" {
		t.Fatalf("loadConfig with corrupt file must not mutate receiverIP, got %q", got)
	}
}

func TestSaveConfig_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	origFile := configFile
	origIP := receiverIP
	defer func() {
		configFile = origFile
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()

	configFile = filepath.Join(dir, "cfg.json")
	saveConfig("172.16.0.5")

	b, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("config file not written: %v", err)
	}
	var c config
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatalf("written file is not valid JSON: %v", err)
	}
	if c.ReceiverIP != "172.16.0.5" {
		t.Fatalf("JSON ReceiverIP = %q, want %q", c.ReceiverIP, "172.16.0.5")
	}
}

// ---------------------------------------------------------------------------
// GET /api/health
// ---------------------------------------------------------------------------

func TestHealthEndpoint_StatusOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}

func TestHealthEndpoint_JSONBody(t *testing.T) {
	origIP := receiverIP
	mu.Lock()
	receiverIP = "192.168.1.99"
	mu.Unlock()
	defer func() {
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	health(rec, req)

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("body is not JSON: %v", err)
	}
	if body["status"] != "OK" {
		t.Errorf("status = %v, want OK", body["status"])
	}
	if body["receiverIP"] != "192.168.1.99" {
		t.Errorf("receiverIP = %v, want 192.168.1.99", body["receiverIP"])
	}
	if _, ok := body["timestamp"]; !ok {
		t.Error("timestamp field missing from health response")
	}
}

// ---------------------------------------------------------------------------
// POST /api/set-receiver-ip
// ---------------------------------------------------------------------------

func TestSetIP_ValidIP(t *testing.T) {
	dir := t.TempDir()
	origFile := configFile
	origIP := receiverIP
	defer func() {
		configFile = origFile
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()
	configFile = filepath.Join(dir, "cfg.json")

	body := strings.NewReader(`{"ip":"10.1.2.3"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/set-receiver-ip", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	setIP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("setIP with valid IP returned %d, want 200", rec.Code)
	}
	var resp map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp["success"] != true {
		t.Errorf("success = %v, want true", resp["success"])
	}
	if resp["ip"] != "10.1.2.3" {
		t.Errorf("ip = %v, want 10.1.2.3", resp["ip"])
	}
	if getIP() != "10.1.2.3" {
		t.Errorf("getIP() = %q after setIP, want 10.1.2.3", getIP())
	}
}

func TestSetIP_InvalidIP(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"bad format", `{"ip":"999.999.999.999"}`},
		{"empty ip",   `{"ip":""}`},
		{"no ip field", `{}`},
		{"hostname",   `{"ip":"yamaha.local"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/set-receiver-ip",
				strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			setIP(rec, req)
			if rec.Code == http.StatusOK {
				t.Errorf("expected non-200, got 200 for body %q", tc.body)
			}
		})
	}
}

func TestSetIP_MalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/set-receiver-ip",
		strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	setIP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("malformed JSON → %d, want 400", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// CORS middleware
// ---------------------------------------------------------------------------

func TestCORS_HeadersPresent(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := cors(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("ACAO = %q, want *", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("Access-Control-Allow-Methods header missing")
	}
}

func TestCORS_OptionsPreflightReturns204(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // must never reach here for OPTIONS
	})
	handler := cors(inner)

	req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS preflight = %d, want 204", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Static file serving — path-traversal guard
// ---------------------------------------------------------------------------

func TestStatic_PathTraversalGuard(t *testing.T) {
	origBase := baseDir
	defer func() { baseDir = origBase }()
	baseDir = t.TempDir()

	malicious := []string{
		"/../etc/passwd",
		"/../../etc/shadow",
	}
	for _, path := range malicious {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		static(rec, req)
		if rec.Code == http.StatusOK {
			t.Errorf("path-traversal %q must not return 200", path)
		}
	}
}

func TestStatic_IndexFallback(t *testing.T) {
	dir := t.TempDir()
	origBase := baseDir
	defer func() { baseDir = origBase }()
	baseDir = dir

	_ = os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>ok</html>"), 0o644)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	static(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("/ → %d, want 200", rec.Code)
	}
}

func TestStatic_MissingFile404(t *testing.T) {
	origBase := baseDir
	defer func() { baseDir = origBase }()
	baseDir = t.TempDir()

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist.txt", nil)
	rec := httptest.NewRecorder()
	static(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing file → %d, want 404", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Receiver proxy — no IP configured
// ---------------------------------------------------------------------------

func TestReceiverProxy_NoIPSet(t *testing.T) {
	origIP := receiverIP
	mu.Lock()
	receiverIP = ""
	mu.Unlock()
	defer func() {
		mu.Lock()
		receiverIP = origIP
		mu.Unlock()
	}()

	req := httptest.NewRequest(http.MethodGet, "/api/receiver/YamahaRemoteControl", nil)
	rec := httptest.NewRecorder()
	receiverProxy(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("proxy with no IP → %d, want 500", rec.Code)
	}
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["error"] == "" {
		t.Error("expected error field in 500 response body")
	}
}
