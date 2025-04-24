package contracts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "contract-test")
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

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
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

	// Create contract with options
	opts := CreateOptions{
		TemplateName: "Basic",
		ContractName: "MyContract",
		OutputDir:    outputDir,
		TemplateFlags: map[string]interface{}{
			"HasMaxSupply": true,
			"MaxSupply":    "2000000",
		},
	}

	// Run the command
	if err := Create(opts); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check if output file exists
	outputFile := filepath.Join(outputDir, "MyContract.sol")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read and verify output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MyContract {
    uint256 public constant MAX_SUPPLY = 2000000;
}`

	if strings.TrimSpace(string(content)) != strings.TrimSpace(expectedContent) {
		t.Errorf("Output content does not match expected.\nGot:\n%s\nWant:\n%s", content, expectedContent)
	}
}
