# Build an Agora Conversational AI Service using Golang

Conversational AI is revolutionizing how people interact with artificial intelligence. Instead of carefully crafting text prompts, users can have natural, real-time voice conversations with AI agents. This opens exciting opportunities for more intuitive and efficient interactions.

Many developers have already invested significant time building custom LLM workflows for text-based agents. Agora's Conversational AI Engine allows you to connect these existing workflows to an Agora channel, enabling real-time voice conversations without abandoning your current AI infrastructure.

In this guide, I'll walk you through building a Go server that handles the connection between your users and Agora's Conversational AI. By the end, you'll have a production-ready backend that can power voice-based AI conversations for your applications.

## Prerequisites

Before getting started, make sure you have:

- Go (version 1.18 or higher)
- Basic knowledge of Go and the Gin framework
- [An Agora account](https://console.agora.io/) - _the first 10k minutes each month are free_
- Conversational AI service [activated on your AppID](https://console.agora.io/)

## Project Setup

Let's start by setting up our Go project with the necessary dependencies. First, create a new directory and initialize a Go module:

```bash
mkdir agora-convo-ai-go-server
cd agora-convo-ai-go-server
go mod init github.com/AgoraIO-Community/convo-ai-go-server
```

Next, we'll add the key dependencies for our server:

```bash
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go get github.com/AgoraIO-Community/go-tokenbuilder
```

Create the initial directory structure, and as we go through the guide, we'll fill these directories with the files we need.

```bash
mkdir -p convoai token_service http_headers validation
touch .env
```

Your project directory should now have a structure like this:

```
├── convoai/
├── token_service/
├── http_headers/
├── validation/
├── .env
├── go.mod
├── go.sum
```

## Server Entry Point

Start by setting up the main application file, which will be the entry point for our server. We'll then load the environment variables, set up the configuration, and initialize the router with the appropriate middleware and routes.

Create the `main.go` file

```bash
touch main.go
```

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadConfig() (*convoai.ConvoAIConfig, error) {
	config := &convoai.ConvoAIConfig{
		// Agora Configuration
		AppID:          os.Getenv("AGORA_APP_ID"),
		AppCertificate: os.Getenv("AGORA_APP_CERTIFICATE"),
		CustomerID:     os.Getenv("AGORA_CUSTOMER_ID"),
		CustomerSecret: os.Getenv("AGORA_CUSTOMER_SECRET"),
		BaseURL:        os.Getenv("AGORA_CONVO_AI_BASE_URL"),
		AgentUID:       os.Getenv("AGENT_UID"),

		// LLM Configuration
		LLMModel: os.Getenv("LLM_MODEL"),
		LLMURL:   os.Getenv("LLM_URL"),
		LLMToken: os.Getenv("LLM_TOKEN"),

		// TTS Configuration
		TTSVendor: os.Getenv("TTS_VENDOR"),
	}

	// Microsoft TTS Configuration
	if msKey := os.Getenv("MICROSOFT_TTS_KEY"); msKey != "" {
		config.MicrosoftTTS = &convoai.MicrosoftTTSConfig{
			Key:       msKey,
			Region:    os.Getenv("MICROSOFT_TTS_REGION"),
			VoiceName: os.Getenv("MICROSOFT_TTS_VOICE_NAME"),
			Rate:      os.Getenv("MICROSOFT_TTS_RATE"),
			Volume:    os.Getenv("MICROSOFT_TTS_VOLUME"),
		}
	}

	// ElevenLabs TTS Configuration
	if elKey := os.Getenv("ELEVENLABS_API_KEY"); elKey != "" {
		config.ElevenLabsTTS = &convoai.ElevenLabsTTSConfig{
			Key:     elKey,
			VoiceID: os.Getenv("ELEVENLABS_VOICE_ID"),
			ModelID: os.Getenv("ELEVENLABS_MODEL_ID"),
		}
	}

	// Modalities Configuration
	config.InputModalities = os.Getenv("INPUT_MODALITIES")
	config.OutputModalities = os.Getenv("OUTPUT_MODALITIES")

	return config, nil
}

func setupServer() *http.Server {
	log.Println("Starting setupServer")
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file. Using existing environment variables.")
	}

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

    // TODO: Validate environment configuration

	// Server Configuration
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// CORS Configuration
	corsAllowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")

	// Set up router with headers
	router := gin.Default()
    //TODO: Register headers

	// TODO: Initialize services & register routes

	// Register healthcheck route
	router.GET("/ping", Ping)

	// Configure and start the HTTP server
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: router,
	}

	log.Println("Server setup completed")
	log.Println("- listening on port", serverPort)
	return server
}

func main() {
	server := setupServer()

	// Start the server in a separate goroutine to handle graceful shutdown.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}

	}()

	// Prepare to handle graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for a shutdown signal.
	<-quit
	log.Println("Shutting down server...")

	// Attempt to gracefully shutdown the server with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// Ping is a handler function that serves as a basic health check endpoint.
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
```

> **Note:** We are loading the PORT from the environment variables, it will default to `8080` if not set in your `.env` file.

Let's test our basic Go server by running:

```bash
go run main.go
```

You should see "Server setup completed" and "- listening on port 8080" in your console. You can now visit <a href="http://localhost:8080/ping">http://localhost:8080/ping</a> to verify the server is working - you should see `{"message": "pong"}` as the response.

To test the server using curl, run:

```bash
curl http://localhost:8080/ping
```

You should see the response: `{"message": "pong"}`.

## Type Definitions

Next, let's define the types needed for our ConvoAI service. Create a file called `convoai-types.go` in the `convoai` directory:

```bash
touch convoai/convoai-types.go
```

Add the following types:

```go
package convoai

// InviteAgentRequest represents the request body for inviting an AI agent
type InviteAgentRequest struct {
	RequesterID      string   `json:"requester_id"`
	ChannelName      string   `json:"channel_name"`
	RtcCodec         *int     `json:"rtc_codec,omitempty"`
	InputModalities  []string `json:"input_modalities,omitempty"`
	OutputModalities []string `json:"output_modalities,omitempty"`
}

// RemoveAgentRequest represents the request body for removing an AI agent
type RemoveAgentRequest struct {
	AgentID string `json:"agent_id"`
}

// TTSVendor represents the text-to-speech vendor type
type TTSVendor string

const (
	TTSVendorMicrosoft  TTSVendor = "microsoft"
	TTSVendorElevenLabs TTSVendor = "elevenlabs"
)

// TTSConfig represents the text-to-speech configuration
type TTSConfig struct {
	Vendor TTSVendor   `json:"vendor"`
	Params interface{} `json:"params"`
}

// AgoraStartRequest represents the request to start a conversation
type AgoraStartRequest struct {
	Name       string     `json:"name"`
	Properties Properties `json:"properties"`
}

// Properties represents the configuration properties for the conversation
type Properties struct {
	Channel          string    `json:"channel"`
	Token            string    `json:"token"`
	AgentRtcUID      string    `json:"agent_rtc_uid"`
	RemoteRtcUIDs    []string  `json:"remote_rtc_uids"`
	EnableStringUID  bool      `json:"enable_string_uid"`
	IdleTimeout      int       `json:"idle_timeout"`
	ASR              ASR       `json:"asr"`
	LLM              LLM       `json:"llm"`
	TTS              TTSConfig `json:"tts"`
	VAD              VAD       `json:"vad"`
	AdvancedFeatures Features  `json:"advanced_features"`
}

// ASR represents the Automatic Speech Recognition configuration
type ASR struct {
	Language string `json:"language"`
	Task     string `json:"task"`
}

// LLM represents the Language Learning Model configuration
type LLM struct {
	URL              string          `json:"url"`
	APIKey           string          `json:"api_key"`
	SystemMessages   []SystemMessage `json:"system_messages"`
	GreetingMessage  string          `json:"greeting_message"`
	FailureMessage   string          `json:"failure_message"`
	MaxHistory       int             `json:"max_history"`
	Params           LLMParams       `json:"params"`
	InputModalities  []string        `json:"input_modalities"`
	OutputModalities []string        `json:"output_modalities"`
}

// SystemMessage represents a system message in the conversation
type SystemMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMParams represents the parameters for the Language Learning Model
type LLMParams struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
}

// VAD represents the Voice Activity Detection configuration
type VAD struct {
	SilenceDurationMS   int     `json:"silence_duration_ms"`
	SpeechDurationMS    int     `json:"speech_duration_ms"`
	Threshold           float64 `json:"threshold"`
	InterruptDurationMS int     `json:"interrupt_duration_ms"`
	PrefixPaddingMS     int     `json:"prefix_padding_ms"`
}

// Features represents advanced features configuration
type Features struct {
	EnableAIVAD bool `json:"enable_aivad"`
	EnableBHVS  bool `json:"enable_bhvs"`
}

// InviteAgentResponse represents the response for an agent invitation
type InviteAgentResponse struct {
	AgentID  string `json:"agent_id"`
	CreateTS int64  `json:"create_ts"`
	Status   string `json:"status"`
}

// RemoveAgentResponse represents the response for an agent removal
type RemoveAgentResponse struct {
	Success bool   `json:"success"`
	AgentID string `json:"agent_id"`
}

// ConvoAIConfig holds all configuration for the ConvoAI service
type ConvoAIConfig struct {
	// Agora Configuration
	AppID          string
	AppCertificate string
	CustomerID     string
	CustomerSecret string
	BaseURL        string
	AgentUID       string

	// LLM Configuration
	LLMModel string
	LLMURL   string
	LLMToken string

	// TTS Configuration
	TTSVendor     string
	MicrosoftTTS  *MicrosoftTTSConfig
	ElevenLabsTTS *ElevenLabsTTSConfig

	// Modalities Configuration
	InputModalities  string
	OutputModalities string
}

// MicrosoftTTSConfig holds Microsoft TTS specific configuration
type MicrosoftTTSConfig struct {
	Key       string `json:"key"`
	Region    string `json:"region"`
	VoiceName string `json:"voice_name"`
	Rate      string `json:"rate"`
	Volume    string `json:"volume"`
}

// ElevenLabsTTSConfig holds ElevenLabs TTS specific configuration
type ElevenLabsTTSConfig struct {
	Key     string `json:"key"`
	VoiceID string `json:"voice_id"`
	ModelID string `json:"model_id"`
}
```

These new types give some insight into all the parts we'll be assembling in the next steps. We'll take the client request, and use it to configure the `AgoraStartRequest` and send it to Agora's Conversational AI Engine. Agora's Convo AI engine will add the agent to the conversation.

## ConvoAI Service

With our types defined, let's implement the agent routes for inviting and removing agents from conversations.

Create the `convoai-service.go` file:

```bash
touch convoai/convoai-service.go
```

Start with importing gin and the `agora-token` library, because we'll need to generate tokens for the agent. Then we'll register and set up the agent routes. These functions will validate the request before passing it to their respective handlers.

```go
package convoai

import (
	"net/http"

	"github.com/AgoraIO-Community/convo-ai-go-server/token_service"
	"github.com/gin-gonic/gin"
)

// ConvoAIService handles AI conversation functionality
type ConvoAIService struct {
	config       *ConvoAIConfig
	tokenService *token_service.TokenService
}

// NewConvoAIService creates a new ConvoAIService instance
func NewConvoAIService(config *ConvoAIConfig, tokenService *token_service.TokenService) *ConvoAIService {
	return &ConvoAIService{
		config:       config,
		tokenService: tokenService,
	}
}

// Register the ConvoAI service routes
func (s *ConvoAIService) RegisterRoutes(router *gin.Engine) {
	agent := router.Group("/agent")
	agent.POST("/invite", s.InviteAgent)
	agent.POST("/remove", s.RemoveAgent)
}

// InviteAgent handles the agent invitation request
func (s *ConvoAIService) InviteAgent(c *gin.Context) {
	var req InviteAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request
	if err := s.validateInviteRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the handler
	response, err := s.HandleInviteAgent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RemoveAgent handles the agent removal request
func (s *ConvoAIService) RemoveAgent(c *gin.Context) {
	var req RemoveAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request
	if err := s.validateRemoveRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the handler
	response, err := s.HandleRemoveAgent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
```

## Invite Agent Handler

Next, we'll implement the invite handler, which needs to handle several key tasks:

- Generate a token for the AI agent to access the RTC channel.
- Configure Text-to-Speech (Microsoft or ElevenLabs)
- Define the AI agent's prompt and greeting message.
- Configure the Voice Activity Detection (VAD), which controls conversation flow
- Sends the start request to Agora's Conversational AI Engine.
- Returns the response to the client that contains the AgentID from Agora's Convo AI Engine response.

Create the file `convoai_handler_invite.go`:

```bash
touch convoai/convoai_handler_invite.go
```

Add the following content:

```go
package convoai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"crypto/rand"

	"github.com/AgoraIO-Community/convo-ai-go-server/token_service"
)

// HandleInviteAgent processes the agent invitation request
func (s *ConvoAIService) HandleInviteAgent(req InviteAgentRequest) (*InviteAgentResponse, error) {
	// Generate token for the agent
	tokenReq := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   req.ChannelName,
		Uid:       "0",
		RtcRole:   "publisher",
	}

	token, err := s.tokenService.GenRtcToken(tokenReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	// Get TTS config based on vendor
	ttsConfig, err := s.getTTSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TTS config: %v", err)
	}

	// Set up system message for AI behavior
	systemMessage := SystemMessage{
		Role:    "system",
		Content: "You are a helpful assistant. Pretend that the text input is audio, and you are responding to it. Speak fast, clearly, and concisely.",
	}

	// Set default modalities if not provided
	inputModalities := req.InputModalities
	if len(inputModalities) == 0 {
		inputModalities = []string{"text"}
	}

	outputModalities := req.OutputModalities
	if len(outputModalities) == 0 {
		outputModalities = []string{"text", "audio"}
	}

	// Build the request body for Agora Conversation AI service
	agoraReq := AgoraStartRequest{
		Name: fmt.Sprintf("agent-%d-%s", time.Now().UnixNano(), randomString(6)),
		Properties: Properties{
			Channel:         req.ChannelName,
			Token:           token,
			AgentRtcUID:     s.config.AgentUID,
			RemoteRtcUIDs:   getRemoteRtcUIDs(req.RequesterID),
			EnableStringUID: isStringUID(req.RequesterID),
			IdleTimeout:     30,
			ASR: ASR{
				Language: "en-US",
				Task:     "conversation",
			},
			LLM: LLM{
				URL:             s.config.LLMURL,
				APIKey:          s.config.LLMToken,
				SystemMessages:  []SystemMessage{systemMessage},
				GreetingMessage: "Hello! How can I assist you today?",
				FailureMessage:  "Please wait a moment.",
				MaxHistory:      10,
				Params: LLMParams{
					Model:       s.config.LLMModel,
					MaxTokens:   1024,
					Temperature: 0.7,
					TopP:        0.95,
				},
				InputModalities:  inputModalities,
				OutputModalities: outputModalities,
			},
			TTS: *ttsConfig,
			VAD: VAD{
				SilenceDurationMS:   480,
				SpeechDurationMS:    15000,
				Threshold:           0.5,
				InterruptDurationMS: 160,
				PrefixPaddingMS:     300,
			},
			AdvancedFeatures: Features{
				EnableAIVAD: false,
				EnableBHVS:  false,
			},
		},
	}

	// Debug logging
	prettyJSON, _ := json.MarshalIndent(agoraReq, "", "  ")
	fmt.Printf("Sending request to start agent: %s\n", string(prettyJSON))

	// Convert request to JSON
	jsonData, err := json.Marshal(agoraReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/%s/join", s.config.BaseURL, s.config.AppID)
	fmt.Printf("URL: %s\n", url)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", s.getBasicAuth())

	// Send the request using a client with a timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v (URL: %s)", err, url)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to start conversation: status=%d, body=%s, url=%s, headers=%v",
			resp.StatusCode, string(body), url, httpReq.Header)
	}

	// Parse the response
	var agoraResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&agoraResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Create the response
	response := &InviteAgentResponse{
		AgentID:  agoraResp["agent_id"].(string),
		CreateTS: time.Now().Unix(),
		Status:   "RUNNING",
	}

	return response, nil
}

// getRemoteRtcUIDs returns the appropriate RemoteRtcUIDs array based on the requesterID
func getRemoteRtcUIDs(requesterID string) []string {
	return []string{requesterID}
}

// Add this helper function
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
```

## Remove Agent Handler

After the agent joins the conversation, we need a way to remove them from the conversation. This is where the remove handler comes in. It takes the `agentID` and sends a request to the Agora's Conversational AI Engine to remove the agent from the channel.

Create the file `convoai_handler_remove.go`:

```bash
touch convoai/convoai_handler_remove.go
```

Add the following:

```go
package convoai

import (
	"fmt"
	"net/http"
	"time"
)

// HandleRemoveAgent processes the agent removal request
func (s *ConvoAIService) HandleRemoveAgent(req RemoveAgentRequest) (*RemoveAgentResponse, error) {
	// Create the HTTP request
	url := fmt.Sprintf("%s/%s/agents/%s/leave", s.config.BaseURL, s.config.AppID, req.AgentID)
	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	auth := s.getBasicAuth()
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", auth)

	// Send the request using a client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to remove agent: %d", resp.StatusCode)
	}

	// Return success response
	response := &RemoveAgentResponse{
		Success: true,
		AgentID: req.AgentID,
	}

	return response, nil
}
```

## Utility Functions

In both the invite and remove routes, we need to use BasicAuthorization in the headers of our requests, so we'll set up a utility function to handle this.

Another utility we need to build is the `getTTSConfig`. I need to call out, because normally you would have a single TTS config. For demo purposes, I've built it this way to show how to implement the configs for all TTS vendors supported by Agora's Convo AI Engine.

Create the file `convoai-utils.go`:

```bash
touch convoai/convoai-utils.go
```

Add the following content:

```go
package convoai

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
)

func (s *ConvoAIService) getBasicAuth() string {
	auth := fmt.Sprintf("%s:%s", s.config.CustomerID, s.config.CustomerSecret)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// Helper function to check if the string is purely numeric (false) or contains any non-digit characters (true)
func isStringUID(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return true // Contains non-digit character
		}
	}
	return false // Contains only digits
}

// getTTSConfig returns the appropriate TTS configuration based on the configured vendor
func (s *ConvoAIService) getTTSConfig() (*TTSConfig, error) {
	switch s.config.TTSVendor {
	case string(TTSVendorMicrosoft):
		if s.config.MicrosoftTTS == nil ||
			s.config.MicrosoftTTS.Key == "" ||
			s.config.MicrosoftTTS.Region == "" ||
			s.config.MicrosoftTTS.VoiceName == "" ||
			s.config.MicrosoftTTS.Rate == "" ||
			s.config.MicrosoftTTS.Volume == "" {
			return nil, fmt.Errorf("missing Microsoft TTS configuration")
		}

		// Convert rate and volume from string to float64
		rate, err := strconv.ParseFloat(s.config.MicrosoftTTS.Rate, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid rate value: %v", err)
		}

		volume, err := strconv.ParseFloat(s.config.MicrosoftTTS.Volume, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid volume value: %v", err)
		}

		return &TTSConfig{
			Vendor: TTSVendorMicrosoft,
			Params: map[string]interface{}{
				"key":        s.config.MicrosoftTTS.Key,
				"region":     s.config.MicrosoftTTS.Region,
				"voice_name": s.config.MicrosoftTTS.VoiceName,
				"rate":       rate,
				"volume":     volume,
			},
		}, nil

	case string(TTSVendorElevenLabs):
		if s.config.ElevenLabsTTS == nil ||
			s.config.ElevenLabsTTS.Key == "" ||
			s.config.ElevenLabsTTS.ModelID == "" ||
			s.config.ElevenLabsTTS.VoiceID == "" {
			return nil, fmt.Errorf("missing ElevenLabs TTS configuration")
		}
		return &TTSConfig{
			Vendor: TTSVendorElevenLabs,
			Params: map[string]interface{}{
				"api_key":  s.config.ElevenLabsTTS.Key,
				"model_id": s.config.ElevenLabsTTS.ModelID,
				"voice_id": s.config.ElevenLabsTTS.VoiceID,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported TTS vendor: %s", s.config.TTSVendor)
	}
}

// validateInviteRequest validates the invite agent request
func (s *ConvoAIService) validateInviteRequest(req *InviteAgentRequest) error {
	if req.RequesterID == "" {
		return errors.New("requester_id is required")
	}

	if req.ChannelName == "" {
		return errors.New("channel_name is required")
	}

	// Validate channel_name length
	if len(req.ChannelName) < 3 || len(req.ChannelName) > 64 {
		return errors.New("channel_name length must be between 3 and 64 characters")
	}

	return nil
}

// validateRemoveRequest validates the remove agent request
func (s *ConvoAIService) validateRemoveRequest(req *RemoveAgentRequest) error {
	if req.AgentID == "" {
		return errors.New("agent_id is required")
	}
	return nil
}
```

## HTTP Headers

To handle all header-related logic, create the `httpHeaders.go` file:

```bash
touch http_headers/httpHeaders.go
```

Add the following content:

```go
package http_headers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HttpHeaders holds configurations for handling requests, such as CORS settings.
type HttpHeaders struct {
	AllowOrigin string // List of origins allowed to access the resources.
}

// NewHttpHeaders initializes and returns a new Middleware object with specified CORS settings.
func NewHttpHeaders(allowOrigin string) *HttpHeaders {
	return &HttpHeaders{AllowOrigin: allowOrigin}
}

// NoCache sets HTTP headers to prevent client-side caching of responses.
func (m *HttpHeaders) NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set multiple cache-related headers to ensure responses are not cached.
		c.Header("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Header("Expires", "-1")
		c.Header("Pragma", "no-cache")
	}
}

// CORShttpHeaders adds CORS (Cross-Origin Resource Sharing) headers to responses and handles pre-flight requests.
// It allows web applications at different domains to interact more securely.
func (m *HttpHeaders) CORShttpHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Check if the origin of the request is allowed to access the resource.
		if !m.isOriginAllowed(origin) {
			// If not allowed, return a JSON error and abort the request.
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Origin not allowed",
			})
			c.Abort()
			return
		}
		// Set CORS headers to allow requests from the specified origin.
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		// Handle pre-flight OPTIONS requests.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// isOriginAllowed checks whether the provided origin is in the list of allowed origins.
func (m *HttpHeaders) isOriginAllowed(origin string) bool {
	if m.AllowOrigin == "*" {
		// Allow any origin if the configured setting is "*".
		return true
	}

	allowedOrigins := strings.Split(m.AllowOrigin, ",")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

// Timestamp adds a timestamp header to responses.
// This can be useful for debugging and logging purposes to track when a response was generated.
func (m *HttpHeaders) Timestamp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Proceed to the next middleware/handler.

		// Add the current timestamp to the response header after handling the request.
		timestamp := time.Now().Format(time.RFC3339)
		c.Writer.Header().Set("X-Timestamp", timestamp)
	}
}
```

### Update Main Server

Let's update our main `main.go` file to add our headers and register the `convoai-service`.

Open the `cmd/main.go` and add:

```typescript
import(
    // Previous imports remain the same
	"github.com/AgoraIO-Community/convo-ai-go-server/convoai"
	"github.com/AgoraIO-Community/convo-ai-go-server/http_headers"
);

// Previous code remains the same..
func setupServer() *http.Server {

    // Previous code remains the same..
    // Set up router with headers
	router := gin.Default()
    // Replace headers TODO:
    var httpHeaders = http_headers.NewHttpHeaders(corsAllowOrigin)
	router.Use(httpHeaders.NoCache())
	router.Use(httpHeaders.CORShttpHeaders())
	router.Use(httpHeaders.Timestamp())


	// Initialize services & register routes
	tokenService := token_service.NewTokenService(config.AppID, config.AppCertificate)
	tokenService.RegisterRoutes(router)

	convoAIService := convoai.NewConvoAIService(config, tokenService)
	convoAIService.RegisterRoutes(router)

// Rest of the code remains the same...
```

By now you've noticed that we added a token service that doesn't exist, ignore the error for now, becuse in the next step we'll implement the token service, which will make it easier to test and integrate with frontend applications.

## Token Generation

In the `convoai-service` we use a token service. While you could tie this to your auth service and have it generate the tokens. For this guide, we'll implement a token service for both the `convoai-service` and our client apps if needed.

Explaining this code is a bit outside the scope of this guide, but if you are new to tokens, I would recommend checking out my guide [Building a Token Server for Agora Applications using Golang](https://www.agora.io/en/blog/how-to-build-a-token-server-using-golang/).

## Token Service

Create the token service and handler files:

```bash
touch token_service/token-service.go
touch token_service/token_handlers.go
```

First, add the token service definition in `token-service.go`:

```go
package token_service

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// TokenService represents the main application token service.
type TokenService struct {
	Server         *http.Server   // The HTTP server for the application
	Sigint         chan os.Signal // Channel to handle OS signals, such as Ctrl+C
	appID          string         // The Agora app ID
	appCertificate string         // The Agora app certificate
}

// TokenRequest is a struct representing the JSON payload structure for token generation requests.
type TokenRequest struct {
	TokenType         string `json:"tokenType"`         // The token type: "rtc", "rtm", or "chat"
	Channel           string `json:"channel,omitempty"` // The channel name (used for RTC and RTM tokens)
	RtcRole           string `json:"role,omitempty"`    // The role of the user for RTC tokens (publisher or subscriber)
	Uid               string `json:"uid,omitempty"`     // The user ID or account (used for RTC, RTM, and some chat tokens)
	ExpirationSeconds int    `json:"expire,omitempty"`  // The token expiration time in seconds (used for all token types)
}

// NewTokenService initializes and returns a TokenService pointer with all configurations set.
func NewTokenService(appIDEnv string, appCertEnv string) *TokenService {
	return &TokenService{
		appID:          appIDEnv,
		appCertificate: appCertEnv,
	}
}

// RegisterRoutes registers the routes for the TokenService.
func (s *TokenService) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/token")
	api.POST("/getNew", s.GetToken)
}

// GetToken handles the HTTP request to generate a token based on the provided TokenRequest.
func (s *TokenService) GetToken(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var tokenReq TokenRequest
	// Parse the request body into a TokenRequest struct
	err := json.NewDecoder(req.Body).Decode(&tokenReq)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleGetToken(tokenReq, respWriter)
}
```

Next, add the token handlers in `token_handlers.go`:

```go
package token_service

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AgoraIO-Community/go-tokenbuilder/chatTokenBuilder"
	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
)

// HandleGetToken handles the HTTP request to generate a token based on the provided tokenType.
func (s *TokenService) HandleGetToken(tokenReq TokenRequest, w http.ResponseWriter) {
	var token string
	var tokenErr error

	switch tokenReq.TokenType {
	case "rtc":
		token, tokenErr = s.GenRtcToken(tokenReq)
	case "rtm":
		token, tokenErr = s.GenRtmToken(tokenReq)
	case "chat":
		token, tokenErr = s.GenChatToken(tokenReq)
	default:
		http.Error(w, "Unsupported tokenType", http.StatusBadRequest)
		return
	}
	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{Token: token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GenRtcToken generates an RTC token based on the provided TokenRequest and returns it.
func (s *TokenService) GenRtcToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.Channel == "" {
		return "", errors.New("invalid: missing channel name")
	}
	if tokenRequest.Uid == "" {
		return "", errors.New("invalid: missing user ID or account")
	}

	var userRole rtctokenbuilder2.Role
	if tokenRequest.RtcRole == "publisher" {
		userRole = rtctokenbuilder2.RolePublisher
	} else {
		userRole = rtctokenbuilder2.RoleSubscriber
	}

	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	uid64, parseErr := strconv.ParseUint(tokenRequest.Uid, 10, 64)
	if parseErr != nil {
		return rtctokenbuilder2.BuildTokenWithAccount(
			s.appID, s.appCertificate, tokenRequest.Channel,
			tokenRequest.Uid, userRole, uint32(tokenRequest.ExpirationSeconds),
		)
	}

	return rtctokenbuilder2.BuildTokenWithUid(
		s.appID, s.appCertificate, tokenRequest.Channel,
		uint32(uid64), userRole, uint32(tokenRequest.ExpirationSeconds),
	)
}

// GenRtmToken generates an RTM (Real-Time Messaging) token based on the provided TokenRequest and returns it.
func (s *TokenService) GenRtmToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.Uid == "" {
		return "", errors.New("invalid: missing user ID or account")
	}
	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	return rtmtokenbuilder2.BuildToken(
		s.appID, s.appCertificate,
		tokenRequest.Uid,
		uint32(tokenRequest.ExpirationSeconds),
		tokenRequest.Channel,
	)
}

// GenChatToken generates a chat token based on the provided TokenRequest and returns it.
func (s *TokenService) GenChatToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	var chatToken string
	var tokenErr error

	if tokenRequest.Uid == "" {
		chatToken, tokenErr = chatTokenBuilder.BuildChatAppToken(
			s.appID, s.appCertificate, uint32(tokenRequest.ExpirationSeconds),
		)
	} else {
		chatToken, tokenErr = chatTokenBuilder.BuildChatUserToken(
			s.appID, s.appCertificate,
			tokenRequest.Uid,
			uint32(tokenRequest.ExpirationSeconds),
		)
	}

	return chatToken, tokenErr
}
```

With the token generation in place, let's add some validation middleware to ensure our API is robust and secure.

## Environment Validation

Create a validation utility to check that all required environment variables are set. Create the file `validation/validation.go`:

```bash
touch validation/validation.go
```

Add the following content:

```go
package validation

import (
	"errors"
	"strings"

	"github.com/AgoraIO-Community/convo-ai-go-server/convoai"
)

// ValidateEnvironment checks if all required environment variables are set
func ValidateEnvironment(config *convoai.ConvoAIConfig) error {
	// Validate Agora Configuration
	if config.AppID == "" || config.AppCertificate == "" {
		return errors.New("config error: Agora credentials (APP_ID, APP_CERTIFICATE) are not set")
	}

	if config.CustomerID == "" || config.CustomerSecret == "" || config.BaseURL == "" {
		return errors.New("config error: Agora Conversation AI credentials (CUSTOMER_ID, CUSTOMER_SECRET, BASE_URL) are not set")
	}

	// Validate LLM Configuration
	if config.LLMURL == "" || config.LLMToken == "" {
		return errors.New("config error: LLM configuration (LLM_URL, LLM_TOKEN) is not set")
	}

	// Validate TTS Configuration
	if config.TTSVendor == "" {
		return errors.New("config error: TTS_VENDOR is not set")
	}

	if err := validateTTSConfig(config); err != nil {
		return err
	}

	// Validate Modalities (optional, using defaults if not set)
	if config.InputModalities != "" && !validateModalities(config.InputModalities) {
		return errors.New("config error: Invalid INPUT_MODALITIES format")
	}
	if config.OutputModalities != "" && !validateModalities(config.OutputModalities) {
		return errors.New("config error: Invalid OUTPUT_MODALITIES format")
	}

	return nil
}

// Validates the TTS configuration based on the vendor
func validateTTSConfig(config *convoai.ConvoAIConfig) error {
	switch config.TTSVendor {
	case "microsoft":
		if config.MicrosoftTTS == nil {
			return errors.New("config error: Microsoft TTS configuration is missing")
		}
		if config.MicrosoftTTS.Key == "" ||
			config.MicrosoftTTS.Region == "" ||
			config.MicrosoftTTS.VoiceName == "" {
			return errors.New("config error: Microsoft TTS configuration is incomplete")
		}
	case "elevenlabs":
		if config.ElevenLabsTTS == nil {
			return errors.New("config error: ElevenLabs TTS configuration is missing")
		}
		if config.ElevenLabsTTS.Key == "" ||
			config.ElevenLabsTTS.VoiceID == "" ||
			config.ElevenLabsTTS.ModelID == "" {
			return errors.New("config error: ElevenLabs TTS configuration is incomplete")
		}
	default:
		return errors.New("config error: Unsupported TTS vendor: " + config.TTSVendor)
	}
	return nil
}

// Checks if the modalities string is properly formatted
func validateModalities(modalities string) bool {
	// map of valid modalities
	validModalities := map[string]bool{
		"text":  true,
		"audio": true,
	}
	// split the modalities string and check if each modality is valid
	for _, modality := range strings.Split(modalities, ",") {
		if !validModalities[strings.TrimSpace(modality)] {
			return false
		}
	}
	return true
}
```

This validation utility ensures that all required environment variables are properly set before the server starts.

Open the `main.go` and update the `setupServer` function to use the validation utility:

```go
	// Load configuration

    // Replace the TODO: comment with the following:
    // Validate environment configuration
    if err := validation.ValidateEnvironment(config); err != nil {
        log.Fatal("FATAL ERROR: ", err)
    }

// Rest of the code remains the same...
```

## Running the Server

Now that we have all the components in place, let's run the server. First, make sure you have set up the `.env` file with all the necessary credentials. The server will automatically load these environment variables at startup.

Build and run the server:

```bash
go build -o server
./server
```

If you've set up everything correctly, you should see the server starting up and listening on the configured port (default is 8080).

## Testing the Server

Before testing the endpoints, make sure you have a client-side app running. You can use any application that implements Agora's video SDK (web, mobile, or desktop). If you don't have an app, you can use [Agora's Voice Demo](https://webdemo.agora.io/basicVoiceCall/index.html), just make sure to make a token request before joining the channel.

Let's test our API endpoints using curl:

### 1. Generate a Token

```bash
curl -X POST http://localhost:8080/token/getNew \
  -H "Content-Type: application/json" \
  -d '{
    "tokenType": "rtc",
    "channel": "test-channel",
    "uid": "1234",
    "role": "publisher"
  }'
```

Expected response:

```json
{
  "token": "007eJxTYBAxNdgrlvnEfm3o..."
}
```

### 2. Invite an AI Agent

```bash
curl -X POST http://localhost:8080/agent/invite \
  -H "Content-Type: application/json" \
  -d '{
    "requester_id": "1234",
    "channel_name": "test-channel",
    "input_modalities": ["text"],
    "output_modalities": ["text", "audio"]
  }'
```

Expected response:

```json
{
  "agent_id": "agent-123abc",
  "create_ts": 1665481725000,
  "status": "RUNNING"
}
```

### 3. Remove an AI Agent

```bash
curl -X POST http://localhost:8080/agent/remove \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent-123abc"
  }'
```

Expected response:

```json
{
  "success": true,
  "agent_id": "agent-123abc"
}
```

## Customizations

Agora Conversational AI Engine supports several customizations.

### Customizing the Agent

In the `convoai_handler_invite.go` file, you can modify the system message to customize the agent's behavior:

```go
systemMessage := SystemMessage{
    Role:    "system",
    Content: "You are a technical support specialist named Alex. Your responses should be friendly but concise, focused on helping users solve their technical problems. Use simple language but don't oversimplify technical concepts.",
}
```

You can also update the greeting message to control the initial message the agent speaks when joining the channel:

```go
LLM: LLM{
    // ... other configurations
    GreetingMessage: "Hello! I'm Alex, your technical support specialist. How can I assist you today?",
    FailureMessage:  "I'm processing your request. Please give me a moment.",
    // ... rest of the configuration
}
```

### Customizing Speech Synthesis

Choose the right voice for your application by exploring the voice libraries:

- For Microsoft Azure TTS: Visit the [Microsoft Azure TTS Voice Gallery](https://speech.microsoft.com/portal/voicegallery)
- For ElevenLabs TTS: Explore the [ElevenLabs Voice Library](https://elevenlabs.io/voice-library)

Update the `.env` file with the appropriate voice settings.

### Fine-tuning Voice Activity Detection

Adjust VAD settings in `convoai_handler_invite.go` to optimize conversation flow:

```go
VAD: VAD{
    SilenceDurationMS:   600,      // How long to wait after silence to end turn
    SpeechDurationMS:    10000,     // Maximum duration for a single speech segment
    Threshold:           0.6,       // Speech detection sensitivity
    InterruptDurationMS: 200,       // How quickly interruptions are detected
    PrefixPaddingMS:     400,       // Audio padding at the beginning of speech
},
```

## Complete Environment Variables Reference

Here's a complete list of environment variables for your `.env` file:

```
# Server Configuration
PORT=8080
CORS_ALLOW_ORIGIN=*

# Agora Configuration
AGORA_APP_ID=your_app_id
AGORA_APP_CERTIFICATE=your_app_certificate
AGORA_CONVO_AI_BASE_URL=https://api.agora.io/api/conversational-ai-agent/v2/projects
AGORA_CUSTOMER_ID=your_customer_id
AGORA_CUSTOMER_SECRET=your_customer_secret
AGENT_UID=Agent

# LLM Configuration
LLM_URL=https://api.openai.com/v1/chat/completions
LLM_TOKEN=your_openai_api_key
LLM_MODEL=gpt-4o-mini

# Input/Output Modalities
INPUT_MODALITIES=text
OUTPUT_MODALITIES=text,audio

# TTS Configuration
TTS_VENDOR=microsoft  # or elevenlabs

# Microsoft TTS Configuration
MICROSOFT_TTS_KEY=your_microsoft_tts_key
MICROSOFT_TTS_REGION=your_microsoft_tts_region
MICROSOFT_TTS_VOICE_NAME=en-US-GuyNeural
MICROSOFT_TTS_RATE=1.0
MICROSOFT_TTS_VOLUME=100.0

# ElevenLabs TTS Configuration
ELEVENLABS_API_KEY=your_elevenlabs_api_key
ELEVENLABS_VOICE_ID=your_elevenlabs_voice_id
ELEVENLABS_MODEL_ID=eleven_monolingual_v1
```

## Next Steps

Congratulations! You've built a Go server that integrates with Agora's Conversational AI Engine. Take this microservice and integrate it with your existing Agora backends.

For more information about [Agora's Conversational AI Engine](https://www.agora.io/en/products/conversational-ai-engine/) check out the [official documentation](https://docs.agora.io/en/).

Happy building!
