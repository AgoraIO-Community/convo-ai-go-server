# ConvoAI Service Entity Relationship Diagram

```mermaid
erDiagram
    ConvoAIService ||--o{ InviteAgentRequest : handles
    ConvoAIService ||--o{ RemoveAgentRequest : handles
    ConvoAIService ||--|| ConvoAIConfig : uses
    ConvoAIConfig ||--o{ TTSConfig : contains
    ConvoAIConfig ||--|| LLMConfig : contains

    ConvoAIService {
        TokenService tokenService
        ConvoAIConfig config
    }

    ConvoAIConfig {
        string AppID
        string AppCertificate
        string CustomerID
        string CustomerSecret
        string BaseURL
        string AgentUID
        string LLMModel
        string LLMURL
        string LLMToken
        string TTSVendor
        string InputModalities
        string OutputModalities
    }

    InviteAgentRequest {
        string ChannelName
        string RequesterID
        string[] InputModalities
        string[] OutputModalities
    }

    RemoveAgentRequest {
        string ChannelName
        string RequesterID
        string AgentID
    }

    TTSConfig {
        string Vendor
        string Key
        string Region
        string VoiceName
        string Rate
        string Volume
    }

    LLMConfig {
        string URL
        string Token
        string Model
        int MaxTokens
        float Temperature
        float TopP
        int MaxHistory
    }

    InviteAgentResponse {
        string AgentID
        int CreateTS
        string Status
    }

    RemoveAgentResponse {
        boolean Success
        string AgentID
    }

    ConvoAIService ||--o{ InviteAgentResponse : generates
    ConvoAIService ||--o{ RemoveAgentResponse : generates
```

## Entity Descriptions

### ConvoAIService

- Main service that handles AI agent operations
- Manages configuration and token service integration

### ConvoAIConfig

- Holds all configuration parameters for the service
- Includes Agora credentials, LLM settings, and TTS configuration

### Request Entities

- **InviteAgentRequest**: Parameters for inviting an AI agent
- **RemoveAgentRequest**: Parameters for removing an AI agent

### Response Entities

- **InviteAgentResponse**: Contains agent details after successful invitation
- **RemoveAgentResponse**: Contains removal operation status

### Configuration Entities

- **TTSConfig**: Text-to-Speech configuration (supports multiple vendors)
- **LLMConfig**: Language Model configuration and parameters

## Relationships

- ConvoAIService handles both invite and remove requests
- Each request generates a corresponding response
- Service uses a single configuration instance
- Configuration can contain multiple TTS configurations
- All operations require proper authentication and configuration
