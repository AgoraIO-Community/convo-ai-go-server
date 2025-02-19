package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AgoraIO-Community/convo-ai-go-server/convoai"
	"github.com/AgoraIO-Community/convo-ai-go-server/http_headers"
	"github.com/AgoraIO-Community/convo-ai-go-server/token_service"
	"github.com/AgoraIO-Community/convo-ai-go-server/validation"
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

	// Validate environment configuration
	if err := validation.ValidateEnvironment(config); err != nil {
		log.Fatal("FATAL ERROR: ", err)
	}

	// Server Configuration
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// CORS Configuration
	corsAllowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")

	// Set up router with headers
	router := gin.Default()
	var httpHeaders = http_headers.NewHttpHeaders(corsAllowOrigin)
	router.Use(httpHeaders.NoCache())
	router.Use(httpHeaders.CORShttpHeaders())
	router.Use(httpHeaders.Timestamp())

	// Initialize services & register routes
	tokenService := token_service.NewTokenService(config.AppID, config.AppCertificate)
	tokenService.RegisterRoutes(router)

	convoAIService := convoai.NewConvoAIService(config, tokenService)
	convoAIService.RegisterRoutes(router)

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
