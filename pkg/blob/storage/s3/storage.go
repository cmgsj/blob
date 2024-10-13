package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
)

const ObjectPrefix = "blobs"

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	storage.Storage
	s3Client     *s3.Client
	bucket       string
	objectPrefix string
}

type StorageOptions struct {
	AccessKey    string
	SecretKey    string
	Bucket       string
	ObjectPrefix string
}

func NewStorage(ctx context.Context, opts StorageOptions) (*Storage, error) {
	s3Client := s3.New(s3.Options{
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     opts.AccessKey,
				SecretAccessKey: opts.SecretKey,
			}, nil
		}),
	})

	_, err := s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(opts.Bucket),
	})
	if err != nil {
		return nil, err
	}

	if opts.ObjectPrefix == "" {
		opts.ObjectPrefix = ObjectPrefix
	} else {
		opts.ObjectPrefix = util.BlobPath(opts.ObjectPrefix, ObjectPrefix)
	}

	return &Storage{
		s3Client:     s3Client,
		bucket:       opts.Bucket,
		objectPrefix: ObjectPrefix,
	}, nil
}
