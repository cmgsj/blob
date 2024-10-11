package docs

import (
	_ "embed"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
	"github.com/cmgsj/go-lib/swagger"
)

//go:embed docs.swagger.json
var swaggerDocs []byte

func SwaggerSchema() swagger.Schema {
	return swagger.Schema{
		Name:    blobv1.BlobService_ServiceDesc.ServiceName,
		Content: swaggerDocs,
	}
}
