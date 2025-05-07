package cmd

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/JacobAndrewSmith92/gobuddy/testutils"
	// "github.com/JacobAndrewSmith/gobuddy/test_utils" // Removed as it is not available
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write valid JSON to the temp file
	configData := Config{
		Token:     "test-token",
		Workspace: "test-workspace",
		Protected: Protected{
			Pipeline: "test-pipeline",
			Branch:   "test-branch",
		},
	}
	data, _ := json.Marshal(configData)
	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Override the configFilePath for testing
	configFilePath = tempFile.Name()

	// Test loadConfig
	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() returned an error: %v", err)
	}

	if config.Token != "test-token" || config.Workspace != "test-workspace" ||
		config.Protected.Pipeline != "test-pipeline" || config.Protected.Branch != "test-branch" {
		t.Errorf("loadConfig() returned incorrect data: %+v", config)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Override the configFilePath for testing
	configFilePath = tempFile.Name()

	// Test saveConfig
	configData := Config{
		Token:     "test-token",
		Workspace: "test-workspace",
		Protected: Protected{
			Pipeline: "test-pipeline",
			Branch:   "test-branch",
		},
	}
	saveConfig(configData)

	// Read the file and verify its contents
	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	var loadedConfig Config
	if err := json.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loadedConfig != configData {
		t.Errorf("saveConfig() wrote incorrect data: %+v", loadedConfig)
	}
}

func TestConfigGetCmd(t *testing.T) {
	// Create a temporary file with a valid config
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	configFilePath = tempFile.Name()
	configData := Config{
		Token:     "test-token",
		Workspace: "test-workspace",
		Protected: Protected{
			Pipeline: "test-pipeline",
			Branch:   "test-branch",
		},
	}
	saveConfig(configData)

	output := testutils.Capture(func() {
		configGetCmd.Run(nil, nil)
	})

	if !strings.Contains(output, "test-token") || !strings.Contains(output, "test-workspace") ||
		!strings.Contains(output, "test-pipeline") || !strings.Contains(output, "test-branch") {
		t.Errorf("configGetCmd output is incorrect: %s", output)
	}
}

func TestHandleMissingConfig(t *testing.T) {
	// Mock user input for "yes"
	input := "yes\n"
	r, w, _ := os.Pipe()
	_, _ = w.Write([]byte(input))
	w.Close()
	originalStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = originalStdin }()

	// Mock setConfig to avoid actual prompts
	// setConfig := func(_, _, _, _ string) {}

	// Capture output
	output := testutils.Capture(func() {
		handleMissingConfig()
	})

	// Assert output contains expected text
	if !strings.Contains(output, "No configuration found.") {
		t.Errorf("Expected 'No configuration found.' in output, got: %s", output)
	}
	if !strings.Contains(output, "Would you like to create one?") {
		t.Errorf("Expected 'Would you like to create one?' in output, got: %s", output)
	}
}
