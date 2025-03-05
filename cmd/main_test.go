package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a test server
	router := setupRouter()
	
	// Create a request to the ping endpoint
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK (200), got %v", w.Code)
	}
	
	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	
	// Check the response content
	if response["message"] != "pong" {
		t.Errorf("Expected message 'pong', got %v", response["message"])
	}
}

func TestSetupServer(t *testing.T) {
	// Save original environment variables
	origAppID := os.Getenv("AGORA_APP_ID")
	origAppCert := os.Getenv("AGORA_APP_CERTIFICATE")
	origCustID := os.Getenv("AGORA_CUSTOMER_ID")
	origCustSecret := os.Getenv("AGORA_CUSTOMER_SECRET")
	origBaseURL := os.Getenv("AGORA_CONVO_AI_BASE_URL")
	origAgentUID := os.Getenv("AGENT_UID")
	origLLMModel := os.Getenv("LLM_MODEL")
	origLLMURL := os.Getenv("LLM_URL")
	origLLMToken := os.Getenv("LLM_TOKEN")
	origTTSVendor := os.Getenv("TTS_VENDOR")
	origMSTTSKey := os.Getenv("MICROSOFT_TTS_KEY")
	origMSTTSRegion := os.Getenv("MICROSOFT_TTS_REGION")
	origMSTTSVoice := os.Getenv("MICROSOFT_TTS_VOICE_NAME")
	origMSTTSRate := os.Getenv("MICROSOFT_TTS_RATE")
	origMSTTSVolume := os.Getenv("MICROSOFT_TTS_VOLUME")
	
	// Restore environment variables after test
	defer func() {
		os.Setenv("AGORA_APP_ID", origAppID)
		os.Setenv("AGORA_APP_CERTIFICATE", origAppCert)
		os.Setenv("AGORA_CUSTOMER_ID", origCustID)
		os.Setenv("AGORA_CUSTOMER_SECRET", origCustSecret)
		os.Setenv("AGORA_CONVO_AI_BASE_URL", origBaseURL)
		os.Setenv("AGENT_UID", origAgentUID)
		os.Setenv("LLM_MODEL", origLLMModel)
		os.Setenv("LLM_URL", origLLMURL)
		os.Setenv("LLM_TOKEN", origLLMToken)
		os.Setenv("TTS_VENDOR", origTTSVendor)
		os.Setenv("MICROSOFT_TTS_KEY", origMSTTSKey)
		os.Setenv("MICROSOFT_TTS_REGION", origMSTTSRegion)
		os.Setenv("MICROSOFT_TTS_VOICE_NAME", origMSTTSVoice)
		os.Setenv("MICROSOFT_TTS_RATE", origMSTTSRate)
		os.Setenv("MICROSOFT_TTS_VOLUME", origMSTTSVolume)
	}()
	
	// Set test environment variables
	os.Setenv("AGORA_APP_ID", "test-app-id")
	os.Setenv("AGORA_APP_CERTIFICATE", "test-app-cert")
	os.Setenv("AGORA_CUSTOMER_ID", "test-customer-id")
	os.Setenv("AGORA_CUSTOMER_SECRET", "test-customer-secret")
	os.Setenv("AGORA_CONVO_AI_BASE_URL", "https://api.example.com")
	os.Setenv("AGENT_UID", "123456")
	os.Setenv("LLM_MODEL", "gpt-3.5-turbo")
	os.Setenv("LLM_URL", "https://api.openai.com/v1/chat/completions")
	os.Setenv("LLM_TOKEN", "test-llm-token")
	os.Setenv("TTS_VENDOR", "microsoft")
	os.Setenv("MICROSOFT_TTS_KEY", "test-ms-key")
	os.Setenv("MICROSOFT_TTS_REGION", "eastus")
	os.Setenv("MICROSOFT_TTS_VOICE_NAME", "en-US-AriaNeural")
	os.Setenv("MICROSOFT_TTS_RATE", "1.0")
	os.Setenv("MICROSOFT_TTS_VOLUME", "1.0")
	
	// Test server setup
	server := setupServer()
	
	// Check that server is not nil
	if server == nil {
		t.Fatal("setupServer() returned nil")
	}
	
	// Check server address
	expectedAddr := ":8080" // Default port
	if server.Addr != expectedAddr {
		t.Errorf("Expected server address %s, got %s", expectedAddr, server.Addr)
	}
	
	// Test with custom port
	os.Setenv("PORT", "9090")
	server = setupServer()
	
	expectedAddr = ":9090"
	if server.Addr != expectedAddr {
		t.Errorf("Expected server address %s, got %s", expectedAddr, server.Addr)
	}
}

// Helper function to initialize router for testing
func setupRouter() http.Handler {
	// Set test environment variables
	os.Setenv("AGORA_APP_ID", "test-app-id")
	os.Setenv("AGORA_APP_CERTIFICATE", "test-app-cert")
	os.Setenv("AGORA_CUSTOMER_ID", "test-customer-id")
	os.Setenv("AGORA_CUSTOMER_SECRET", "test-customer-secret")
	os.Setenv("AGORA_CONVO_AI_BASE_URL", "https://api.example.com")
	os.Setenv("AGENT_UID", "123456")
	os.Setenv("LLM_MODEL", "gpt-3.5-turbo")
	os.Setenv("LLM_URL", "https://api.openai.com/v1/chat/completions")
	os.Setenv("LLM_TOKEN", "test-llm-token")
	os.Setenv("TTS_VENDOR", "microsoft")
	os.Setenv("MICROSOFT_TTS_KEY", "test-ms-key")
	os.Setenv("MICROSOFT_TTS_REGION", "eastus")
	os.Setenv("MICROSOFT_TTS_VOICE_NAME", "en-US-AriaNeural")
	os.Setenv("MICROSOFT_TTS_RATE", "1.0")
	os.Setenv("MICROSOFT_TTS_VOLUME", "1.0")
	
	server := setupServer()
	return server.Handler
}