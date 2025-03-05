package http_headers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(allowOrigin string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	headers := NewHttpHeaders(allowOrigin)
	
	router.Use(headers.NoCache())
	router.Use(headers.CORShttpHeaders())
	router.Use(headers.Timestamp())
	
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})
	
	return router
}

func TestNewHttpHeaders(t *testing.T) {
	headers := NewHttpHeaders("http://example.com")
	
	if headers.AllowOrigin != "http://example.com" {
		t.Errorf("NewHttpHeaders() AllowOrigin = %v, want %v", headers.AllowOrigin, "http://example.com")
	}
}

func TestNoCache(t *testing.T) {
	router := setupRouter("*")
	
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	
	router.ServeHTTP(resp, req)
	
	if resp.Header().Get("Cache-Control") != "private, no-cache, no-store, must-revalidate" {
		t.Errorf("NoCache() Cache-Control = %v, want %v", 
			resp.Header().Get("Cache-Control"), 
			"private, no-cache, no-store, must-revalidate")
	}
	
	if resp.Header().Get("Expires") != "-1" {
		t.Errorf("NoCache() Expires = %v, want %v", resp.Header().Get("Expires"), "-1")
	}
	
	if resp.Header().Get("Pragma") != "no-cache" {
		t.Errorf("NoCache() Pragma = %v, want %v", resp.Header().Get("Pragma"), "no-cache")
	}
}

func TestCORSHeadersWildcard(t *testing.T) {
	router := setupRouter("*")
	
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	resp := httptest.NewRecorder()
	
	router.ServeHTTP(resp, req)
	
	if resp.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("CORShttpHeaders() Access-Control-Allow-Origin = %v, want %v", 
			resp.Header().Get("Access-Control-Allow-Origin"), 
			"http://example.com")
	}
	
	if resp.Code != http.StatusOK {
		t.Errorf("CORShttpHeaders() status code = %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestCORSHeadersAllowed(t *testing.T) {
	router := setupRouter("http://example.com,http://allowed.com")
	
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	resp := httptest.NewRecorder()
	
	router.ServeHTTP(resp, req)
	
	if resp.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("CORShttpHeaders() Access-Control-Allow-Origin = %v, want %v", 
			resp.Header().Get("Access-Control-Allow-Origin"), 
			"http://example.com")
	}
	
	if resp.Code != http.StatusOK {
		t.Errorf("CORShttpHeaders() status code = %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestCORSHeadersNotAllowed(t *testing.T) {
	router := setupRouter("http://example.com")
	
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://not-allowed.com")
	resp := httptest.NewRecorder()
	
	router.ServeHTTP(resp, req)
	
	if resp.Code != http.StatusForbidden {
		t.Errorf("CORShttpHeaders() status code = %v, want %v", resp.Code, http.StatusForbidden)
	}
}

func TestCORSHeadersOptions(t *testing.T) {
	router := setupRouter("*")
	
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	resp := httptest.NewRecorder()
	
	router.ServeHTTP(resp, req)
	
	if resp.Code != http.StatusNoContent {
		t.Errorf("CORShttpHeaders() OPTIONS status code = %v, want %v", resp.Code, http.StatusNoContent)
	}
}

func TestTimestamp(t *testing.T) {
	router := setupRouter("*")
	
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	
	router.ServeHTTP(resp, req)
	
	// Parse the timestamp from the response header
	timestamp := resp.Header().Get("X-Timestamp")
	if timestamp == "" {
		t.Errorf("Timestamp() X-Timestamp header not set")
		return
	}
	
	// Parse the timestamp
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("Timestamp() X-Timestamp format error: %v", err)
	}
	
	// We won't check the exact time range since it depends on execution time
	// and can cause flaky tests. Just verify it's a valid timestamp.
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name        string
		allowOrigin string
		origin      string
		want        bool
	}{
		{
			name:        "Wildcard allows any origin",
			allowOrigin: "*",
			origin:      "http://example.com",
			want:        true,
		},
		{
			name:        "Specific origin allowed",
			allowOrigin: "http://example.com",
			origin:      "http://example.com",
			want:        true,
		},
		{
			name:        "Multiple origins - allowed",
			allowOrigin: "http://example.com,http://allowed.com",
			origin:      "http://allowed.com",
			want:        true,
		},
		{
			name:        "Origin not allowed",
			allowOrigin: "http://example.com",
			origin:      "http://not-allowed.com",
			want:        false,
		},
		{
			name:        "Empty origin not allowed",
			allowOrigin: "http://example.com",
			origin:      "",
			want:        false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := NewHttpHeaders(tt.allowOrigin)
			if got := headers.isOriginAllowed(tt.origin); got != tt.want {
				t.Errorf("isOriginAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}