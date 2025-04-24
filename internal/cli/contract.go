package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Manage smart contracts",
	Long:  `Commands for creating and managing smart contracts.`,
}

var contractListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available contract templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read template configuration
		configPath := "templates/contracts/config.json"
		configData, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read template config: %w", err)
		}

		var config TemplateConfig
		if err := json.Unmarshal(configData, &config); err != nil {
			return fmt.Errorf("failed to parse template config: %w", err)
		}

		fmt.Println("Available templates:")
		for name, tmpl := range config.Templates {
			fmt.Printf("- %s: %s\n", name, tmpl.Description)
		}
		return nil
	},
}

var contractCreateCmd = &cobra.Command{
	Use:   "create [template] [name]",
	Short: "Create a new contract from template",
	Long: `Create a new smart contract from a template.
Available templates: ERC20, ERC721, Basic

Example:
  kinetic contract create ERC20 MyToken --output-dir ./contracts
  kinetic contract create ERC721 MyNFT --output-dir ./src/contracts --has-max-supply
  kinetic contract create Basic MyContract --output-dir ./solidity`,
	Args: cobra.ExactArgs(2),
	RunE: runContractCreate,
}

var contractDeployCmd = &cobra.Command{
	Use:   "deploy [contract]",
	Short: "Deploy a contract",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contract := args[0]
		network, _ := cmd.Flags().GetString("network")
		fmt.Printf("Deploying contract '%s' to '%s'...\n", contract, network)
		// TODO: Implement contract deployment logic
		return nil
	},
}

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

func runContractCreate(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	contractName := args[1]

	// Get the user's current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Get output directory from flag, default to current working directory if not specified
	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return fmt.Errorf("failed to get output directory flag: %w", err)
	}
	if outputDir == "" {
		outputDir = cwd
	} else {
		// If relative path, make it relative to current working directory
		if !filepath.IsAbs(outputDir) {
			outputDir = filepath.Join(cwd, outputDir)
		}
	}

	// Read template configuration - use absolute path from project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(cwd))) // Go up to project root
	configPath := filepath.Join(projectRoot, "templates", "contracts", "config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read template config: %w", err)
	}

	var config TemplateConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse template config: %w", err)
	}

	// Validate template name
	templateConfig, ok := config.Templates[templateName]
	if !ok {
		return fmt.Errorf("invalid template name. Available templates: %s", strings.Join(getTemplateNames(config), ", "))
	}

	// Create template data with defaults
	templateData := map[string]interface{}{
		"ContractName": contractName,
	}

	// Set default values for all options
	for name, opt := range templateConfig.Options {
		if opt.Default != nil {
			templateData[name] = opt.Default
		}
	}

	// Override defaults with flags
	for name, opt := range templateConfig.Options {
		if name == "ContractName" {
			continue
		}

		if opt.Type == "boolean" {
			if val, _ := cmd.Flags().GetBool(kebabCase(name)); cmd.Flags().Changed(kebabCase(name)) {
				templateData[name] = val
			}
		} else if opt.Type == "string" {
			if val, _ := cmd.Flags().GetString(kebabCase(name)); cmd.Flags().Changed(kebabCase(name)) {
				templateData[name] = val
			}
		}
	}

	// Read template file - use absolute path from project root
	templatePath := filepath.Join(projectRoot, "templates", "contracts", fmt.Sprintf("%s.sol.tmpl", templateName))
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file in the specified directory
	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.sol", contractName))
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

	// Display the path relative to where the user ran the command
	relPath, err := filepath.Rel(cwd, outputPath)
	if err != nil {
		// Fall back to absolute path if relative path fails
		relPath = outputPath
	}
	if relPath == "" {
		relPath = filepath.Base(outputPath)
	}

	fmt.Printf("Contract created successfully: %s\n", relPath)
	return nil
}

func getTemplateNames(config TemplateConfig) []string {
	names := make([]string, 0, len(config.Templates))
	for name := range config.Templates {
		names = append(names, name)
	}
	return names
}

func kebabCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "Has", "has-"))
}

func init() {
	contractCmd.AddCommand(contractListCmd)
	contractCmd.AddCommand(contractCreateCmd)
	contractCmd.AddCommand(contractDeployCmd)

	// Update output directory flag description
	contractCreateCmd.Flags().StringP("output-dir", "o", "", "Output directory for generated contracts (default: current directory)")

	// Add template-specific flags
	contractCreateCmd.Flags().Bool("has-cap", false, "Add maximum supply cap (ERC20)")
	contractCreateCmd.Flags().Bool("has-max-supply", false, "Add maximum supply limit (ERC721)")
	contractCreateCmd.Flags().Bool("has-base-uri", true, "Use base URI for metadata (ERC721)")
	contractCreateCmd.Flags().Bool("has-custom-uri", false, "Allow custom URI per token (ERC721)")
	contractCreateCmd.Flags().Bool("has-storage", true, "Include storage functionality (Basic)")
	contractCreateCmd.Flags().Bool("has-events", true, "Include events (Basic)")
	contractCreateCmd.Flags().Bool("has-whitelist", false, "Include whitelist functionality (Basic)")
	contractCreateCmd.Flags().Bool("has-initial-setup", false, "Include initial setup code (Basic)")
	contractCreateCmd.Flags().Bool("has-emergency-stop", false, "Include emergency stop functionality (Basic)")
	contractCreateCmd.Flags().Bool("has-upgradeable", false, "Include upgrade functionality (Basic)")
	contractCreateCmd.Flags().Bool("has-emergency-withdraw", false, "Include emergency withdraw functionality (Basic)")

	contractCreateCmd.Flags().Bool("is-mintable", false, "Allow minting new tokens")
	contractCreateCmd.Flags().Bool("is-burnable", false, "Allow burning tokens")
	contractCreateCmd.Flags().Bool("is-pausable", false, "Allow pausing transfers")
	contractCreateCmd.Flags().Bool("only-owner-can-mint", true, "Only owner can mint tokens (ERC721)")

	contractDeployCmd.Flags().StringP("network", "n", "local", "Target network (local, fuji, mainnet)")
	contractDeployCmd.Flags().StringP("private-key", "k", "", "Private key for deployment")
}
