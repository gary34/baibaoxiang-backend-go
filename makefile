GitBranch=$(shell cat .git/HEAD | awk -F '/' '{print $$3}')
GitCommit=$(shell cat .git/refs/heads/${GitBranch})
GO_CFLAGS=-ldflags "-X main.BuildDate=`date -u +%Y-%m-%d,%H:%M:%S` \
					-X main.GitBranch=$(GitBranch) \
					-X main.GitCommit=$(GitCommit)"

build:
		@echo "\t > compiling golang sources..."
		@export GOPATH=$(shell pwd);\
				go build -o bin/server $(GO_CFLAGS) main
		@echo "\t > done <"