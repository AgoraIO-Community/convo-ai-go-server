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

	// TODO: Remove debug logging
	fmt.Printf("Request headers: %v\n", httpReq.Header)

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
	// if requesterID == "0" {
	// 	return []string{"*"}
	// }
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
