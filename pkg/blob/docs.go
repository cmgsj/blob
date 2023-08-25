package blob

import (
	_ "embed"

	"github.com/cmgsj/go-lib/openapi"
)

func init() {
	openapi.Must(openapi.RegisterSchema(schema))
}

var (
	//go:embed blob.swagger.json
	docs []byte

	schema = openapi.Schema{
		Name:        ServiceName,
		ContentJSON: docs,
	}
)
