DOCKER_COMPOSE_BINARY := $(shell docker compose version > /dev/null 2>&1 && echo "docker compose" || (which docker-compose > /dev/null 2>&1 && echo "docker-compose" || (echo "docker compose not found. Aborting." >&2; exit 1)))

## Colors
COLOR_GREEN=\033[0;32m
COLOR_RED=\033[0;31m
COLOR_BLUE=\033[0;34m
COLOR_PURPLE=\033[0;35m
COLOR_END=\033[0m

.DEFAULT_GOAL := help

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Lists available targets
	@echo
	@echo "Makefile usage:"
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[1;32m%-25s\033[0m %s\n", $$1, $$2}' | sort
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

tests: unit-tests ## Executes tests

unit-tests: ## Executes unit test
	@echo
	@echo "$(COLOR_BLUE) Executing unit test$(COLOR_END)"
	@echo
	@docker run --rm -v "${PWD}":/app -w /app golang:${GOLANG_VERSION}-alpine go test -count=1 -cover ./internal/... && echo "$(COLOR_GREEN) Unit test: OK$(COLOR_END)" || echo "$(COLOR_RED)Unit test: some test failed$(COLOR_END)"
#
# Environment targets

serve: ## Start the server
	@echo
	@echo " Starting the server"
	@echo
	@RANSIDBLE_SERVER_LOG_LEVEL=debug RANSIDBLE_SERVER_WORKER_POOL_SIZE=1 RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH=test/projects  go run cmd/main.go serve

run-task-1: ## Make a request to create an ansible-playbook task
	@echo
	@echo " Making a request to the server"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-1 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

run-task-2: ## Make a request to create an ansible-playbook task
	@echo
	@echo " Making a request to the server"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-2 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local", "dependencies": {"collections": {"requirements_file": "requirements.yml"}}}'

run-task-3: ## Make a request to create an ansible-playbook task
	@echo
	@echo " Making a request to the server"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-3 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

run-task-4: ## Make a request to create an ansible-playbook task
	@echo
	@echo " Making a request to the server"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-4 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

get-task: ## Get the task status
	@echo
	@echo " Getting the task status"
	@echo
	curl -XGET 0.0.0.0:8080/tasks/$(TASK_ID)

validate-openapi: ## Check the openapi spec
	@echo
	@echo " Checking the openapi spec"
	@echo
	@docker run --rm -v "${PWD}/api/openapi.yaml":/openapi.yaml jeanberu/swagger-cli swagger-cli validate /openapi.yaml
