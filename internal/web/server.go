package web

import (
	"net/http"

	"github.com/RyanSStephens/TF-NLP-Agent/internal/ai"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/nlp"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/security"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/terraform"
	"github.com/gin-gonic/gin"
)

// Server represents the web server
type Server struct {
	router      *gin.Engine
	aiProvider  ai.Provider
	nlpEngine   *nlp.Engine
	tfGenerator *terraform.Generator
	secScanner  *security.Scanner
}

// GenerateRequest represents a generation request
type GenerateRequest struct {
	Description string `json:"description" binding:"required"`
	Provider    string `json:"provider,omitempty"`
}

// GenerateResponse represents a generation response
type GenerateResponse struct {
	Configuration string             `json:"configuration"`
	Issues        []security.Issue   `json:"issues,omitempty"`
	Costs         map[string]float64 `json:"estimated_costs,omitempty"`
	Success       bool               `json:"success"`
	Error         string             `json:"error,omitempty"`
}

// NewServer creates a new web server
func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)

	server := &Server{
		router:      gin.Default(),
		aiProvider:  ai.NewProvider("openai"),
		nlpEngine:   nlp.NewEngine(),
		tfGenerator: terraform.NewGenerator(),
		secScanner:  security.NewScanner(),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// API routes
	api := s.router.Group("/api/v1")
	{
		api.POST("/generate", s.handleGenerate)
		api.POST("/validate", s.handleValidate)
		api.GET("/health", s.handleHealth)
	}

	// Static files and web interface
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("web/templates/*")
	s.router.GET("/", s.handleIndex)
}

// handleGenerate handles Terraform configuration generation
func (s *Server) handleGenerate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenerateResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Parse natural language input
	parsed, err := s.nlpEngine.Parse(req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenerateResponse{
			Success: false,
			Error:   "Failed to parse description: " + err.Error(),
		})
		return
	}

	// Override cloud provider if specified
	if req.Provider != "" {
		parsed.CloudProvider = req.Provider
	}

	// Generate configuration using AI
	config, err := s.aiProvider.GenerateConfig(parsed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenerateResponse{
			Success: false,
			Error:   "Failed to generate configuration: " + err.Error(),
		})
		return
	}

	// Validate and format
	validated, err := s.tfGenerator.Validate(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenerateResponse{
			Success: false,
			Error:   "Failed to validate configuration: " + err.Error(),
		})
		return
	}

	// Security scan
	issues, err := s.secScanner.Scan(validated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenerateResponse{
			Success: false,
			Error:   "Security scan failed: " + err.Error(),
		})
		return
	}

	// Cost estimation
	costs, err := s.tfGenerator.EstimateCost(validated)
	if err != nil {
		// Don't fail on cost estimation errors
		costs = make(map[string]float64)
	}

	c.JSON(http.StatusOK, GenerateResponse{
		Configuration: validated,
		Issues:        issues,
		Costs:         costs,
		Success:       true,
	})
}

// handleValidate handles configuration validation
func (s *Server) handleValidate(c *gin.Context) {
	var req struct {
		Configuration string `json:"configuration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Validate configuration
	_, err := s.tfGenerator.Validate(req.Configuration)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Security scan
	issues, err := s.secScanner.Scan(req.Configuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Security scan failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"issues":  issues,
	})
}

// handleHealth handles health checks
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
	})
}

// handleIndex serves the main web interface
func (s *Server) handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "TF-NLP-Agent",
	})
}

// Start starts the web server
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
