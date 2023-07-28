package main

import (
	"fmt"
	"os"

	"github.com/RyanSStephens/TF-NLP-Agent/internal/ai"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/nlp"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/security"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/terraform"
	"github.com/RyanSStephens/TF-NLP-Agent/internal/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version = "1.0.0"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tf-nlp-agent",
	Short: "Terraform Natural Language Processing Agent",
	Long: `TF-NLP-Agent is a tool that converts natural language descriptions 
into functional Terraform configurations using AI and NLP techniques.`,
	Version: version,
}

var generateCmd = &cobra.Command{
	Use:   "generate [description]",
	Short: "Generate Terraform configuration from natural language",
	Long: `Generate Terraform configuration from a natural language description.
	
Example:
  tf-nlp-agent generate "Create an AWS VPC with public and private subnets"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		description := args[0]

		// Initialize components
		aiProvider := ai.NewProvider(viper.GetString("ai.provider"))
		nlpEngine := nlp.NewEngine()
		tfGenerator := terraform.NewGenerator()
		securityScanner := security.NewScanner()

		// Process the description
		fmt.Printf("Processing: %s\n", description)

		// Parse the natural language input
		parsed, err := nlpEngine.Parse(description)
		if err != nil {
			return fmt.Errorf("failed to parse description: %w", err)
		}

		// Generate Terraform configuration using AI
		config, err := aiProvider.GenerateConfig(parsed)
		if err != nil {
			return fmt.Errorf("failed to generate configuration: %w", err)
		}

		// Validate and format the configuration
		validated, err := tfGenerator.Validate(config)
		if err != nil {
			return fmt.Errorf("failed to validate configuration: %w", err)
		}

		// Security scan if enabled
		if viper.GetBool("security.scan_enabled") {
			issues, err := securityScanner.Scan(validated)
			if err != nil {
				return fmt.Errorf("security scan failed: %w", err)
			}

			if len(issues) > 0 {
				fmt.Println("Security issues found:")
				for _, issue := range issues {
					fmt.Printf("  - %s: %s\n", issue.Severity, issue.Message)
				}

				if viper.GetBool("security.fail_on_high") && hasHighSeverityIssues(issues) {
					return fmt.Errorf("high severity security issues found")
				}
			}
		}

		// Output the configuration
		outputFile := cmd.Flag("output").Value.String()
		if outputFile != "" {
			err = os.WriteFile(outputFile, []byte(validated), 0644)
			if err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Printf("Configuration written to: %s\n", outputFile)
		} else {
			fmt.Println("\nGenerated Terraform Configuration:")
			fmt.Println("=" + string(make([]rune, 50)) + "=")
			fmt.Println(validated)
		}

		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  "Start the web server for interactive Terraform configuration generation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		port := cmd.Flag("port").Value.String()

		server := web.NewServer()
		fmt.Printf("Starting web server on port %s\n", port)
		fmt.Printf("Open your browser to http://localhost:%s\n", port)

		return server.Start(":" + port)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a Terraform configuration file",
	Long:  "Validate syntax and security of an existing Terraform configuration file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		content, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		tfGenerator := terraform.NewGenerator()
		securityScanner := security.NewScanner()

		// Validate syntax
		_, err = tfGenerator.Validate(string(content))
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		// Security scan
		issues, err := securityScanner.Scan(string(content))
		if err != nil {
			return fmt.Errorf("security scan failed: %w", err)
		}

		fmt.Printf("Validation successful for: %s\n", filename)

		if len(issues) > 0 {
			fmt.Println("Security issues found:")
			for _, issue := range issues {
				fmt.Printf("  - %s: %s\n", issue.Severity, issue.Message)
			}
		} else {
			fmt.Println("No security issues found.")
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tf-nlp-agent.yaml)")

	// Generate command flags
	generateCmd.Flags().StringP("output", "o", "", "output file for generated configuration")
	generateCmd.Flags().StringP("provider", "p", "aws", "cloud provider (aws, azure, gcp)")

	// Serve command flags
	serveCmd.Flags().StringP("port", "p", "8080", "port to run the web server on")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(validateCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tf-nlp-agent")
	}

	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-4")
	viper.SetDefault("terraform.default_provider", "aws")
	viper.SetDefault("terraform.validate", true)
	viper.SetDefault("terraform.format", true)
	viper.SetDefault("security.scan_enabled", true)
	viper.SetDefault("security.fail_on_high", false)
	viper.SetDefault("templates.path", "./templates")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}
}

func hasHighSeverityIssues(issues []security.Issue) bool {
	for _, issue := range issues {
		if issue.Severity == "HIGH" || issue.Severity == "CRITICAL" {
			return true
		}
	}
	return false
}
