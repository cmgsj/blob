package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

const BlobStorageFolder = "blobs"

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	minioClient *minio.Client
	bucketName  string
	folder      string
}

func NewStorage(ctx context.Context, uri string, opts *minio.Options) (*Storage, error) {
	if opts == nil {
		opts = &minio.Options{}
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid minio uri %q: host is required", uri)
	}

	var endpoint string
	var bucketName string
	var bucketPrefix string

	switch u.Scheme {
	case "minio":
		path := strings.Split(strings.Trim(u.Path, "/"), "/")

		if len(path) < 1 {
			return nil, fmt.Errorf("invalid minio uri %q: bucket is required", uri)
		}

		bucketName = path[0]

		if len(path) > 1 {
			bucketPrefix = strings.Join(path[1:], "/")
		}

		endpoint = u.Host

	default:
		return nil, fmt.Errorf("invalid minio uri %q: unknown scheme", uri)
	}

	minioClient, err := minio.New(endpoint, opts)
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("invalid minio uri %q: bucket does not exist", uri)
	}

	folder := BlobStorageFolder

	if bucketPrefix != "" {
		folder = util.BlobPath(bucketPrefix, folder)
	}

	return &Storage{
		minioClient: minioClient,
		bucketName:  bucketName,
		folder:      folder,
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	path = util.BlobPath(s.folder, path)

	objects := s.minioClient.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix:    path,
		Recursive: true,
	})

	var blobNames []string

	for object := range objects {
		blobNames = append(blobNames, util.BlobPath(strings.TrimPrefix(object.Key, s.folder)))
	}

	slices.Sort(blobNames)

	return blobNames, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	path := util.BlobPath(s.folder, name)

	object, err := s.minioClient.GetObject(ctx, s.bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	err = object.Close()
	if err != nil {
		return nil, err
	}

	info, err := object.Stat()
	if err != nil {
		return nil, err
	}

	return &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: info.LastModified.Unix(),
	}, nil
}

func (s *Storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	path := util.BlobPath(s.folder, name)

	_, err := s.minioClient.PutObject(ctx, s.bucketName, path, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	path := util.BlobPath(s.folder, name)

	err := s.minioClient.RemoveObject(ctx, s.bucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func IsStorage(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "minio":
		return true
	}

	return false
}
