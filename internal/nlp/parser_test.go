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
		{"random infrastructure", "aws"}, // default fallback
	}

	for _, tt := range tests {
		result := engine.detectCloudProvider(tt.input)
		if result != tt.expected {
			t.Errorf("detectCloudProvider(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestExtractResources(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name     string
		input    string
		expected []string // resource types
	}{
		{
			name:     "VPC and database",
			input:    "create vpc with mysql database",
			expected: []string{"network", "database"},
		},
		{
			name:     "Compute instances",
			input:    "deploy ec2 instances with load balancer",
			expected: []string{"compute", "network"},
		},
		{
			name:     "Storage bucket",
			input:    "create s3 bucket for file storage",
			expected: []string{"storage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources := engine.extractResources(tt.input)

			if len(resources) != len(tt.expected) {
				t.Errorf("extractResources(%q) returned %d resources, want %d",
					tt.input, len(resources), len(tt.expected))
				return
			}

			for i, resource := range resources {
				if resource.Type != tt.expected[i] {
					t.Errorf("extractResources(%q)[%d].Type = %v, want %v",
						tt.input, i, resource.Type, tt.expected[i])
				}
			}
		})
	}
}

func TestDetermineIntent(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		input    string
		expected string
	}{
		{"create new vpc", "create"},
		{"build infrastructure", "create"},
		{"update existing database", "modify"},
		{"scale the application", "modify"},
		{"delete old resources", "delete"},
		{"remove the vpc", "delete"},
		{"setup monitoring", "create"}, // default
	}

	for _, tt := range tests {
		result := engine.determineIntent(tt.input)
		if result != tt.expected {
			t.Errorf("determineIntent(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
