.DEFAULT_GOAL := helper
GIT_COMMIT ?= $(shell git rev-parse --short=12 HEAD || echo "NoGit")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
TEXT_RED = \033[0;31m
TEXT_BLUE = \033[0;34;1m
TEXT_GREEN = \033[0;32;1m
TEXT_NOCOLOR = \033[0m
DOCKER_IMAGE_NAME = config-puller

helper: # Adapted from: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@echo "Available targets..." # @ will not output shell command part to stdout that Makefiles normally do but will execute and display the output.
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

command:
	echo "Application run: $(ENV_FILE_LOCATION)"
	docker run -it --rm -e AWS_SDK_LOAD_CONFIG=true -v ${HOME}/.aws:/root/.aws -v $(pwd)/output/:/output/ --env-file $(ENV_FILE_LOCATION) $(DOCKER_IMAGE_NAME)

.PHONY: test
test: ## Builds and then runs tests against the application
	go test -coverprofile ./coverage.html .

build:
	docker build -t $(DOCKER_IMAGE_NAME) .

dev: build ## Runs a dev version of the application
	$(MAKE) command ENV_FILE_LOCATION="./test/envfile.txt"

clean: ## Cleans up any old/unneeded items
