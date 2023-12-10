export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

update-dependencies:
	go get -u ./...
	go mod tidy

install-tools:
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

lint:
	golangci-lint run ./... --timeout 300s

test:
	gotestsum --junitfile junit.xml.out -- -coverprofile=cover.out  ./...

semantic-release:
	semantic-release
