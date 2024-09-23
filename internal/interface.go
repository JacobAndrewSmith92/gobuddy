// buddy/interface.go

package buddy

// BuddyAPI defines the interface for interacting with Buddy
type BuddyAPI interface {
	FetchProjects() ([]Project, error)
	FetchBranches(project string) ([]Branch, error)
	FetchPipelines(project string) ([]Pipeline, error)
	FetchProjectByName(name string) (*Project, error)
	FetchBranchByName(project, name string) (*Branch, error)
	FetchPipelineByID(project string, id int) (*Pipeline, error)
	RunPipeline(project string, pipelineID int, branch string) (*PipelineExecutionResponse, error)
	CheckPipelineStatus(project string, pipeline int, executionID int) (*string, error)
}

type ProjectResponse struct {
	URL      string    `json:"url"`
	HTMLURL  string    `json:"html_url"`
	Projects []Project `json:"projects"` // This contains the array of Project objects
}

// The Project struct
type Project struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

// BranchResponse struct to capture the full response
type BranchResponse struct {
	URL      string   `json:"url"`
	HTMLURL  string   `json:"html_url"`
	Branches []Branch `json:"branches"` // This contains the array of Branch objects
}

// Branch struct to represent an individual Git branch
type Branch struct {
	URL     string `json:"url,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
	Name    string `json:"name,omitempty"`
	Default bool   `json:"default,omitempty"`
}

type PipelineResponse struct {
	URL       string     `json:"url"`
	HTMLURL   string     `json:"html_url"`
	Pipelines []Pipeline `json:"pipelines"`
}

// Pipeline struct to represent an individual Git branch
type Pipeline struct {
	URL      string   `json:"url"`
	HTMLURL  string   `json:"html_url"`
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Priority string   `json:"priority,omitempty"`
	Refs     []string `json:"refs,omitempty"`
}

// Committer represents the committer object
type Committer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
	Admin bool   `json:"admin,omitempty"`
}

type Author struct {
	Email string `json:"email"`
}

// Revision represents the to_revision and from_revision object
type Revision struct {
	URL        string    `json:"url,omitempty"`
	HTMLURL    string    `json:"html_url,omitempty"`
	Revision   string    `json:"revision,omitempty"`
	AuthorDate string    `json:"author_date,omitempty"`
	CommitDate string    `json:"commit_date,omitempty"`
	Message    string    `json:"message,omitempty"`
	Committer  Committer `json:"committer,omitempty"`
	Author     Author    `json:"author,omitempty"`
}

// PipelineExecutionRequest represents the payload to trigger the pipeline execution
type PipelineExecutionRequest struct {
	ToRevision Revision `json:"to_revision"`
	Branch     Branch   `json:"branch"`
}

// Creator struct for the user who triggered the pipeline
type Creator struct {
	URL            string `json:"url"`
	HTMLURL        string `json:"html_url"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	Admin          bool   `json:"admin"`
	WorkspaceOwner bool   `json:"workspace_owner"`
}

// PipelineExecutionResponse struct for the full response of a pipeline execution
type PipelineExecutionResponse struct {
	URL          string   `json:"url"`
	HTMLURL      string   `json:"html_url"`
	ID           int      `json:"id"`
	StartDate    string   `json:"start_date"`
	FinishDate   *string  `json:"finish_date"`
	TriggeredOn  string   `json:"triggered_on"`
	Priority     string   `json:"priority"`
	Refresh      bool     `json:"refresh"`
	ClearCache   bool     `json:"clear_cache"`
	Status       string   `json:"status"`
	Comment      string   `json:"comment"`
	Branch       Branch   `json:"branch"`
	FromRevision Revision `json:"from_revision"`
	ToRevision   Revision `json:"to_revision"`
	Creator      Creator  `json:"creator"`
	Pipeline     Pipeline `json:"pipeline"`
	// ActionExecutions []ActionExecution `json:"action_executions"`
}

type ErrorDetail struct {
	Message string `json:"message,omitempty"`
}

// ErrorResponse represents the full error response structure
type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}
