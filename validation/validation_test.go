package validation

import (
	"testing"

	"github.com/AgoraIO-Community/convo-ai-go-server/convoai"
)

func TestValidateEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		config  *convoai.ConvoAIConfig
		wantErr bool
	}{
		{
			name: "Valid config with Microsoft TTS",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				AgentUID:       "123456",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
				TTSVendor:      "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:       "ms-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
					Rate:      "1.0",
					Volume:    "1.0",
				},
				InputModalities:  "text,audio",
				OutputModalities: "text,audio",
			},
			wantErr: false,
		},
		{
			name: "Valid config with ElevenLabs TTS",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				AgentUID:       "123456",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
				TTSVendor:      "elevenlabs",
				ElevenLabsTTS: &convoai.ElevenLabsTTSConfig{
					Key:     "el-key",
					VoiceID: "voice-id",
					ModelID: "model-id",
				},
				InputModalities:  "text",
				OutputModalities: "text,audio",
			},
			wantErr: false,
		},
		{
			name: "Missing Agora credentials",
			config: &convoai.ConvoAIConfig{
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
				TTSVendor:      "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:       "ms-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing Agora Conversation AI credentials",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
				TTSVendor:      "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:       "ms-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing LLM configuration",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				LLMModel:       "gpt-3.5-turbo",
				TTSVendor:      "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:       "ms-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing TTS vendor",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
			},
			wantErr: true,
		},
		{
			name: "Unsupported TTS vendor",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
				TTSVendor:      "unsupported",
			},
			wantErr: true,
		},
		{
			name: "Invalid input modalities",
			config: &convoai.ConvoAIConfig{
				AppID:          "app-id",
				AppCertificate: "app-cert",
				CustomerID:     "customer-id",
				CustomerSecret: "customer-secret",
				BaseURL:        "https://api.example.com",
				LLMModel:       "gpt-3.5-turbo",
				LLMURL:         "https://api.openai.com/v1/chat/completions",
				LLMToken:       "llm-token",
				TTSVendor:      "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:       "ms-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
				},
				InputModalities: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnvironment(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTTSConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *convoai.ConvoAIConfig
		wantErr bool
	}{
		{
			name: "Valid Microsoft TTS",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:       "ms-key",
					Region:    "eastus",
					VoiceName: "en-US-AriaNeural",
					Rate:      "1.0",
					Volume:    "1.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid ElevenLabs TTS",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "elevenlabs",
				ElevenLabsTTS: &convoai.ElevenLabsTTSConfig{
					Key:     "el-key",
					VoiceID: "voice-id",
					ModelID: "model-id",
				},
			},
			wantErr: false,
		},
		{
			name: "Microsoft TTS missing config",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "microsoft",
			},
			wantErr: true,
		},
		{
			name: "Microsoft TTS incomplete config",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "microsoft",
				MicrosoftTTS: &convoai.MicrosoftTTSConfig{
					Key:    "ms-key",
					Region: "eastus",
				},
			},
			wantErr: true,
		},
		{
			name: "ElevenLabs TTS missing config",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "elevenlabs",
			},
			wantErr: true,
		},
		{
			name: "ElevenLabs TTS incomplete config",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "elevenlabs",
				ElevenLabsTTS: &convoai.ElevenLabsTTSConfig{
					Key:     "el-key",
					VoiceID: "voice-id",
				},
			},
			wantErr: true,
		},
		{
			name: "Unsupported TTS vendor",
			config: &convoai.ConvoAIConfig{
				TTSVendor: "unsupported",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTTSConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTTSConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateModalities(t *testing.T) {
	tests := []struct {
		name       string
		modalities string
		want       bool
	}{
		{
			name:       "Valid single modality",
			modalities: "text",
			want:       true,
		},
		{
			name:       "Valid multiple modalities",
			modalities: "text,audio",
			want:       true,
		},
		{
			name:       "Valid modalities with spaces",
			modalities: "text, audio",
			want:       true,
		},
		{
			name:       "Invalid modality",
			modalities: "invalid",
			want:       false,
		},
		{
			name:       "Mixed valid and invalid modalities",
			modalities: "text,invalid",
			want:       false,
		},
		{
			name:       "Empty modalities",
			modalities: "",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateModalities(tt.modalities); got != tt.want {
				t.Errorf("validateModalities() = %v, want %v", got, tt.want)
			}
		})
	}
}