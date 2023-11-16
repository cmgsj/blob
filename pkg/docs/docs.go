package docs

import (
	_ "embed"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/go-lib/openapi"
)

//go:embed openapi.swagger.json
var docs []byte

func OpenapiSchema() openapi.Schema {
	return openapi.Schema{
		Name:        blobv1.BlobService_ServiceDesc.ServiceName,
		ContentJSON: docs,
	}
}
