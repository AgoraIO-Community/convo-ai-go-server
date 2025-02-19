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
	Key       string
	Region    string
	VoiceName string
	Rate      string
	Volume    string
}

// ElevenLabsTTSConfig holds ElevenLabs TTS specific configuration
type ElevenLabsTTSConfig struct {
	APIKey  string
	VoiceID string
	ModelID string
}
