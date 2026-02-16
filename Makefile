DOCKER_COMPOSE_BINARY := $(shell docker compose version > /dev/null 2>&1 && echo "docker compose" || (which docker-compose > /dev/null 2>&1 && echo "docker-compose" || (echo "docker compose not found. Aborting." >&2; exit 1)))

## Colors
COLOR_GREEN=\033[0;32m
COLOR_RED=\033[0;31m
COLOR_BLUE=\033[0;34m
COLOR_PURPLE=\033[0;35m
COLOR_END=\033[0m

## Serve target configuration

### SERVECONF_RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS: Address to listen for HTTP requests
SERVECONF_RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS ?= :8081
### SERVECONF_DOCKER_RUN_PUBLISH_PORT: Port to publish the server
SERVECONF_DOCKER_RUN_PUBLISH_PORT ?= 8080
### SERVECONF_RANSIDBLE_SERVER_LOG_LEVEL: Log level for the server
SERVECONF_RANSIDBLE_SERVER_LOG_LEVEL ?= debug
### SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_LOCAL_PATH: Path to store the project files when using local repository type
SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_LOCAL_PATH ?= /repository/projects
### SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_TYPE: Type of project repository to use (local or git)
SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_TYPE ?= local
### SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_LOCAL_PATH: Path to store the project files when using local storage type
SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_LOCAL_PATH ?= /storage/projects
### SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_TYPE: Type of project storage to use (local or s3)
SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_TYPE ?= local
### SERVECONF_RANSIDBLE_SERVER_WORKER_POOL_SIZE: Number of workers to process the tasks
SERVECONF_RANSIDBLE_SERVER_WORKER_POOL_SIZE ?= 1

.DEFAULT_GOAL := help

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Lists available targets
	@echo
	@echo "Makefile usage:"
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$'  $(filter-out .env, $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[1;32m%-25s\033[0m %s\n", $$1, $$2}' | sort
	@echo

#
# Development lifecycle targets

static-analysis: vet golint ## Execute static analysis

ci-go-tools-docker-image: ## Build the docker image
	@echo
	@echo "$(COLOR_BLUE) Building the docker image $(COLOR_END)"
	@echo
	@docker build -t ci-go-tools-docker-image -f build/Dockerfile.ci .

vet: ci-go-tools-docker-image ## Executes the go vet to report any suspicious constructs
	@echo
	@echo "$(COLOR_BLUE) Executing go vet $(COLOR_END)"
	@echo
	@docker run --rm -v "${PWD}":/app -w /app ci-go-tools-docker-image go vet ./internal/... && echo "$(COLOR_GREEN) go vet: all files linted$(COLOR_END)" || echo "$(COLOR_RED)go vet: some files not linted$(COLOR_END)"

golint: ci-go-tools-docker-image ## Executes Go linter (golint)
	@echo
	@echo "$(COLOR_BLUE) Executing golint$(COLOR_END)"
	@echo
	@docker run --rm -v "${PWD}":/app -w /app ci-go-tools-docker-image golint ./internal/... && echo "$(COLOR_GREEN) golint: all files linted$(COLOR_END)" || echo "$(COLOR_RED)golint: some files not linted$(COLOR_END)"

tests: unit-tests functional-test validate-openapi ## Executes tests

unit-tests: ## Executes unit test
	@echo
	@echo "$(COLOR_BLUE) Executing unit test$(COLOR_END)"
	@echo
	@docker run --rm -v "${PWD}":/app -w /app golang:${GOLANG_VERSION}-alpine go test -count=1 -cover ./internal/... && echo "$(COLOR_GREEN) Unit test: OK$(COLOR_END)" || echo "$(COLOR_RED)Unit test: some test failed$(COLOR_END)"

functional-test: ## Execute functional test
	@echo
	@echo "$(COLOR_BLUE) Executing functional test$(COLOR_END)"
	@echo
	@docker run --rm -v "${PWD}":/app -w /app golang:${GOLANG_VERSION}-alpine go test -count=1 -cover ./test/functional/... && echo "$(COLOR_GREEN) Functional test: OK$(COLOR_END)" || echo "$(COLOR_RED)Functional test: some test failed$(COLOR_END)"

validate-openapi: ## Check the openapi spec
	@echo
	@echo "$(COLOR_BLUE) Checking the openapi spec$(COLOR_END)"
	@echo
	@docker run --rm -v "${PWD}/api/openapi.yaml":/openapi.yaml jeanberu/swagger-cli swagger-cli validate /openapi.yaml

#
# Environment targets

serve: ## Start a Ransidble server
	@echo
	@echo "$(COLOR_BLUE) Starting a Ransidble server$(COLOR_END)"
	@echo
	@echo "$(COLOR_BLUE)  Docker run publish port: $(COLOR_GREEN)$(SERVECONF_DOCKER_RUN_PUBLISH_PORT)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Log level:%-25s$(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_LOG_LEVEL)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Project local storage path: $(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Worker pool size: $(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_WORKER_POOL_SIZE)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Project repository type: $(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_TYPE)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Project repository local path: $(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_LOCAL_PATH)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Project storage type: $(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_TYPE)$(COLOR_END)"
	@echo "$(COLOR_BLUE)  Project storage local path: $(COLOR_GREEN)$(SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_LOCAL_PATH)$(COLOR_END)"
	@echo
	@$(DOCKER_COMPOSE_BINARY) run \
		--build \
		--env RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS=$(SERVECONF_RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS) \
		--env RANSIDBLE_SERVER_LOG_LEVEL=$(SERVECONF_RANSIDBLE_SERVER_LOG_LEVEL) \
		--env RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH=$(SERVECONF_RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH) \
		--env RANSIDBLE_SERVER_PROJECT_REPOSITORY_LOCAL_PATH=$(SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_LOCAL_PATH) \
		--env RANSIDBLE_SERVER_PROJECT_REPOSITORY_TYPE=$(SERVECONF_RANSIDBLE_SERVER_PROJECT_REPOSITORY_TYPE) \
		--env RANSIDBLE_SERVER_PROJECT_STORAGE_LOCAL_PATH=$(SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_LOCAL_PATH) \
		--env RANSIDBLE_SERVER_PROJECT_STORAGE_TYPE=$(SERVECONF_RANSIDBLE_SERVER_PROJECT_STORAGE_TYPE) \
		--env RANSIDBLE_SERVER_WORKER_POOL_SIZE=$(SERVECONF_RANSIDBLE_SERVER_WORKER_POOL_SIZE) \
		--interactive \
		--publish $(SERVECONF_DOCKER_RUN_PUBLISH_PORT):$(shell echo $${SERVECONF_RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS##*:}) \
		--workdir /usr/src/app \
		ransidble-server go run cmd/main.go serve

create-projects: create-project-1 create-project-2 create-project-3 create-project-4 ## Create projects

create-project-1: ## Create a project with the name project-1
	@echo
	@echo " $(COLOR_BLUE)Creating project project-1$(COLOR_END)"
	@echo
	curl -iX POST 0.0.0.0:8080/projects -H 'Content-Type: multipart/form-data' -F 'metadata={"format":"targz","storage":"local"};type=application/json' -F 'file=@test/fixtures/projects/project-1.tar.gz'

create-project-2: ## Create a project with the name project-2
	@echo
	@echo " $(COLOR_BLUE)Creating project project-2$(COLOR_END)"
	@echo
	curl -iX POST 0.0.0.0:8080/projects -H 'Content-Type: multipart/form-data' -F 'metadata={"format":"targz","storage":"local"};type=application/json' -F 'file=@test/fixtures/projects/project-2.tar.gz'

create-project-3: ## Create a project with the name project-3
	@echo
	@echo " $(COLOR_BLUE)Creating project project-3$(COLOR_END)"
	@echo
	curl -iX POST 0.0.0.0:8080/projects -H 'Content-Type: multipart/form-data' -F 'metadata={"format":"targz","storage":"local"};type=application/json' -F 'file=@test/fixtures/projects/project-3.tar.gz'

create-project-4: ## Create a project with the name project-4
	@echo
	@echo " $(COLOR_BLUE)Creating project project-4$(COLOR_END)"
	@echo
	curl -iX POST 0.0.0.0:8080/projects -H 'Content-Type: multipart/form-data' -F 'metadata={"format":"targz","storage":"local"};type=application/json' -F 'file=@test/fixtures/projects/project-4.tar.gz'

list-projects: ## List the projects
	@echo
	@echo " $(COLOR_BLUE)Listing the projects$(COLOR_END)"
	@echo
	curl -s -XGET 0.0.0.0:8080/projects | jq

run-task-1: ## Make a request to create an ansible-playbook task
	@echo
	@echo " $(COLOR_BLUE)Making a request to the server$(COLOR_END)"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-1 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

run-task-2: ## Make a request to create an ansible-playbook task
	@echo
	@echo " $(COLOR_BLUE)Making a request to the server$(COLOR_END)"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-2 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local", "dependencies": {"collections": {"requirements_file": "requirements.yml", "force_with_deps": true}}}'

run-task-3: ## Make a request to create an ansible-playbook task
	@echo
	@echo " $(COLOR_BLUE)Making a request to the server$(COLOR_END)"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-3 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

run-task-4: ## Make a request to create an ansible-playbook task
	@echo
	@echo " $(COLOR_BLUE)Making a request to the server$(COLOR_END)"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-4 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

get-task: ## Get the task status
	@echo
	@echo " $(COLOR_BLUE)Getting the task status$(COLOR_END)"
	@echo
	curl -XGET 0.0.0.0:8080/tasks/$(TASK_ID)
