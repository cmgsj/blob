install: gen
	go install -ldflags "-X github.com/cmgsj/blob/pkg/version.Version=1.0.0" .

build: gen
	go build -ldflags "-X github.com/cmgsj/blob/pkg/version.Version=1.0.0" .

gen:
	buf format -w && buf generate --exclude-path vendor
	cp pkg/gen/proto/blob/v1/blob.swagger.json pkg/blob
