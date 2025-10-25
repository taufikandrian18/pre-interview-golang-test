package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	cache  Cache
	closer Closer
}

type CacheResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewServer(cache Cache, closer Closer) *Server {
	return &Server{
		cache:  cache,
		closer: closer,
	}
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed. Use POST", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		s.sendError(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	var value interface{}
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var jsonValue interface{}
		if err := json.NewDecoder(r.Body).Decode(&jsonValue); err != nil {
			s.sendError(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		value = jsonValue
	} else {
		// Handle form data or plain text
		if err := r.ParseForm(); err != nil {
			s.sendError(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		valueStr := r.FormValue("value")
		if valueStr == "" {
			s.sendError(w, "Value parameter is required", http.StatusBadRequest)
			return
		}
		value = valueStr
	}

	err := s.cache.Set(key, value)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.sendSuccess(w, fmt.Sprintf("Successfully set key '%s'", key), value)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed. Use GET", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		s.sendError(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	value, exists, err := s.cache.Get(key)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		s.sendError(w, fmt.Sprintf("Key '%s' not found", key), http.StatusNotFound)
		return
	}

	s.sendSuccess(w, fmt.Sprintf("Found key '%s'", key), value)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		s.sendError(w, "Method not allowed. Use DELETE", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		s.sendError(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	err := s.cache.Delete(key)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.sendSuccess(w, fmt.Sprintf("Successfully deleted key '%s'", key), nil)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed. Use GET", http.StatusMethodNotAllowed)
		return
	}

	stats := map[string]interface{}{
		"message": "Cache server is running",
		"time":    time.Now().Format(time.RFC3339),
		"type":    fmt.Sprintf("%T", s.cache),
	}

	s.sendSuccess(w, "Cache stats", stats)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.sendSuccess(w, "OK", map[string]string{"status": "healthy"})
}

func (s *Server) sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := CacheResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := CacheResponse{
		Success: false,
		Error:   message,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/cache/set", s.handleSet)
	mux.HandleFunc("/api/cache/get", s.handleGet)
	mux.HandleFunc("/api/cache/delete", s.handleDelete)
	mux.HandleFunc("/api/cache/stats", s.handleStats)
	mux.HandleFunc("/health", s.handleHealth)

	// Root endpoint with usage instructions
	mux.HandleFunc("/", s.handleRoot)

	return mux
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Cache API Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .method { color: #007acc; font-weight: bold; }
        code { background: #e8e8e8; padding: 2px 4px; border-radius: 3px; }
        pre { background: #f8f8f8; padding: 15px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>Cache API Server</h1>
    <p>Interactive API for testing the cache implementation</p>

    <h2>Available Endpoints</h2>

    <div class="endpoint">
        <span class="method">POST</span> <code>/api/cache/set?key={key}</code>
        <p>Set a value in the cache. Send value in request body.</p>
        <p><strong>Examples:</strong></p>
        <pre>curl -X POST "http://localhost:8080/api/cache/set?key=name" \
     -H "Content-Type: application/json" \
     -d '"John Doe"'</pre>
        <pre>curl -X POST "http://localhost:8080/api/cache/set?key=age" \
     -d "value=25"</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> <code>/api/cache/get?key={key}</code>
        <p>Get a value from the cache.</p>
        <p><strong>Example:</strong></p>
        <pre>curl "http://localhost:8080/api/cache/get?key=name"</pre>
    </div>

    <div class="endpoint">
        <span class="method">DELETE</span> <code>/api/cache/delete?key={key}</code>
        <p>Delete a key from the cache.</p>
        <p><strong>Example:</strong></p>
        <pre>curl -X DELETE "http://localhost:8080/api/cache/delete?key=name"</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> <code>/api/cache/stats</code>
        <p>Get cache statistics and information.</p>
        <p><strong>Example:</strong></p>
        <pre>curl "http://localhost:8080/api/cache/stats"</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> <code>/health</code>
        <p>Health check endpoint.</p>
        <pre>curl "http://localhost:8080/health"</pre>
    </div>

    <h2>Quick Test Commands</h2>
    <p>Copy and paste these commands to test the cache:</p>
    <pre>
# Set some values
curl -X POST "http://localhost:8080/api/cache/set?key=user" -H "Content-Type: application/json" -d '"Alice"'
curl -X POST "http://localhost:8080/api/cache/set?key=age" -d "value=30"

# Get values
curl "http://localhost:8080/api/cache/get?key=user"
curl "http://localhost:8080/api/cache/get?key=age"

# Get stats
curl "http://localhost:8080/api/cache/stats"

# Delete a key
curl -X DELETE "http://localhost:8080/api/cache/delete?key=age"

# Try to get deleted key
curl "http://localhost:8080/api/cache/get?key=age"
    </pre>

    <p><strong>Server is running with:</strong> ` + fmt.Sprintf("%T", s.cache) + `</p>
</body>
</html>`

	w.Write([]byte(html))
}

func runServer() {
	var cache Cache
	var closer Closer

	// You can switch between cache types here
	useTTL := true

	if useTTL {
		ttlCache, err := NewTTLCache(30 * time.Second) // 30 second TTL for demo
		if err != nil {
			log.Fatal("Failed to create TTL cache:", err)
		}
		cache = ttlCache
		closer = ttlCache
	} else {
		simpleCache := NewSimpleCache()
		cache = simpleCache
		closer = nil
	}

	if closer != nil {
		defer closer.Close()
	}

	server := NewServer(cache, closer)
	mux := server.setupRoutes()

	port := "8080"
	fmt.Printf("🚀 Cache API Server starting on http://localhost:%s\n", port)
	fmt.Printf("📖 Open http://localhost:%s in your browser for API documentation\n", port)
	fmt.Printf("🛑 Press Ctrl+C to stop the server\n\n")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
