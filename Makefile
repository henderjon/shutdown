BIN=vasili
HEAD=$(shell ([ -n "$${CI_TAG}" ] && echo "$$CI_TAG" || exit 1) || git describe --dirty --long --tags 2> /dev/null || git rev-parse --short HEAD)
TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S %z %Z')
DEPLOYMENT_PATH=s3://myon-deployment/pre-release/$(BIN)/$(BIN)-$(HEAD)

LDFLAGS="-X 'main.buildVersion=$(HEAD)' -X 'main.buildTimestamp=$(TIMESTAMP)' -X 'main.compiledBy=$(shell go version)'" # `-s -w` removes some debugging info that might not be necessary in production (smaller binaries)

all: print

.PHONY: dep
dep:
	go mod vendor

.PHONY: test-vendor
test-vendor:
	go test -mod=vendor -coverprofile=coverage.out -covermode=count

.PHONY: test
test:
	go test -coverprofile=coverage.out -covermode=count

.PHONY: race
race:
	go test -race

.PHONY: test-report
test-report: test
	go tool cover -html=coverage.out
