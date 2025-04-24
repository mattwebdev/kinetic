package cli

import (
	"fmt"
	"strings"

	"github.com/kinetic-dev/kinetic/internal/contracts"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Manage smart contracts",
	Long:  `Commands for creating and managing smart contracts.`,
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

func runContractCreate(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	contractName := args[1]

	// Get output directory from flag
	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return fmt.Errorf("failed to get output directory flag: %w", err)
	}

	// Collect template flags
	templateFlags := make(map[string]interface{})
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Name != "output-dir" {
			templateFlags[camelCase(f.Name)] = f.Value
		}
	})

	// Create contract using the contracts package
	opts := contracts.CreateOptions{
		TemplateName:  templateName,
		ContractName:  contractName,
		OutputDir:     outputDir,
		TemplateFlags: templateFlags,
	}

	if err := contracts.Create(opts); err != nil {
		return err
	}

	fmt.Printf("Contract created successfully: %s\n", contractName)
	return nil
}

func camelCase(s string) string {
	parts := strings.Split(s, "-")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func init() {
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
