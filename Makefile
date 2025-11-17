SHELL := /bin/bash

MODULE := $$(go list -m)

.PHONY: default
default: tidy fmt generate build

.PHONY: tools
tools: tidy
	@go -C tools install tool

.PHONY: update
update:
	@go -C tools get tool
	@go get $$(go mod edit -json | jq -r '.Require[] | select(.Indirect | not) | .Path')
	@buf dep prune
	@buf dep update
	@$(MAKE) tidy
	@$(MAKE) tools
	@$(MAKE) build

.PHONY: tidy
tidy:
	@go -C tools mod tidy
	@go -C tools mod download
	@go mod tidy
	@go mod download

.PHONY: fmt
fmt: fmt/go fmt/proto

.PHONY: fmt/go
fmt/go:
	@golangci-lint fmt ./...

.PHONY: fmt/proto
fmt/proto:
	@buf format --write .

.PHONY: generate
generate: generate/go generate/proto generate/docs

.PHONY: generate/go
generate/go:
	@go generate ./...

.PHONY: generate/proto
generate/proto:
	@rm -rf pkg/proto
	@buf generate

.PHONY: generate/docs
generate/docs:
	@find pkg/docs/assets/openapi -type f -name '*.openapi.json' -delete
	@rm -rf pkg/docs/assets/swagger
	@find pkg/proto/blob/api -type f -name '*.openapi.json'  | while read -r file; do \
		mkdir -p $$(dirname "pkg/docs/assets/openapi/$${file#pkg/proto/}"); \
		cp "$$file" "pkg/docs/assets/openapi/$${file#pkg/proto/}"; \
	done
	@go run ./cmd/internal/swagger-gen -c swagger-gen.yaml

.PHONY: lint
lint: lint/go lint/proto

.PHONY: lint/go
lint/go:
	@govulncheck ./...
	@golangci-lint run ./...

.PHONY: lint/proto
lint/proto:
	@buf lint
	@buf breaking --against "https://$(MODULE).git#branch=main"

.PHONY: test
test:
	@go test -coverprofile=cover.out -race ./...

.PHONY: cover/html
cover/html: test
	@go tool cover -html=cover.out

.PHONY: cover/func
cover/func: test
	@go tool cover -func=cover.out

.PHONY: pprof/http
pprof/http:
	@go tool pprof -http=localhost:8081 http://localhost:8080/debug/pprof/profile
	@open http://localhost:8081

.PHONY: build
build:
	@echo "building blob"
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -extldflags='-static'" -o ./bin/blob ./cmd/blob

.PHONY: clean
clean:
	@rm -rf bin
