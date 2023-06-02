BASE_PATH := $(shell pwd)
CMD_PATH := "cmd"
CMD_DIRS := $(shell find $(CMD_PATH)/* -type d)
GO_COMPILE:=GOOS=linux GOARCH=amd64 go build 

.SILENT:
init:
	echo "ℹ️  INITIALIZING PROJECT..."
	npm install
	go mod tidy

ifdef profile
	cdk bootstrap --profile ${profile}
else
	cdk bootstrap
endif

.PHONY: clean compile_all
clean:
	echo "ℹ️  CLEANING ALL BUILD FILES..."
	for base in $(CMD_DIRS); do \
		dirname=$$(basename $$base); \
		cd $(BASE_PATH)/$$base && \
		if [ -f $(BASE_PATH)/$$base/$$dirname ]; then \
			echo "- $$dirname"; \
			rm $$dirname; \
		fi; \
	done
	echo "\n"

compile_all: clean
	echo "ℹ️  STARTING TO COMPILE ALL..."
	for base in $(CMD_DIRS); do \
		dirname=$$(basename $$base); \
		cd $(BASE_PATH)/$$base && \
		if [ -s main.go ]; then \
			echo "- $$dirname"; \
			cd $(BASE_PATH)/$$base && \
			$(GO_COMPILE) -o $$dirname main.go; \
		fi; \
	done
	echo "\n"

deploy: compile_all
	echo "🚀 Deploying stack..."

ifdef profile
	cdk deploy --profile ${profile}
else
	cdk deploy
endif