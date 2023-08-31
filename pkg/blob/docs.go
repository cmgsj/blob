package blob

import (
	_ "embed"

	"github.com/cmgsj/go-lib/openapi"
)

var (
	//go:embed blob.swagger.json
	docs []byte

	OpenapiSchema = openapi.Schema{
		Name:        ServiceName,
		ContentJSON: docs,
	}
)
