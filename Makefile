#.PHONY: help

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

build: ## build and up docker containers
	@docker-compose up --build -d

start: ## run cart server
	@docker-compose up -d

test: ## run tests
	#@go test ./... -v
	@go DEVOM_API_URL=http://localhost:8030/api/v1
	#custom pkg tesst
	#go test ./pkg/server -run TestParseDevotionals -v
	#go test ./pkg/fs -run TestDocFeeder -v

coverage:
	mkdir -p .build/test_results
	@go test -coverprofile=.build/test_results/coverage.out ./...
	@go tool cover -func=.build/test_results/coverage.out

