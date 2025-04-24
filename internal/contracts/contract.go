package contracts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateConfig represents the structure of the template configuration file
type TemplateConfig struct {
	Templates map[string]struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Options     map[string]struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Required    bool   `json:"required"`
			Default     any    `json:"default,omitempty"`
		} `json:"options"`
	} `json:"templates"`
}

// CreateOptions holds the options for contract creation
type CreateOptions struct {
	TemplateName  string
	ContractName  string
	OutputDir     string
	TemplateFlags map[string]interface{}
}

// Create generates a new contract from a template
func Create(opts CreateOptions) error {
	// Get the user's current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Handle output directory
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = cwd
	} else if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(cwd, outputDir)
	}

	// Read template configuration
	configPath := filepath.Join("templates", "contracts", "config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read template config: %w", err)
	}

	var config TemplateConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse template config: %w", err)
	}

	// Validate template name
	templateConfig, ok := config.Templates[opts.TemplateName]
	if !ok {
		return fmt.Errorf("invalid template name. Available templates: %s", strings.Join(getTemplateNames(config), ", "))
	}

	// Create template data with defaults
	templateData := map[string]interface{}{
		"ContractName": opts.ContractName,
	}

	// Set default values for all options
	for name, opt := range templateConfig.Options {
		if opt.Default != nil {
			templateData[name] = opt.Default
		}
	}

	// Override defaults with provided flags
	for name, value := range opts.TemplateFlags {
		templateData[name] = value
	}

	// Read template file
	templatePath := filepath.Join("templates", "contracts", fmt.Sprintf("%s.sol.tmpl", opts.TemplateName))
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file in the specified directory
	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.sol", opts.ContractName))
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Parse and execute template
	tmpl, err := template.New("contract").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(outputFile, templateData); err != nil {
		return fmt.Errorf("failed to generate contract: %w", err)
	}

	return nil
}

// GetTemplateNames returns a list of available template names
func getTemplateNames(config TemplateConfig) []string {
	names := make([]string, 0, len(config.Templates))
	for name := range config.Templates {
		names = append(names, name)
	}
	return names
}
