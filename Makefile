.DEFAULT_GOAL := help

help: ## Lists available targets
	@echo
	@echo "Makefile usage:"
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[1;32m%-20s\033[0m %s\n", $$1, $$2}' | sort
	@echo

serve: ## Start the server
	@echo
	@echo " Starting the server"
	@echo
	@RANSIDBLE_SERVER_LOG_LEVEL=debug RANSIDBLE_SERVER_WORKER_POOL_SIZE=1 RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH=test/projects  go run cmd/main.go serve

run-task-1: ## Make a request to create an ansible-playbook task
	@echo
	@echo " Making a request to the server"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/task/ansible-playbook/project-1 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'

run-task-2: ## Make a request to create an ansible-playbook task
	@echo
	@echo " Making a request to the server"
	@echo
	curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/task/ansible-playbook/project-2 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local", "dependencies": {"collections": {"requirements_file": "requirements.yml"}}}'


get-task: ## Get the task status
	@echo
	@echo " Getting the task status"
	@echo
	curl -XGET 0.0.0.0:8080/task/$(TASK_ID)
