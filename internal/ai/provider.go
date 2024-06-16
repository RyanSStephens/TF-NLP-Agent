package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/RyanSStephens/TF-NLP-Agent/internal/nlp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Provider represents an AI provider interface
type Provider interface {
	GenerateConfig(parsed *nlp.ParsedInput) (string, error)
	GenerateTerraform(input *nlp.ParsedInput) (string, error)
}

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	client *openai.Client
	model  string
}

// NewProvider creates a new AI provider based on the provider type
func NewProvider(providerType string) Provider {
	switch strings.ToLower(providerType) {
	case "openai":
		return NewOpenAIProvider()
	default:
		return NewOpenAIProvider() // Default to OpenAI
	}
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider() *OpenAIProvider {
	client := openai.NewClient(
		option.WithAPIKey(""), // Will be set from environment or config
	)

	return &OpenAIProvider{
		client: client,
		model:  "gpt-4",
	}
}

// GenerateConfig generates Terraform configuration using OpenAI
func (p *OpenAIProvider) GenerateConfig(parsed *nlp.ParsedInput) (string, error) {
	prompt := buildPrompt(parsed)

	resp, err := p.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a Terraform expert. Generate clean, secure, and production-ready Terraform configurations based on user requirements. Always include proper resource naming, tags, and security best practices."),
			openai.UserMessage(prompt),
		}),
		Model: openai.F(p.model),
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// Extract Terraform code from response (remove markdown formatting if present)
	terraformCode := extractTerraformCode(content)

	return terraformCode, nil
}

// buildPrompt constructs the prompt for the AI model
func buildPrompt(parsed *nlp.ParsedInput) string {
	var prompt strings.Builder

	prompt.WriteString("Generate a Terraform configuration based on the following requirements:\n\n")
	prompt.WriteString(fmt.Sprintf("Description: %s\n", parsed.OriginalText))

	if parsed.CloudProvider != "" {
		prompt.WriteString(fmt.Sprintf("Cloud Provider: %s\n", parsed.CloudProvider))
	}

	if len(parsed.Resources) > 0 {
		prompt.WriteString("Resources identified:\n")
		for _, resource := range parsed.Resources {
			prompt.WriteString(fmt.Sprintf("- %s: %s\n", resource.Type, resource.Name))
		}
	}

	if len(parsed.Requirements) > 0 {
		prompt.WriteString("Requirements:\n")
		for _, req := range parsed.Requirements {
			prompt.WriteString(fmt.Sprintf("- %s\n", req))
		}
	}

	prompt.WriteString("\nPlease provide a complete, working Terraform configuration that:\n")
	prompt.WriteString("1. Follows Terraform best practices\n")
	prompt.WriteString("2. Includes proper resource naming and tagging\n")
	prompt.WriteString("3. Implements security best practices\n")
	prompt.WriteString("4. Is production-ready\n")
	prompt.WriteString("5. Includes necessary variables and outputs\n")
	prompt.WriteString("\nReturn only the Terraform configuration code without explanations.")

	return prompt.String()
}

// extractTerraformCode extracts Terraform code from AI response
func extractTerraformCode(content string) string {
	// Remove markdown code blocks if present
	if strings.Contains(content, "```") {
		lines := strings.Split(content, "\n")
		var result []string
		inCodeBlock := false

		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = !inCodeBlock
				continue
			}
			if inCodeBlock {
				result = append(result, line)
			}
		}

		if len(result) > 0 {
			return strings.Join(result, "\n")
		}
	}

	return content
}

// GenerateTerraform creates Terraform configuration from parsed input
func (p *OpenAIProvider) GenerateTerraform(input *nlp.ParsedInput) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input cannot be nil")
	}

	// Create context for AI generation
	context := p.buildContext(input)

	// Generate configuration using AI
	config, err := p.generateWithAI(context)
	if err != nil {
		return "", fmt.Errorf("failed to generate configuration: %w", err)
	}

	// Validate and format the generated configuration
	formattedConfig, err := p.validateAndFormat(config)
	if err != nil {
		// Log error but don't fail completely - return raw config
		fmt.Printf("Warning: Failed to format configuration: %v\n", err)
		return config, nil
	}

	return formattedConfig, nil
}
