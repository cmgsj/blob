package docs

import (
	_ "embed"

	"github.com/cmgsj/go-lib/swagger"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

//go:embed docs.swagger.json
var swaggerDocs []byte

func SwaggerSchema() swagger.Schema {
	return swagger.Schema{
		Name:    blobv1.BlobService_ServiceDesc.ServiceName,
		Content: swaggerDocs,
	}
}
