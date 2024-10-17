SHELL := /bin/bash

MODULE := $$(go list -m)
VERSION := 1.0.0
SWAGGER_UI_VERSION :=

.PHONY: default
default: fmt install

.PHONY: fmt
fmt:
	@find . -type f -name "*.go" ! -path "./pkg/gen/*" ! -path "./vendor/*" | while read -r file; do \
		go fmt "$${file}" 2>&1 | grep -v "is a program, not an importable package"; \
		goimports -w -local $(MODULE) "$${file}"; \
	done

.PHONY: generate
generate: generate/buf generate/swagger

.PHONY: generate/buf
generate/buf:
	@rm -rf pkg/gen; \
	find swagger/dist -type f -name '*.swagger.json' -delete; \
	buf format --write; \
	buf lint; \
	buf breaking --against "https://$(MODULE).git#branch=main"; \
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
