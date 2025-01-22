SHELL := /bin/bash

MODULE := $$(go list -m)
VERSION := 1.0.0
SWAGGER_UI_VERSION :=

.PHONY: default
default: tidy fmt generate install

.PHONY: tools
tools: tidy
	@rm -f bin/*; \
	go -C internal/tools list -e -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go -C internal/tools install

.PHONY: update
update:
	@go list -m -f '{{if and (not .Main) (not .Indirect)}}{{.Path}}{{end}}' all | xargs go get; \
	go -C internal/tools list -m -f '{{if and (not .Main) (not .Indirect)}}{{.Path}}{{end}}' all | xargs go -C internal/tools get; \
	$(MAKE) tidy

.PHONY: tidy
tidy:
	@go mod tidy; \
	go -C internal/tools mod tidy

.PHONY: fmt
fmt: fmt/buf
	@go fmt ./...; \
	go -C internal/tools fmt ./... 2>&1 | grep -v 'is a program, not an importable package' || true; \
	goimports -w -local $(MODULE) .; \
	goimports -w -local $(MODULE) internal/tools; \
	tagalign -fix -sort -order "json,yaml,validate" --strict ./... 2>&1 | grep -v 'proto' || true; \

.PHONY: fmt/buf
fmt/buf:
	@buf format --write .

.PHONY: generate
generate: generate/buf generate/swagger
	@go generate ./...

.PHONY: generate/buf
generate/buf:
	@rm -rf pkg/gen; \
	find swagger/dist -type f -name '*.swagger.json' -delete; \
	buf generate

.PHONY: generate/swagger
generate/swagger:
	@version=$(SWAGGER_UI_VERSION); \
	if [[ -z "$${version}" ]]; then \
		version="$$(curl -sSL https://api.github.com/repos/swagger-api/swagger-ui/releases/latest | jq -r '.tag_name' | sed 's/^v//')"; \
	fi; \
	rm -rf /tmp/swagger-ui.tar.gz; \
	curl -sSLo /tmp/swagger-ui.tar.gz "https://github.com/swagger-api/swagger-ui/archive/refs/tags/v$${version}.tar.gz"; \
	rm -rf /tmp/swagger-ui; \
	mkdir -p /tmp/swagger-ui; \
	tar -xzf /tmp/swagger-ui.tar.gz -C /tmp/swagger-ui; \
	mkdir -p swagger/dist; \
	find swagger/dist -type f -not -name '*.swagger.json' -delete; \
	cp -r /tmp/swagger-ui/swagger-ui-$${version}/dist/ swagger/dist/; \
	urls="    urls: ["; \
	for file in "$$(find swagger/dist -type f -name "*.swagger.json")"; do \
		path="$${file#swagger/dist/}"; \
		urls+="\n      { name: \"$${path}\", url: \"$${path}\" },\n"; \
	done; \
	urls+="    ],"; \
	line="$$(cat swagger/dist/swagger-initializer.js | grep -n "url" | cut -d: -f1)"; \
	before="$$(head -n "$$(($${line} - 1))" swagger/dist/swagger-initializer.js)"; \
	after="$$(tail -n +"$$(($${line} + 1))" swagger/dist/swagger-initializer.js)"; \
	echo -e "$${before}\n$${urls}\n$${after}" >swagger/dist/swagger-initializer.js

.PHONY: lint
lint:
	@go vet ./...; \
	golangci-lint run ./...; \
	govulncheck ./...; \
	buf lint; \
	buf breaking --against "https://$(MODULE).git#branch=main"

.PHONY: test
test:
	@go test -v ./...

.PHONY: build
build:
	@$(MAKE) binary cmd=build version=$(VERSION)

.PHONY: install
install:
	@$(MAKE) binary cmd=install version=$(VERSION)

.PRONY: binary
binary: generate
	@if [[ -z "$${cmd}" ]]; then \
		echo "must set cmd env var"; \
		exit 1; \
	fi; \
	if [[ "$${cmd}" != "build" && "$${cmd}" != "install" ]]; then \
		echo "unknown cmd '$${cmd}'"; \
		exit 1; \
	fi; \
	if [[ -z "$${version}" ]]; then \
		version="$$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')"; \
	fi; \
	ldflags="-s -w -extldflags='-static'"; \
	if [[ -n "$${version}" ]]; then \
		ldflags+=" -X '$(MODULE)/pkg/cmd/blob.version=$${version}'"; \
	fi; \
	flags=(-trimpath -ldflags="$${ldflags}"); \
	if [[ "$${cmd}" == "build" ]]; then \
		flags+=(-o "bin/blob"); \
	fi; \
	echo "$${cmd}ing blob@$${version} $$(go env GOOS)/$$(go env GOARCH) cgo=$$(go env CGO_ENABLED)"; \
	go "$${cmd}" "$${flags[@]}" ./cmd/blob
