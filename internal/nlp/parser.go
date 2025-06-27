package nlp

import (
	"fmt"
	"regexp"
	"strings"
)

// ParsedInput represents the structured output from natural language parsing
type ParsedInput struct {
	OriginalText  string
	CloudProvider string
	Resources     []Resource
	Requirements  []string
	Intent        string
}

// Resource represents an identified infrastructure resource
type Resource struct {
	Type       string
	Name       string
	Properties map[string]string
	Attributes []string
}

// Engine handles natural language processing
type Engine struct {
	cloudProviders map[string][]string
	resourceTypes  map[string][]string
}

// NewEngine creates a new NLP engine
func NewEngine() *Engine {
	return &Engine{
		cloudProviders: map[string][]string{
			"aws":   {"aws", "amazon", "ec2", "s3", "rds", "vpc", "lambda"},
			"azure": {"azure", "microsoft", "vm", "storage", "sql"},
			"gcp":   {"gcp", "google", "compute", "storage", "cloud"},
		},
		resourceTypes: map[string][]string{
			"compute":    {"vm", "instance", "server", "compute", "ec2"},
			"storage":    {"storage", "bucket", "s3", "blob", "disk"},
			"network":    {"vpc", "network", "subnet", "security group", "firewall", "load balancer", "alb", "nlb"},
			"database":   {"database", "db", "rds", "sql", "mysql", "postgres", "mongodb"},
			"container":  {"container", "kubernetes", "k8s", "docker", "ecs", "aks", "gke"},
			"serverless": {"lambda", "function", "serverless", "azure functions", "cloud functions"},
		},
	}
}

// Parse processes natural language input and extracts structured information
func (e *Engine) Parse(input string) (*ParsedInput, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	parsed := &ParsedInput{
		OriginalText: input,
		Resources:    []Resource{},
		Requirements: []string{},
	}

	// Detect cloud provider
	parsed.CloudProvider = e.detectCloudProvider(input)

	// Extract resources
	parsed.Resources = e.extractResources(input)

	// Extract requirements
	parsed.Requirements = e.extractRequirements(input)

	// Determine intent
	parsed.Intent = e.determineIntent(input)

	return parsed, nil
}

// detectCloudProvider identifies the cloud provider from the input
func (e *Engine) detectCloudProvider(input string) string {
	for provider, keywords := range e.cloudProviders {
		for _, keyword := range keywords {
			if strings.Contains(input, keyword) {
				return provider
			}
		}
	}
	return "aws" // Default to AWS
}

// extractResources identifies infrastructure resources mentioned in the input
func (e *Engine) extractResources(input string) []Resource {
	var resources []Resource

	for resourceType, keywords := range e.resourceTypes {
		for _, keyword := range keywords {
			if strings.Contains(input, keyword) {
				resource := Resource{
					Type:       resourceType,
					Name:       e.generateResourceName(resourceType),
					Properties: make(map[string]string),
					Attributes: []string{},
				}

				// Extract specific attributes based on context
				resource.Attributes = e.extractAttributes(input, resourceType)

				resources = append(resources, resource)
				break // Only add each resource type once
			}
		}
	}

	return resources
}

// extractRequirements identifies specific requirements from the input
func (e *Engine) extractRequirements(input string) []string {
	var requirements []string

	// Security requirements
	securityPatterns := []string{
		"secure", "security", "encrypted", "ssl", "tls", "https",
		"private", "public", "firewall", "access control",
	}

	for _, pattern := range securityPatterns {
		if strings.Contains(input, pattern) {
			requirements = append(requirements, "Security: "+pattern)
		}
	}

	// Scalability requirements
	scalabilityPatterns := []string{
		"scalable", "auto scaling", "high availability", "redundant",
		"multi-az", "multi-region", "load balanced",
	}

	for _, pattern := range scalabilityPatterns {
		if strings.Contains(input, pattern) {
			requirements = append(requirements, "Scalability: "+pattern)
		}
	}

	// Performance requirements
	performancePatterns := []string{
		"fast", "performance", "optimized", "cached", "cdn",
	}

	for _, pattern := range performancePatterns {
		if strings.Contains(input, pattern) {
			requirements = append(requirements, "Performance: "+pattern)
		}
	}

	// Extract numerical requirements (e.g., "3 servers", "100GB storage")
	numericRequirements := e.extractNumericRequirements(input)
	requirements = append(requirements, numericRequirements...)

	return requirements
}

// extractNumericRequirements finds numerical specifications in the input
func (e *Engine) extractNumericRequirements(input string) []string {
	var requirements []string

	// Pattern for numbers followed by units or resources
	patterns := []string{
		`(\d+)\s*(gb|tb|mb)\s*(storage|disk|memory|ram)`,
		`(\d+)\s*(cpu|core|vcpu)`,
		`(\d+)\s*(instance|server|vm|node)`,
		`(\d+)\s*(port|ports)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(input, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				requirements = append(requirements,
					fmt.Sprintf("Specification: %s %s %s", match[1], match[2], match[3]))
			}
		}
	}

	return requirements
}

// extractAttributes extracts specific attributes for a resource type
func (e *Engine) extractAttributes(input string, resourceType string) []string {
	var attributes []string

	switch resourceType {
	case "compute":
		if strings.Contains(input, "linux") || strings.Contains(input, "ubuntu") || strings.Contains(input, "centos") {
			attributes = append(attributes, "os:linux")
		}
		if strings.Contains(input, "windows") {
			attributes = append(attributes, "os:windows")
		}

	case "network":
		if strings.Contains(input, "public") {
			attributes = append(attributes, "access:public")
		}
		if strings.Contains(input, "private") {
			attributes = append(attributes, "access:private")
		}

	case "database":
		if strings.Contains(input, "mysql") {
			attributes = append(attributes, "engine:mysql")
		}
		if strings.Contains(input, "postgres") {
			attributes = append(attributes, "engine:postgresql")
		}
	}

	return attributes
}

// determineIntent identifies the primary intent of the request
func (e *Engine) determineIntent(input string) string {
	createWords := []string{"create", "setup", "build", "deploy", "provision"}
	modifyWords := []string{"update", "modify", "change", "scale", "resize"}
	deleteWords := []string{"delete", "remove", "destroy", "terminate"}

	for _, word := range createWords {
		if strings.Contains(input, word) {
			return "create"
		}
	}

	for _, word := range modifyWords {
		if strings.Contains(input, word) {
			return "modify"
		}
	}

	for _, word := range deleteWords {
		if strings.Contains(input, word) {
			return "delete"
		}
	}

	// Default intent if no specific action is detected
	return "create"
}

// generateResourceName creates a resource name based on the type
func (e *Engine) generateResourceName(resourceType string) string {
	switch resourceType {
	case "compute":
		return "main_instance"
	case "storage":
		return "main_storage"
	case "network":
		return "main_network"
	case "database":
		return "main_database"
	case "container":
		return "main_cluster"
	case "serverless":
		return "main_function"
	default:
		return "main_" + resourceType
	}
}
