all: test race

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
