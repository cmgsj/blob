package azblob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go-extensions/pkg/errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	azblobClient *service.Client
	bucket       string
	objectPrefix string
}

type DriverOptions struct {
	URI         string
	AccountName string
	AccountKey  string
}

func NewDriver(ctx context.Context, opts DriverOptions) (*Driver, error) {
	u, err := url.Parse(opts.URI)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid azblob uri %q: host is required", opts.URI)
	}

	var (
		bucket       string
		objectPrefix string
		endpoint     string
	)

	switch u.Scheme {
	case "http", "https":
		path := strings.Split(strings.Trim(u.Path, "/"), "/")

		if len(path) < 3 {
			return nil, fmt.Errorf("invalid azblob uri %q: bucket is required", opts.URI)
		}

		bucket = path[2]

		if len(path) > 3 {
			objectPrefix = strings.Join(path[3:], "/")
		}

		endpoint = fmt.Sprintf("%s://%s/%s/%s/", u.Scheme, u.Host, path[0], path[1])

	default:
		return nil, fmt.Errorf("invalid azblob uri %q: unknown scheme", opts.URI)
	}

	credential, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{})
	if err != nil {
		return nil, err
	}

	azblobClient, err := service.NewClient(endpoint, credential, &service.ClientOptions{})
	if err != nil {
		return nil, err
	}

	return &Driver{
		azblobClient: azblobClient,
		bucket:       bucket,
		objectPrefix: objectPrefix,
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
