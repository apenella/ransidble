# Ransidble Roadmap

## Enhancements
- Store the execution output and make it available to the user
- Add pagination to the get projects and tasks
- Do not accept local connection
- At this moment, the project repository loads the projects from the local storage. In the future, the plan is to have a database where you need to create a project before running it. (internal/handler/cli/serve: serve.go)
- The following cases are not supported yet in tar extraction, it would be required when fetching files from a git source (internal/infrastructure/tar: tar.go)
- Remove the ability to execute playbooks directly on the server (currently temporarily enabled) (internal/infrastructure/executor: ansiblePlaybook.go)
- Return an error if there are no playbooks to run before calling the executor (internal/infrastructure/executor: ansiblePlaybook.go)
- workingDir should be set to a directory when the project is a directory and a file otherwise (internal/domain/core/service/workspace: workspace.go)
- id generator should be injected as a dependency (for task creation) (internal/domain/core/service/task: createTaskAnsiblePlaybookService.go)

## Features
- Create a use case for loading projects from the file system
- Implement a client to interact with the API
- Upload projects in multiple formats
  - tar.gz (multiple: see internal/infrastructure/unpack/tarGzipFormat.go, internal/domain/core/entity/project.go, etc.)
  - oci image
- Implement authentication (delegated to an external service)
- Implement authorization (delegated to an external service)

## Ideas
- RolesPath: support specifying the path where roles should be installed on the local filesystem (internal/infrastructure/executor: ansiblePlaybook.go, internal/domain/core/model/request/ansiblePlaybookParameters.go, internal/domain/core/entity/ansiblePlaybookParameters.go)
