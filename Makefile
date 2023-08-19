install: gen
	go install -ldflags "-X github.com/cmgsj/blob/pkg/version.Version=1.0.0" ./cmd/blob

build: gen
	go build -ldflags "-X github.com/cmgsj/blob/pkg/version.Version=1.0.0" ./cmd/blob

gen:
	buf format -w && buf generate --exclude-path vendor
