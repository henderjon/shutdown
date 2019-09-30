export CC_TEST_REPORTER_ID = a9f420a89d2147cbd748cff6efce62c046672ddf83c0df67353a0f6e82d2ce2d
COVERAGEOUTFILE=c.out

all: test race

.PHONY: dep
dep:
	go mod vendor

.PHONY: test-vendor
test-vendor:
	go test -mod=vendor -coverprofile=coverage.out -covermode=count

.PHONY: test
test: dep
	go test -coverprofile=coverage.out -covermode=count

.PHONY: race
race: dep
	go test -race

.PHONY: test-report
test-report: test
	go tool cover -html=coverage.out

.PHONY: travis
travis:
	# install deps pre-1.13
	go get -u github.com/google/go-cmp/cmp
	go test -coverprofile $(COVERAGEOUTFILE) ./...

.PHONY: cclimate-linux
cclimate-linux:
	curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
	# curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-darwin-amd64 > ./cc-test-reporter
	chmod +x ./cc-test-reporter
	./cc-test-reporter before-build
	# install deps pre-1.13
	go get -u github.com/google/go-cmp/cmp
	go test -coverprofile $(COVERAGEOUTFILE) ./...
	./cc-test-reporter after-build --exit-code $(TRAVIS_TEST_RESULT)
