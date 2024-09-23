package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	buddy "github.com/JacobAndrewSmith92/gobuddy/internal"
	"github.com/JacobAndrewSmith92/gobuddy/internal/util"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// var projects = []string{}
// var branches = []string{"master", "main", "development", "feature-1", "bugfix-42"}
// var pipelines = []string{"CD", "Deploy to Staging"}

var branchFlag string
var pipelineFlag string
var currentFlag bool

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [project]",
	Short: "Select a project, branch, and pipeline for deployment",
	Long:  `This command allows you to choose a project, a git branch, and a pipeline for deployment. The project can be provided as an argument, and the branch or pipeline can be provided via flags or interactively selected.`,
	Args:  cobra.MaximumNArgs(1), // Accept one optional argument for project
	Run: func(_ *cobra.Command, args []string) {
		var project, branch string
		var pipeline buddy.Pipeline
		config, err := loadConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v\n", err)
		}

		apiClient := buddy.NewBuddyClient(config.Token, config.Workspace)

		if err != nil {
			log.Fatalf("Error creating api client: %v", err)
		}

		if currentFlag {
			branch, project, err = util.GetBranchAndDirectory()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Printf("Using current project %s and branch: %s\n", project, branch)
		}

		if len(args) > 0 || (currentFlag && project != "") {
			if project == "" {
				project = args[0]
			}
			fmt.Printf("Looking up project: %s\n", project)
			projectFound, err := apiClient.FetchProjectByName(project)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Println("Project found.", project)
			project = projectFound.Name
		} else {
			projects, err := apiClient.FetchProjects()
			if err != nil {
				log.Fatalf("Error fetching projects 2: %v", err)
			}
			project = searchProject(projects)
		}

		if branchFlag != "" || currentFlag {
			if branch == "" {
				branch = branchFlag
			}
			fmt.Printf("Looking up branch: %s\n", branch)
			branchFound, err := apiClient.FetchBranchByName(project, branch)
			if err != nil {
				log.Fatalf("Error: %v", err.Errors)
			}
			log.Println("Branch found.", branchFound.Name)
			branch = branchFound.Name
		} else if branch == "" {
			branches, err := apiClient.FetchBranches(project)
			if err != nil {
				log.Fatalf("Error fetching branches: %v", err)
			}
			branch = searchBranch(branches)
		}

		if pipelineFlag != "" {
			fmt.Printf("Using pipeline passed as flag: %s\n", pipelineFlag)
			pipelineFound, err := apiClient.FetchPipelineByID(project, pipelineFlag)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Println("Pipeline found.", pipelineFound.ID, pipelineFound.Name)
			pipeline = *pipelineFound
			return
		} else {
			pipelines, err := apiClient.FetchPipelines(project)
			if err != nil {
				log.Fatalf("Error fetching pipelines: %v", err)
			}
			pipeline = searchPipeline(pipelines, branch)
		}

		if pipeline.Name == config.Protected.Pipeline {
			red := color.New(color.FgRed).SprintFunc()
			log.Fatalf(red("Error: Unable to deploy protected pipeline: %s"), config.Protected.Pipeline)
			return
		} else if branch == config.Protected.Branch {
			red := color.New(color.FgRed).SprintFunc()
			log.Fatalf(red("Error: Unable to deploy protected branch: %s"), config.Protected.Branch)
			return
		}

		// Final output with color
		cyan := color.New(color.FgCyan).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		log.Printf("You selected project: %s\n", cyan(bold(project)))
		log.Printf("You selected branch: %s\n", cyan(bold(branch)))
		log.Printf("You selected pipeline: %s(%s)", cyan(bold(pipeline.Name)), cyan(bold(pipeline.ID)))

		if !confirmDeployment() {
			log.Println("Deployment canceled.")
			return
		}

		// Proceed with deployment logic (e.g., calling the Buddy API)
		log.Println("Proceeding with deployment...")
		// Call the function to deploy the pipeline

		execution, err := apiClient.RunPipeline(project, pipeline.ID, branch)

		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Printf("Pipeline execution successfully! \nTriggered On: %s\nStatus: %s\n", cyan(execution.TriggeredOn), cyan(execution.Status))
		log.Printf("Executed By: %s\n", cyan(execution.Creator.Name))
		log.Printf("Checkout the execution at: %s", cyan(execution.HTMLURL))

		for {
			ok, err := checkStatus()
			if err != nil {
				log.Printf("Unable to check status: %v\n Checkout the pipeline: %s", err, execution.Pipeline.URL)
			}

			if ok {
				status, err := apiClient.CheckPipelineStatus(project, pipeline.ID, execution.ID)
				if err != nil {
					log.Printf("Error: %v", err)
					break
				}
				success := color.New(color.FgGreen).SprintFunc()
				inProgress := color.New(color.FgYellow).SprintFunc()
				failed := color.New(color.FgRed).SprintFunc()

				if *status == "SUCCESSFUL" {
					log.Printf("Current status: %s", success(*status))
					log.Println("Goodbye!")
					break
				} else if *status == "INPROGRESS" {
					log.Printf("Current status: %s", inProgress(*status))
					log.Printf("\nWaiting...")
					time.Sleep(7 * time.Second) // Adjust the sleep duration as needed
				} else if *status == "FAILED" {
					log.Printf("Current status: %s", failed(*status))
					log.Println("Goodbye!")
					break
				}
			} else {
				log.Println("Goodbye!")
				break
			}
		}
	},
}

func init() {
	// Add branch and pipeline flags
	deployCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Branch to deploy")
	deployCmd.Flags().StringVarP(&pipelineFlag, "pipeline", "p", "", "Pipeline to deploy (production or staging)")
	deployCmd.Flags().BoolVarP(&currentFlag, "current", "c", false, "Use the current Git branch for deployment")
	rootCmd.AddCommand(deployCmd)
}

// Function to search and select project interactively
func searchProject(projectsArray []buddy.Project) string {
	var projectNames []string

	// Loop over the []Project slice and extract the Name
	for _, project := range projectsArray {
		projectNames = append(projectNames, project.Name)
	}

	prompt := promptui.Select{
		Label: "Select a Project",
		Items: projectNames,
		Searcher: func(input string, index int) bool {
			project := projectNames[index]
			return containsIgnoreCase(project, input)
		},
		StartInSearchMode: true, // Start the prompt in search mode
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | bold }}",
			Active:   "▸ {{ . | cyan | bold }}",
			Inactive: "  {{ . | cyan }}",
			Selected: "✔  {{ . | cyan | bold }}",
		},
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return projectNames[i]
}

// Function to search and select branch interactively
func searchBranch(branchesArray []buddy.Branch) string {
	var branchNames []string

	// Loop over the []Branch slice and extract the Name
	for _, branch := range branchesArray {
		branchNames = append(branchNames, branch.Name)
	}

	prompt := promptui.Select{
		Label: "Select a Branch",
		Items: branchNames,
		Searcher: func(input string, index int) bool {
			branch := branchNames[index]
			return containsIgnoreCase(branch, input)
		},
		StartInSearchMode: true, // Start the prompt in search mode
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | bold }}",
			Active:   "▸ {{ . | green | bold }}",
			Inactive: "  {{ . | green }}",
			Selected: "✔  {{ . | green | bold }}",
		},
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return branchNames[i]
}

// Function to select pipeline interactively (production or staging)
func searchPipeline(pipelinesArray []buddy.Pipeline, _ string) buddy.Pipeline {
	var availablePipelines []string
	for _, pipelines := range pipelinesArray {
		availablePipelines = append(availablePipelines, pipelines.Name)
	}

	prompt := promptui.Select{
		Label: "Select Pipeline",
		Items: availablePipelines,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | bold }}",
			Active:   "▸ {{ . | magenta | bold }}",
			Inactive: "  {{ . | magenta }}",
			Selected: "✔  {{ . | magenta | bold }}",
		},
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	selectedPipeline := filterPipelineByName(pipelinesArray, availablePipelines[i])
	return *selectedPipeline
}

func filterPipelineByName(pipelines []buddy.Pipeline, name string) *buddy.Pipeline {
	for _, pipeline := range pipelines {
		if pipeline.Name == name {
			return &pipeline // Return the pipeline if the name matches
		}
	}
	return nil // Return nil if no match is found
}

func confirmDeployment() bool {
	prompt := promptui.Prompt{
		Label: "Are you sure you want to deploy (yes/no)",
		Validate: func(input string) error {
			if strings.ToLower(input) != "yes" && strings.ToLower(input) != "no" {
				return fmt.Errorf("please type 'yes' or 'no'")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return strings.ToLower(result) == "yes"
}

func checkStatus() (bool, error) {
	prompt := promptui.Prompt{
		Label: "Want to check the status (yes/no)",
		Validate: func(input string) error {
			if strings.ToLower(input) != "yes" && strings.ToLower(input) != "no" {
				return fmt.Errorf("please type 'yes' or 'no'")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return false, err
	}

	return strings.ToLower(result) == "yes", nil
}

// Helper function to do case-insensitive search
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
