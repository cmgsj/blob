install: gen
	@go install -ldflags "-s -w -X github.com/cmgsj/blob/pkg/version.Version=1.0.0" .

build: gen
	@go build -ldflags "-X github.com/cmgsj/blob/pkg/version.Version=1.0.0" .

gen:
	@buf format --write proto && buf generate --template proto/buf.gen.yaml proto
