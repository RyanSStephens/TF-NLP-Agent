package nlp

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	if engine.cloudProviders == nil {
		t.Error("cloudProviders not initialized")
	}

	if engine.resourceTypes == nil {
		t.Error("resourceTypes not initialized")
	}
}

func TestParse(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "AWS VPC request",
			input:    "Create an AWS VPC with public and private subnets",
			expected: "aws",
		},
		{
			name:     "Azure request",
			input:    "Set up Azure virtual network with load balancer",
			expected: "azure",
		},
		{
			name:     "GCP request",
			input:    "Deploy Google Cloud compute instances",
			expected: "google",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if result.CloudProvider != tt.expected {
				t.Errorf("Parse() CloudProvider = %v, want %v", result.CloudProvider, tt.expected)
			}
		})
	}
}

func TestDetectCloudProvider(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		input    string
		expected string
	}{
		{"create aws vpc", "aws"},
		{"azure virtual machine", "azure"},
		{"google cloud storage", "google"},
		{"random infrastructure", "aws"}, // default
	}

	for _, tt := range tests {
		result := engine.detectCloudProvider(tt.input)
		if result != tt.expected {
			t.Errorf("detectCloudProvider(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
