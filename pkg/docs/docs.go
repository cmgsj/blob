package docs

import (
	_ "embed"

	"github.com/cmgsj/go-lib/openapi"
)

var (
	//go:embed openapi.swagger.json
	docs []byte

	OpenapiSchema = openapi.Schema{
		Name:        "openapi",
		ContentJSON: docs,
	}
)
