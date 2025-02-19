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
