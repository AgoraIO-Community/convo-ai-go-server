# ConvoAI Service Flow Diagram

```mermaid
flowchart LR
    subgraph Client[Client Layer]
        A[HTTP Client]
    end

    subgraph Server[API Layer]
        B[Router]
        C[ConvoAI Service]
        D[Token Service]
    end

    subgraph Logic[Agent Logic]
        E[Invite Handler]
        F[Remove Handler]
    end

    subgraph External[External Services]
        G[Agora ConvoAI API]
    end

    %% Main Flow
    A <--> B
    B <--> C
    C <--> E

    %% Service Interactions
    F <--> C
    E <--> D
    E & F <--> G

    %% Styling
    classDef default fill:#f9f,stroke:#333,stroke-width:2px,color:#000;
    classDef external fill:#bbf,stroke:#333,stroke-width:2px,color:#000;
    classDef config fill:#bfb,stroke:#333,stroke-width:2px,color:#000;

    class G external
    class A config
```

## Service Flow Description

1. **Client Request**

   - Client sends HTTP request to invite or remove an AI agent

2. **Router Processing**

   - Gin router directs request to appropriate ConvoAI Service handler

3. **Invite Flow**

   - Validates request
   - Generates RTC token
   - Configures TTS and LLM settings
   - Calls Agora ConvoAI API to start agent
   - Returns agent details to client

4. **Remove Flow**

   - Validates request
   - Calls Agora ConvoAI API to remove agent
   - Returns success status to client

5. **Configuration**
   - Uses environment variables for service configuration
   - Supports multiple TTS vendors
   - Configurable LLM parameters
