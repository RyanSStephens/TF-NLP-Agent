# TF-NLP-Agent

A Terraform configuration generator that uses natural language processing and AI to create functional Terraform configurations from simple English descriptions.

## Overview

TF-NLP-Agent bridges the gap between infrastructure requirements expressed in natural language and production-ready Terraform configurations. By leveraging advanced AI models, it transforms plain English descriptions into properly structured, secure, and deployable Infrastructure as Code.

## Features

- **Natural Language Processing**: Parse infrastructure requirements from plain English
- **Multi-Cloud Support**: Generate configurations for AWS, Azure, and Google Cloud
- **Security Scanning**: Built-in security analysis with comprehensive rule set
- **Web Interface**: REST API and web UI for easy integration
- **CLI Tool**: Command-line interface for automation and CI/CD pipelines
- **Terraform Validation**: Syntax and logical validation of generated configurations
- **Template System**: Extensible template system for custom resource types
- **Configuration Management**: YAML-based configuration with environment support

## Quick Start

1. **Install the CLI:**
   ```bash
   go install github.com/RyanSStephens/TF-NLP-Agent/cmd/agent@latest
   ```

2. **Set up your configuration:**
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your OpenAI API key and preferences
   ```

### Basic Usage

#### CLI
```bash
# Generate Terraform config from natural language
./tf-nlp-agent generate "Create an AWS VPC with public and private subnets"

# Validate generated configuration
./tf-nlp-agent validate output.tf

# Start web server
./tf-nlp-agent serve --port 8080
```

#### Web Interface
```bash
# Start the web server
./tf-nlp-agent serve

# Open browser to http://localhost:8080
```

## Configuration

Create a `config.yaml` file:

```yaml
ai:
  provider: "openai"
  model: "gpt-4"
  api_key: "${OPENAI_API_KEY}"

terraform:
  default_provider: "aws"
  validate: true
  format: true

security:
  scan_enabled: true
  fail_on_high: true

templates:
  path: "./templates"
  custom_path: "./custom-templates"
```

## Examples

### Example 1: Simple Web Application Infrastructure
```
Input: "I need a web application setup with a load balancer, auto-scaling group, and RDS database in AWS"

Output: Complete Terraform configuration with:
- Application Load Balancer
- Auto Scaling Group with Launch Template
- RDS MySQL instance with proper security groups
- VPC with public/private subnets
- Security groups with least privilege access
```

### Example 2: Kubernetes Cluster
```
Input: "Create a production-ready Kubernetes cluster on GCP with 3 nodes"

Output: Terraform configuration for:
- GKE cluster with proper node pools
- Network configuration
- IAM roles and service accounts
- Monitoring and logging setup
```

## Architecture

The TF-NLP-Agent follows a modular architecture:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   NLP Parser    │───▶│   AI Provider   │───▶│ Terraform Gen   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Security Scanner│    │  Web Interface  │    │   CLI Tool      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+ (for web UI development)
- Terraform 1.0+

### Building
```bash
# Build CLI tool
go build -o tf-nlp-agent ./cmd/agent

# Build web assets
cd web && npm install && npm run build

# Run tests
go test ./...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Support for more cloud providers (Oracle Cloud, DigitalOcean)
- [ ] Integration with Terraform Cloud
- [ ] Cost optimization suggestions
- [ ] Infrastructure drift detection
- [ ] Multi-language support
- [ ] Plugin system for custom providers
- [ ] Integration with CI/CD pipelines
- [ ] Advanced security compliance checks

## Support

- Documentation: [Wiki](https://github.com/RyanSStephens/TF-NLP-Agent/wiki)
- Issues: [GitHub Issues](https://github.com/RyanSStephens/TF-NLP-Agent/issues)
- Discussions: [GitHub Discussions](https://github.com/RyanSStephens/TF-NLP-Agent/discussions) 