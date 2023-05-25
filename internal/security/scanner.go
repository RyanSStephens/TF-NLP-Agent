package security

import (
	"regexp"
	"strings"
)

// Issue represents a security issue found in Terraform configuration
type Issue struct {
	Severity    string
	Message     string
	Resource    string
	Line        int
	Rule        string
	Remediation string
}

// Scanner handles security scanning of Terraform configurations
type Scanner struct {
	rules []SecurityRule
}

// SecurityRule represents a security rule to check
type SecurityRule struct {
	ID          string
	Name        string
	Severity    string
	Pattern     *regexp.Regexp
	Message     string
	Remediation string
}

// NewScanner creates a new security scanner with default rules
func NewScanner() *Scanner {
	scanner := &Scanner{
		rules: []SecurityRule{},
	}

	scanner.loadDefaultRules()
	return scanner
}

// Scan analyzes a Terraform configuration for security issues
func (s *Scanner) Scan(config string) ([]Issue, error) {
	var issues []Issue

	lines := strings.Split(config, "\n")

	for lineNum, line := range lines {
		for _, rule := range s.rules {
			if rule.Pattern.MatchString(line) {
				issue := Issue{
					Severity:    rule.Severity,
					Message:     rule.Message,
					Resource:    extractResourceName(line),
					Line:        lineNum + 1,
					Rule:        rule.ID,
					Remediation: rule.Remediation,
				}
				issues = append(issues, issue)
			}
		}
	}

	// Additional context-aware checks
	contextIssues := s.performContextualScans(config)
	issues = append(issues, contextIssues...)

	return issues, nil
}

// loadDefaultRules loads the default security rules
func (s *Scanner) loadDefaultRules() {
	rules := []SecurityRule{
		{
			ID:          "SEC001",
			Name:        "Public S3 Bucket",
			Severity:    "HIGH",
			Pattern:     regexp.MustCompile(`acl\s*=\s*"public-read"`),
			Message:     "S3 bucket configured with public read access",
			Remediation: "Remove public ACL and use bucket policies for controlled access",
		},
		{
			ID:          "SEC002",
			Name:        "Unencrypted Storage",
			Severity:    "MEDIUM",
			Pattern:     regexp.MustCompile(`resource\s+"aws_s3_bucket"`),
			Message:     "S3 bucket may not have encryption enabled",
			Remediation: "Enable server-side encryption for S3 buckets",
		},
		{
			ID:          "SEC003",
			Name:        "Open Security Group",
			Severity:    "CRITICAL",
			Pattern:     regexp.MustCompile(`cidr_blocks\s*=\s*\["0\.0\.0\.0/0"\]`),
			Message:     "Security group allows access from anywhere (0.0.0.0/0)",
			Remediation: "Restrict CIDR blocks to specific IP ranges",
		},
		{
			ID:          "SEC004",
			Name:        "Unencrypted EBS Volume",
			Severity:    "MEDIUM",
			Pattern:     regexp.MustCompile(`resource\s+"aws_ebs_volume"`),
			Message:     "EBS volume may not have encryption enabled",
			Remediation: "Enable encryption for EBS volumes",
		},
		{
			ID:          "SEC005",
			Name:        "Public RDS Instance",
			Severity:    "HIGH",
			Pattern:     regexp.MustCompile(`publicly_accessible\s*=\s*true`),
			Message:     "RDS instance is publicly accessible",
			Remediation: "Set publicly_accessible to false for RDS instances",
		},
		{
			ID:          "SEC006",
			Name:        "Weak Password Policy",
			Severity:    "MEDIUM",
			Pattern:     regexp.MustCompile(`password\s*=\s*"[^"]{1,7}"`),
			Message:     "Password appears to be too short",
			Remediation: "Use strong passwords with at least 8 characters",
		},
		{
			ID:          "SEC007",
			Name:        "Hardcoded Secrets",
			Severity:    "CRITICAL",
			Pattern:     regexp.MustCompile(`(password|secret|key)\s*=\s*"[^$][^"]*"`),
			Message:     "Potential hardcoded secret or password",
			Remediation: "Use variables or AWS Secrets Manager for sensitive data",
		},
		{
			ID:          "SEC008",
			Name:        "Missing HTTPS",
			Severity:    "MEDIUM",
			Pattern:     regexp.MustCompile(`protocol\s*=\s*"HTTP"`),
			Message:     "Load balancer listener using HTTP instead of HTTPS",
			Remediation: "Use HTTPS protocol for load balancer listeners",
		},
		{
			ID:          "SEC009",
			Name:        "Default VPC Usage",
			Severity:    "LOW",
			Pattern:     regexp.MustCompile(`default\s*=\s*true.*vpc`),
			Message:     "Using default VPC may not follow security best practices",
			Remediation: "Create custom VPC with proper network segmentation",
		},
		{
			ID:          "SEC010",
			Name:        "Missing Backup",
			Severity:    "MEDIUM",
			Pattern:     regexp.MustCompile(`backup_retention_period\s*=\s*0`),
			Message:     "Database backup retention period is set to 0",
			Remediation: "Enable automated backups with appropriate retention period",
		},
	}

	s.rules = rules
}

// performContextualScans performs more complex security checks that require context
func (s *Scanner) performContextualScans(config string) []Issue {
	var issues []Issue

	// Check for missing encryption on storage resources
	if strings.Contains(config, "aws_s3_bucket") && !strings.Contains(config, "server_side_encryption") {
		issues = append(issues, Issue{
			Severity:    "MEDIUM",
			Message:     "S3 bucket missing server-side encryption configuration",
			Rule:        "SEC011",
			Remediation: "Add server_side_encryption_configuration block",
		})
	}

	// Check for missing versioning on S3 buckets
	if strings.Contains(config, "aws_s3_bucket") && !strings.Contains(config, "versioning") {
		issues = append(issues, Issue{
			Severity:    "LOW",
			Message:     "S3 bucket missing versioning configuration",
			Rule:        "SEC012",
			Remediation: "Enable versioning for S3 buckets",
		})
	}

	// Check for missing MFA delete on S3 buckets
	if strings.Contains(config, "aws_s3_bucket") && !strings.Contains(config, "mfa_delete") {
		issues = append(issues, Issue{
			Severity:    "LOW",
			Message:     "S3 bucket missing MFA delete protection",
			Rule:        "SEC013",
			Remediation: "Enable MFA delete for S3 buckets containing sensitive data",
		})
	}

	// Check for EC2 instances without security groups
	if strings.Contains(config, "aws_instance") && !strings.Contains(config, "security_groups") && !strings.Contains(config, "vpc_security_group_ids") {
		issues = append(issues, Issue{
			Severity:    "HIGH",
			Message:     "EC2 instance missing security group configuration",
			Rule:        "SEC014",
			Remediation: "Assign appropriate security groups to EC2 instances",
		})
	}

	// Check for RDS instances without encryption
	if strings.Contains(config, "aws_db_instance") && !strings.Contains(config, "storage_encrypted") {
		issues = append(issues, Issue{
			Severity:    "MEDIUM",
			Message:     "RDS instance missing storage encryption",
			Rule:        "SEC015",
			Remediation: "Enable storage encryption for RDS instances",
		})
	}

	return issues
}

// extractResourceName extracts the resource name from a Terraform line
func extractResourceName(line string) string {
	// Pattern to match resource declarations: resource "type" "name"
	re := regexp.MustCompile(`resource\s+"([^"]+)"\s+"([^"]+)"`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 3 {
		return matches[1] + "." + matches[2]
	}

	return ""
}

// AddCustomRule adds a custom security rule to the scanner
func (s *Scanner) AddCustomRule(rule SecurityRule) {
	s.rules = append(s.rules, rule)
}

// GetRules returns all loaded security rules
func (s *Scanner) GetRules() []SecurityRule {
	return s.rules
}
