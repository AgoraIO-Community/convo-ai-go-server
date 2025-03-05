package convoai

import (
	"net/http"

	"github.com/AgoraIO-Community/convo-ai-go-server/token_service"
	"github.com/gin-gonic/gin"
)

// TokenGenerator is an interface for token generation services
type TokenGenerator interface {
	GenRtcToken(req token_service.TokenRequest) (string, error)
}

// ConvoAIService handles AI conversation functionality
type ConvoAIService struct {
	config       *ConvoAIConfig
	tokenService TokenGenerator
}

// NewConvoAIService creates a new ConvoAIService instance
func NewConvoAIService(config *ConvoAIConfig, tokenService TokenGenerator) *ConvoAIService {
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