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
	// Empty string is valid (will use defaults)
	if modalities == "" {
		return true
	}
	
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
