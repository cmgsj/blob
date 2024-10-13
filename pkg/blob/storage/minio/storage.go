package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	minioClient  *minio.Client
	bucket       string
	objectPrefix string
}

type StorageOptions struct {
	Address      string
	AccessKey    string
	SecretKey    string
	Bucket       string
	ObjectPrefix string
	Secure       bool
}

func NewStorage(ctx context.Context, opts StorageOptions) (*Storage, error) {
	minioClient, err := minio.New(opts.Address, &minio.Options{
		Creds:  credentials.NewStatic(opts.AccessKey, opts.SecretKey, "", credentials.SignatureDefault),
		Secure: opts.Secure,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(ctx, opts.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("minio bucket %q does not exist", opts.Bucket)
	}

	return &Storage{
		minioClient:  minioClient,
		bucket:       opts.Bucket,
		objectPrefix: util.BlobObjectPrefix(opts.ObjectPrefix),
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	path = util.BlobPath(s.objectPrefix, path)

	objects := s.minioClient.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    path,
		Recursive: true,
	})

	var blobNames []string

	for object := range objects {
		blobNames = append(blobNames, util.BlobPath(strings.TrimPrefix(object.Key, s.objectPrefix)))
	}

	slices.Sort(blobNames)

	return blobNames, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	path := util.BlobPath(s.objectPrefix, name)

	object, err := s.minioClient.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	content, err := io.ReadAll(object)
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
	path := util.BlobPath(s.objectPrefix, name)

	_, err := s.minioClient.PutObject(ctx, s.bucket, path, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	path := util.BlobPath(s.objectPrefix, name)

	err := s.minioClient.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
