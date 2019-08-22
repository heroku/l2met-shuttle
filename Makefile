GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GO_PACKAGES := $(shell go list ./... | grep -v /vendor/)

travis: tidy test

test:
	go test -mod=vendor -race -v $(GO_PACKAGES)

# Setup & Code Cleanliness
setup: hooks tidy

hooks:
	ln -fs ../../bin/git-pre-commit.sh .git/hooks/pre-commit

tidy: goimports
	test -z "$$(goimports -l -d $(GO_FILES) | tee /dev/stderr)"
	go vet -mod=vendor ./...
	go mod vendor
	go mod tidy
	go mod verify

precommit: tidy test

goimports:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
