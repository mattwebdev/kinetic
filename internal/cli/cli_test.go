package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// testCommand is a helper function to execute a command and capture its output
func testCommand(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return strings.TrimSpace(buf.String()), err
}

func TestContractCreateCommand(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template directory structure
	templatesDir := filepath.Join(tmpDir, "templates", "contracts")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	// Create test template file
	templateContent := `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract {{.ContractName}} {
    {{if .HasMaxSupply}}
    uint256 public constant MAX_SUPPLY = {{.MaxSupply}};
    {{end}}
}`
	if err := os.WriteFile(
		filepath.Join(templatesDir, "Basic.sol.tmpl"),
		[]byte(templateContent),
		0644,
	); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create test config file
	configContent := `{
		"templates": {
			"Basic": {
				"description": "A basic smart contract template",
				"options": {
					"HasMaxSupply": {
						"type": "boolean",
						"description": "Whether to include a max supply cap",
						"default": false
					},
					"MaxSupply": {
						"type": "string",
						"description": "The maximum supply cap",
						"default": "1000000"
					}
				}
			}
		}
	}`
	if err := os.WriteFile(
		filepath.Join(templatesDir, "config.json"),
		[]byte(configContent),
		0644,
	); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set up test environment
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(filepath.Dir(tmpDir)); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()

	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		wantFiles []string
	}{
		{
			name:      "basic contract without flags",
			args:      []string{"create", "Basic", "MyContract"},
			wantErr:   false,
			wantFiles: []string{"MyContract.sol"},
		},
		{
			name:      "contract with output directory",
			args:      []string{"create", "Basic", "MyContract2", "--output-dir", "contracts"},
			wantErr:   false,
			wantFiles: []string{"contracts/MyContract2.sol"},
		},
		{
			name:      "contract with max supply",
			args:      []string{"create", "Basic", "MyContract3", "--has-max-supply", "--max-supply", "2000000"},
			wantErr:   false,
			wantFiles: []string{"MyContract3.sol"},
		},
		{
			name:      "invalid template name",
			args:      []string{"create", "NonExistent", "MyContract"},
			wantErr:   true,
			wantFiles: nil,
		},
		{
			name:      "missing contract name",
			args:      []string{"create", "Basic"},
			wantErr:   true,
			wantFiles: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(contractCmd)

			// Execute the command
			_, err := testCommand(t, cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("command execution error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if expected files were created
			if tt.wantFiles != nil {
				for _, file := range tt.wantFiles {
					path := filepath.Join(tmpDir, file)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("expected file %s was not created", file)
					}
				}
			}
		})
	}
}

func TestContractDeployCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "deploy to local network",
			args:    []string{"deploy", "MyContract", "--network", "local"},
			wantErr: false,
		},
		{
			name:    "deploy to fuji network",
			args:    []string{"deploy", "MyContract", "--network", "fuji"},
			wantErr: false,
		},
		{
			name:    "missing contract name",
			args:    []string{"deploy"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(contractCmd)

			_, err := testCommand(t, cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("command execution error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "help flag",
			args:    []string{"--help"},
			wantErr: false,
		},
		{
			name:    "version flag",
			args:    []string{"--version"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testCommand(t, rootCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("command execution error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlagHandling(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "config flag",
			args:    []string{"--config", "custom-config.yaml"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "custom-config.yaml") {
					t.Error("custom config path not acknowledged")
				}
			},
		},
		{
			name:    "invalid flag",
			args:    []string{"--invalid-flag"},
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := testCommand(t, rootCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("command execution error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}
