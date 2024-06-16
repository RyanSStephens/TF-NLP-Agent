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

// GenerateConfig creates Terraform configuration from parsed natural language input
func (p *OpenAIProvider) GenerateConfig(parsed *nlp.ParsedInput) (string, error) {
	if parsed == nil {
		return "", fmt.Errorf("parsed input cannot be nil")
	}

	if p.client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	// Build the prompt for the AI model
	prompt := p.buildPrompt(parsed)

	// Make the API call to OpenAI
	response, err := p.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(p.model),
	})

	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from OpenAI")
	}

	// Extract the generated Terraform configuration
	content := response.Choices[0].Message.Content

	// Clean up the response (remove markdown formatting if present)
	content = p.cleanResponse(content)

	return content, nil
}

// buildPrompt constructs the prompt for the AI model
func (p *OpenAIProvider) buildPrompt(parsed *nlp.ParsedInput) string {
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

// cleanResponse cleans up the response from OpenAI
func (p *OpenAIProvider) cleanResponse(content string) string {
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
