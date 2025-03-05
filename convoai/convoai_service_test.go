package convoai

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AgoraIO-Community/convo-ai-go-server/token_service"
	"github.com/gin-gonic/gin"
)

// MockTokenService mocks the token service for testing
type MockTokenService struct {
	AppID          string
	AppCertificate string
	MockToken      string
}

// GenRtcToken implements the token generation method for MockTokenService
func (m *MockTokenService) GenRtcToken(req token_service.TokenRequest) (string, error) {
	return "mock-rtc-token", nil
}

// Create a new mock token service
func NewMockTokenService() *MockTokenService {
	return &MockTokenService{
		AppID:          "test-app-id",
		AppCertificate: "test-app-cert",
		MockToken:      "mock-token-12345",
	}
}

// Mock HTTP client for testing
type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Mock the success response
	mockResponse := `{"agent_id": "test-agent-123"}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
	}, nil
}

// NewTestConvoAIService creates a ConvoAIService instance for testing
func NewTestConvoAIService() *ConvoAIService {
	config := &ConvoAIConfig{
		AppID:          "test-app-id",
		AppCertificate: "test-app-cert",
		CustomerID:     "test-customer-id",
		CustomerSecret: "test-customer-secret",
		BaseURL:        "https://api.example.com",
		AgentUID:       "123456",
		LLMModel:       "gpt-3.5-turbo",
		LLMURL:         "https://api.openai.com/v1/chat/completions",
		LLMToken:       "test-llm-token",
		TTSVendor:      "microsoft",
		MicrosoftTTS: &MicrosoftTTSConfig{
			Key:       "test-ms-key",
			Region:    "eastus",
			VoiceName: "en-US-AriaNeural",
			Rate:      "1.0",
			Volume:    "1.0",
		},
	}

	mockTokenService := NewMockTokenService()
	
	return &ConvoAIService{
		config:       config,
		tokenService: mockTokenService,
	}
}

func TestValidateInviteRequest(t *testing.T) {
	service := NewTestConvoAIService()

	tests := []struct {
		name    string
		request InviteAgentRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			request: InviteAgentRequest{
				RequesterID: "123",
				ChannelName: "test-channel",
			},
			wantErr: false,
		},
		{
			name: "Missing requester_id",
			request: InviteAgentRequest{
				ChannelName: "test-channel",
			},
			wantErr: true,
		},
		{
			name: "Missing channel_name",
			request: InviteAgentRequest{
				RequesterID: "123",
			},
			wantErr: true,
		},
		{
			name: "Channel name too short",
			request: InviteAgentRequest{
				RequesterID: "123",
				ChannelName: "te",
			},
			wantErr: true,
		},
		{
			name: "Channel name too long",
			request: InviteAgentRequest{
				RequesterID: "123",
				ChannelName: strings.Repeat("a", 65),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateInviteRequest(&tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInviteRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRemoveRequest(t *testing.T) {
	service := NewTestConvoAIService()

	tests := []struct {
		name    string
		request RemoveAgentRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			request: RemoveAgentRequest{
				AgentID: "agent-123",
			},
			wantErr: false,
		},
		{
			name:    "Missing agent_id",
			request: RemoveAgentRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateRemoveRequest(&tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRemoveRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTTSConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *ConvoAIConfig
		wantErr   bool
		wantVendor TTSVendor
	}{
		{
			name: "Microsoft TTS",
			config: &ConvoAIConfig{
				TTSVendor: "microsoft",
				MicrosoftTTS: &MicrosoftTTSConfig{
					Key:       "test-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
					Rate:      "1.0",
					Volume:    "1.0",
				},
			},
			wantErr:   false,
			wantVendor: TTSVendorMicrosoft,
		},
		{
			name: "ElevenLabs TTS",
			config: &ConvoAIConfig{
				TTSVendor: "elevenlabs",
				ElevenLabsTTS: &ElevenLabsTTSConfig{
					Key:     "test-key",
					VoiceID: "voice-id",
					ModelID: "model-id",
				},
			},
			wantErr:   false,
			wantVendor: TTSVendorElevenLabs,
		},
		{
			name: "Microsoft TTS missing config",
			config: &ConvoAIConfig{
				TTSVendor: "microsoft",
			},
			wantErr: true,
		},
		{
			name: "ElevenLabs TTS missing config",
			config: &ConvoAIConfig{
				TTSVendor: "elevenlabs",
			},
			wantErr: true,
		},
		{
			name: "Unsupported vendor",
			config: &ConvoAIConfig{
				TTSVendor: "unsupported",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ConvoAIService{
				config: tt.config,
			}

			config, err := service.getTTSConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("getTTSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && config.Vendor != tt.wantVendor {
				t.Errorf("getTTSConfig() vendor = %v, want %v", config.Vendor, tt.wantVendor)
			}
		})
	}
}

func TestIsStringUID(t *testing.T) {
	tests := []struct {
		name  string
		uid   string
		want  bool
	}{
		{
			name: "Numeric UID",
			uid:  "12345",
			want: false,
		},
		{
			name: "String UID",
			uid:  "user123",
			want: true,
		},
		{
			name: "Mixed UID",
			uid:  "123abc",
			want: true,
		},
		{
			name: "Empty UID",
			uid:  "",
			want: false,
		},
		{
			name: "Zero UID",
			uid:  "0",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStringUID(tt.uid); got != tt.want {
				t.Errorf("isStringUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInviteAgentValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewTestConvoAIService()

	// Setup
	router := gin.New()
	service.RegisterRoutes(router)

	// Test validation cases only
	tests := []struct {
		name           string
		requestBody    string
		wantStatusCode int
	}{
		{
			name:           "Invalid JSON",
			requestBody:    `{"requester_id": "123"`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Missing required field",
			requestBody:    `{"channel_name": "test-channel"}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/agent/invite", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}
		})
	}
}

func TestGetBasicAuth(t *testing.T) {
	service := NewTestConvoAIService()
	
	auth := service.getBasicAuth()
	
	// Expected value: Basic dGVzdC1jdXN0b21lci1pZDp0ZXN0LWN1c3RvbWVyLXNlY3JldA==
	expected := "Basic dGVzdC1jdXN0b21lci1pZDp0ZXN0LWN1c3RvbWVyLXNlY3JldA=="
	
	if auth != expected {
		t.Errorf("getBasicAuth() = %v, want %v", auth, expected)
	}
}