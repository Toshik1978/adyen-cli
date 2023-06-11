.PHONY: generate quality.lint quality.tests quality.tests.coverage app.dependencies.download app.build.darwin.amd64 app.build.darwin.arm64 app.build.linux.amd64 app.build.windows.amd64 app.build clean
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

BUILDSTAMP = $(shell date +%Y-%m-%d_%H:%M:%S)
COMMIT = $(shell git describe --tags --exact-match --match "v*.*.*" || git describe --match "v*.*.*" --tags || git describe --tags || git rev-parse --short HEAD)

app.build.darwin.amd64:
	@echo "+ $@"
	rm -rf bin/adyen-cli-darwin-amd64

	GOOS=darwin GOARCH=amd64 go build -v -ldflags "\
		-X main.Buildstamp=$(BUILDSTAMP) \
		-X main.Commit=$(COMMIT) \
	" -o bin/adyen-cli-darwin-amd64 cmd/main.go

app.build.darwin.arm64:
	@echo "+ $@"
	rm -rf bin/adyen-cli-darwin-arm64

	GOOS=darwin GOARCH=arm64 go build -v -ldflags "\
		-X main.Buildstamp=$(BUILDSTAMP) \
		-X main.Commit=$(COMMIT) \
	" -o bin/adyen-cli-darwin-arm64 cmd/main.go

app.build.linux.amd64:
	@echo "+ $@"
	rm -rf bin/adyen-cli-linux-amd64

	GOOS=linux GOARCH=amd64 go build -v -ldflags "\
		-X main.Buildstamp=$(BUILDSTAMP) \
		-X main.Commit=$(COMMIT) \
	" -o bin/adyen-cli-linux-amd64 cmd/main.go

app.build.windows.amd64:
	@echo "+ $@"
	rm -rf bin/adyen-cli-win64.exe

	GOOS=windows GOARCH=amd64 go build -v -ldflags "\
		-X main.Buildstamp=$(BUILDSTAMP) \
		-X main.Commit=$(COMMIT) \
	" -o bin/adyen-cli-win64.exe cmd/main.go

app.build: app.build.darwin.amd64

dist.build: app.build.darwin.amd64 app.build.darwin.arm64 app.build.linux.amd64 app.build.windows.amd64
	@echo "+ $@"
	cd bin && zip adyen-cli-darwin-amd64.zip adyen-cli-darwin-amd64
	cd bin && minisign -S -s ~/.minisign/adyen-cli.key -m adyen-cli-darwin-amd64.zip
	cd bin && zip adyen-cli-darwin-arm64.zip adyen-cli-darwin-arm64
	cd bin && minisign -S -s ~/.minisign/adyen-cli.key -m adyen-cli-darwin-arm64.zip
	cd bin && zip adyen-cli-linux-amd64.zip adyen-cli-linux-amd64
	cd bin && minisign -S -s ~/.minisign/adyen-cli.key -m adyen-cli-linux-amd64.zip
	cd bin && zip adyen-cli-win64.zip adyen-cli-win64.exe
	cd bin && minisign -S -s ~/.minisign/adyen-cli.key -m adyen-cli-win64.zip

clean:
	@echo "+ $@"
	go clean -testcache
	rm -rf bin/adyen-cli-darwin-amd64
	rm -rf bin/adyen-cli-darwin-arm64
	rm -rf bin/adyen-cli-linux-amd64
	rm -rf bin/adyen-cli-win64.exe
