package azblob

import (
	"bytes"
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go-extensions/pkg/errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

const DriverType = "azblob"

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	azblobClient *service.Client
	bucket       string
	objectPrefix string
}

type DriverOptions struct {
	URI string
}

func NewDriver(ctx context.Context, o DriverOptions) (*Driver, error) {
	uri, err := driver.ParseURI(DriverType, o.URI)
	if err != nil {
		return nil, err
	}

	credential, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{})
	if err != nil {
		return nil, err
	}

	azblobClient, err := service.NewClient(uri.Host, credential, &service.ClientOptions{})
	if err != nil {
		return nil, err
	}

	return &Driver{
		azblobClient: azblobClient,
		bucket:       uri.Bucket,
		objectPrefix: uri.ObjectPrefix,
	}, nil
}

func (d *Driver) Bucket() string {
	return d.bucket
}

func (d *Driver) ObjectPrefix() string {
	return d.objectPrefix
}

func (d *Driver) BucketExists(ctx context.Context, bucket string) (bool, error) {
	_, err := d.azblobClient.NewContainerClient(d.bucket).GetProperties(ctx, &container.GetPropertiesOptions{})
	if err != nil {
		if d.IsObjectNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	pager := d.azblobClient.NewContainerClient(d.bucket).NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix: &path,
	})

	var objectNames []string

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, blob := range response.Segment.BlobItems {
			objectNames = append(objectNames, *blob.Name)
		}
	}

	return objectNames, nil
}

func (d *Driver) GetObject(ctx context.Context, name string) ([]byte, error) {
	response, err := d.azblobClient.NewContainerClient(d.bucket).NewBlockBlobClient(name).DownloadStream(ctx, &blob.DownloadStreamOptions{})
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *Driver) PutObject(ctx context.Context, name string, content []byte) error {
	_, err := d.azblobClient.NewContainerClient(d.bucket).NewBlockBlobClient(name).UploadStream(ctx, bytes.NewReader(content), &blockblob.UploadStreamOptions{})

	return err
}

func (d *Driver) DeleteObject(ctx context.Context, name string) error {
	_, err := d.azblobClient.NewContainerClient(d.bucket).NewBlobClient(name).Delete(ctx, &blob.DeleteOptions{})

	return err
}

func (d *Driver) IsObjectNotFound(err error) bool {
	return errors.IsNotFoundErr(err)
}
