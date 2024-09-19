package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Configuration structure
type Config struct {
	Token     string `json:"token"`
	Workspace string `json:"workspace"`
}

var configFilePath = filepath.Join(os.Getenv("HOME"), ".gobuddy_config.json")

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your token and workspace",
	Long:  `This command allows you to configure your authorization token and workspace for Buddy API requests.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [token|workspace] [value]",
	Short: "Set or update your token and workspace",
	Long:  `This subcommand allows you to set or update your authorization token and workspace. Pass "token" or "workspace" followed by the value to update.`,
	Args:  cobra.MinimumNArgs(0), // No minimum args; prompts if args are missing
	Run: func(_ *cobra.Command, args []string) {
		setConfigFromArgs(args)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current configuration",
	Long:  `This subcommand will display the currently saved token and workspace.`,
	Run: func(_ *cobra.Command, _ []string) {
		config, err := loadConfig()
		if err != nil && os.IsNotExist(err) {
			handleMissingConfig()
			return
		} else if err != nil {
			log.Fatalf("Failed to load configuration: %v\n", err)
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println(bold("Current Configuration:"))
		fmt.Printf("Token: %s\n", cyan(config.Token))
		fmt.Printf("Workspace: %s\n", cyan(config.Workspace))
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the current configuration",
	Long:  `This subcommand will clear the currently saved configuration.`,
	Run: func(_ *cobra.Command, _ []string) {
		confirmReset()
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}

// Save the configuration
func saveConfig(config Config) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v\n", err)
	}

	err = os.WriteFile(configFilePath, data, 0600)
	if err != nil {
		log.Fatalf("Failed to write config file: %v\n", err)
	}
}

// Load the configuration
func loadConfig() (Config, error) {
	config := Config{}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Handle case where config doesn't exist during 'config get'
func handleMissingConfig() {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(red("No configuration found."))

	prompt := promptui.Prompt{
		Label: yellow("Would you like to create one? (yes/no)"),
		Validate: func(input string) error {
			if input != "yes" && input != "no" {
				return fmt.Errorf("please type 'yes' or 'no'")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if result == "yes" {
		setConfig("", "")
	} else {
		fmt.Println("No configuration created.")
	}
}

// Set or update configuration fields from arguments
func setConfigFromArgs(args []string) {
	config, err := loadConfig()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to load existing config: %v\n", err)
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	// Parse arguments and update fields
	if len(args) >= 2 {
		key := strings.ToLower(args[0])
		value := args[1]

		switch key {
		case "token":
			config.Token = value
			fmt.Printf("Token updated to: %s\n", yellow(value))
		case "workspace":
			config.Workspace = value
			fmt.Printf("Workspace updated to: %s\n", yellow(value))
		default:
			log.Fatalf("Invalid argument: %s. Use 'token' or 'workspace'.", key)
		}
	} else if len(args) == 0 {
		// Prompt for both token and workspace if no args are provided
		setConfig("", "")
		return
	} else {
		log.Fatalf("Invalid number of arguments. You must provide a key (token|workspace) and a value.")
	}

	saveConfig(config)
	fmt.Println(green("Configuration updated successfully!"))
}

// Prompt-based configuration setup
func setConfig(tokenFlag, workspaceFlag string) {
	config, err := loadConfig()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to load existing config: %v\n", err)
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	// Prompt for token if not provided
	if tokenFlag == "" && config.Token == "" {
		tokenPrompt := promptui.Prompt{
			Label: yellow("Enter your Buddy API token"),
			Mask:  '*',
		}
		token, err := tokenPrompt.Run()
		if err != nil {
			log.Fatalf("Failed to read token: %v\n", err)
		}
		config.Token = token
	}

	// Prompt for workspace if not provided
	if workspaceFlag == "" && config.Workspace == "" {
		workspacePrompt := promptui.Prompt{
			Label: yellow("Enter your Buddy workspace"),
		}
		workspace, err := workspacePrompt.Run()
		if err != nil {
			log.Fatalf("Failed to read workspace: %v\n", err)
		}
		config.Workspace = workspace
	}

	saveConfig(config)
	fmt.Println(green("Configuration saved successfully!"))
}

// Confirm reset
func confirmReset() {
	confirm := promptui.Prompt{
		Label: "Are you sure you want to reset the configuration? (yes/no)",
		Validate: func(input string) error {
			if input != "yes" && input != "no" {
				return fmt.Errorf("please type 'yes' or 'no'")
			}
			return nil
		},
	}

	result, err := confirm.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if result == "yes" {
		err := os.Remove(configFilePath)
		if err != nil {
			log.Fatalf("Failed to reset configuration: %v\n", err)
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(green("Configuration has been reset."))
	} else {
		fmt.Println("Reset canceled.")
	}
}
