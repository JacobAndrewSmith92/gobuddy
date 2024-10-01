/*
Package buddy provides a client for interacting with the Buddy API.
The available methods allow fetching projects, branches, pipelines, and executing pipelines.
*/
package buddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BuddyClient represents the actual Buddy API client
type BuddyClient struct {
	Token     string
	Workspace string
}

// NewBuddyClient initializes a new BuddyClient with token and workspace from the config
func NewBuddyClient(token, workspace string) *BuddyClient {
	return &BuddyClient{
		Token:     token,
		Workspace: workspace,
	}
}

// FetchProjects fetches projects from the Buddy API
func (c *BuddyClient) FetchProjects() ([]Project, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects?per_page=100", c.Workspace)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching projects: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projectResponse ProjectResponse
	err = json.Unmarshal(body, &projectResponse)
	if err != nil {
		return nil, err
	}

	return projectResponse.Projects, nil
}

// FetchProjectByName fetches a project by name from the Buddy API
func (c *BuddyClient) FetchProjectByName(name string) (*Project, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s", c.Workspace, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project %s not found - status: %s", name, resp.Status)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching project: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var project Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// FetchBranches fetches branches for a specific project
func (c *BuddyClient) FetchBranches(project string) ([]Branch, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/repository/branches", c.Workspace, project)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching branches: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var branchResponse BranchResponse
	err = json.Unmarshal(body, &branchResponse)
	if err != nil {
		return nil, err
	}

	return branchResponse.Branches, nil
}

// FetchBranchByName fetches a branch by name for a given project
func (c *BuddyClient) FetchBranchByName(project, branch string) (*Branch, *ErrorResponse) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/repository/branches/%s", c.Workspace, project, branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &ErrorResponse{
			Errors: []ErrorDetail{
				{
					Message: fmt.Sprintf("Failed to build request to get branch: %v", err),
				},
			},
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, &ErrorResponse{
			Errors: []ErrorDetail{
				{
					Message: fmt.Sprintf("Failed to fetch branch: %v", err),
				},
			},
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &ErrorResponse{
			Errors: []ErrorDetail{
				{
					Message: fmt.Sprintf("Branch %s not found in project %s", branch, project),
				},
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ErrorResponse{
			Errors: []ErrorDetail{
				{
					Message: fmt.Sprintf("Failed to read response body: %v", err),
				},
			},
		}
	}

	var branchResponse Branch
	err = json.Unmarshal(body, &branchResponse)
	if err != nil {
		return nil, &ErrorResponse{
			Errors: []ErrorDetail{
				{
					Message: fmt.Sprintf("Failed to unmarshal branch response: %v", err),
				},
			},
		}
	}

	return &branchResponse, nil
}

// FetchPipelines fetches pipelines for a specific project
func (c *BuddyClient) FetchPipelines(project string) ([]Pipeline, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/pipelines", c.Workspace, project)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching pipelines: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pipelineResponse PipelineResponse
	err = json.Unmarshal(body, &pipelineResponse)
	if err != nil {
		return nil, err
	}

	return pipelineResponse.Pipelines, nil
}

// FetchPipelineByID fetches a pipeline for a specific project by ID
// - can be used if dev knows the pipeline ID or I need to do some more logic to map name to ID
func (c *BuddyClient) FetchPipelineByID(project, id string) (*Pipeline, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/pipelines/%s", c.Workspace, project, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching pipelines: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pipelineResponse Pipeline
	err = json.Unmarshal(body, &pipelineResponse)
	if err != nil {
		return nil, err
	}

	return &pipelineResponse, nil
}

// RunPipeline triggers the execution of a pipeline
func (c *BuddyClient) RunPipeline(project string, pipelineID int, branch string) (*PipelineExecutionResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/pipelines/%d/executions", c.Workspace, project, pipelineID)

	requestBody := PipelineExecutionRequest{
		ToRevision: Revision{
			Revision: "HEAD",
		},
		Branch: Branch{
			Name: branch,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Response Body: %s\n", resp.Body)
		return nil, fmt.Errorf("error executing pipeline: %s", resp.Status)
	}

	var executionResponse PipelineExecutionResponse
	err = json.NewDecoder(resp.Body).Decode(&executionResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &executionResponse, nil
}

// CheckPipelineStatus fetches the status of a pipeline execution
func (c *BuddyClient) CheckPipelineStatus(project string, pipeline int, executionID int) (*string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/pipelines/%d/executions/%d", c.Workspace, project, pipeline, executionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching pipeline status: %s", resp.Status)
	}

	var executionResponse PipelineExecutionResponse
	err = json.NewDecoder(resp.Body).Decode(&executionResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &executionResponse.Status, nil
}
