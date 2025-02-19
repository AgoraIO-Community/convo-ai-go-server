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
	// if s == "0" {
	// 	return true // 0 sets the requester as a wildcard
	// }
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
