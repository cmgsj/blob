package mongodb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

const BlobStorageCollection = "blobs"

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	collection *mongo.Collection
}

func NewStorage(ctx context.Context, uri string, opts ...*options.ClientOptions) (*Storage, error) {
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}

	if cs.Database == "" {
		return nil, fmt.Errorf("invalid mongodb uri %q: database is required", uri)
	}

	mongoClient, err := mongo.Connect(ctx, append(opts, options.Client().ApplyURI(uri))...)
	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Storage{
		collection: mongoClient.Database(cs.Database).Collection(BlobStorageCollection),
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	cursor, err := s.collection.Find(ctx, bson.M{
		"name": bson.M{
			"$regex": fmt.Sprintf("^%s", util.BlobNamePrefix(path)),
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mbs []Blob

	err = cursor.All(ctx, &mbs)
	if err != nil {
		return nil, err
	}

	var blobNames []string

	for _, b := range mbs {
		blobNames = append(blobNames, b.Name)
	}

	slices.Sort(blobNames)

	return blobNames, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	result := s.collection.FindOne(ctx, bson.M{
		"name": name,
	})

	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, storage.ErrBlobNotFound
		}

		return nil, err
	}

	var mb Blob

	err = result.Decode(&mb)
	if err != nil {
		return nil, err
	}

	return &blobv1.Blob{
		Name:      mb.Name,
		Content:   mb.Content,
		UpdatedAt: mb.UpdatedAt,
	}, nil
}

func (s *Storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	_, err := s.collection.InsertOne(ctx, Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{
		"name": name,
	})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return storage.ErrBlobNotFound
	}

	return nil
}

func IsStorage(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "mongodb", "mongodb+srv":
		return true
	}

	return false
}
