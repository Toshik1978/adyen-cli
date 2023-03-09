.PHONY: generate quality.lint quality.tests quality.tests.coverage app.dependencies.download app.build clean
.DEFAULT_GOAL := all

all: quality.lint quality.tests app.build

generate:
	@echo "+ $@"
	GO111MODULE=on go generate ./...

quality.lint:
	@echo "+ $@"
	./scripts/quality.lint.sh

quality.tests:
	@echo "+ $@"
	GO111MODULE=on go test -v ./...

quality.tests.coverage:
	@echo "+ $@"
	GO111MODULE=on go test -race -coverprofile=coverage.txt -covermode=atomic ./...

app.dependencies.download:
	@echo "+ $@"
	GO111MODULE=on go mod download -x

app.build:
	@echo "+ $@"
	rm -rf bin/csv2adyen
	GO111MODULE=on GOGC=off go build -v -ldflags "\
		-X main.Buildstamp=$(shell date +%Y/%m/%d_%H:%M:%S) \
		-X main.Commit=$(shell git rev-parse --short HEAD) \
	" -o bin/csv2adyen cmd/main.go

clean:
	@echo "+ $@"
	go clean -testcache
	rm -rf bin/csv2adyen
