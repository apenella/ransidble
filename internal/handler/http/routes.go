package http

// Route constants for API endpoints
const (
	// ProjectBasePath is the base path for all project-related endpoints
	ProjectBasePath = "/projects"
	// CreateProjectPath is the endpoint to create a new project
	CreateProjectPath = "/projects/:id"
	// GetProjectPath is the endpoint to get a project by ID
	GetProjectPath = "/projects/:id"
	// GetProjectsPath is the endpoint to list all projects
	GetProjectsPath = "/projects"
	// DeleteProjectPath is the endpoint to delete a project by ID
	DeleteProjectPath = "/projects/:id"

	// TaskBasePath is the base path for all task-related endpoints
	TaskBasePath = "/tasks"
	// CreateTaskAnsiblePlaybookPath is the endpoint to create a new Ansible playbook task
	CreateTaskAnsiblePlaybookPath = "/tasks/ansible-playbook/:project_id"
	// GetTaskPath is the endpoint to get a task by ID
	GetTaskPath = "/tasks/:id"
	// GetTasksPath is the endpoint to list all tasks
	GetTasksPath = "/tasks"

	// GetHealthPath is the endpoint to check the health of the service
	GetHealthPath = "/health"
)
