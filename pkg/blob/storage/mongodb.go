package storage

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/cmgsj/blob/pkg/blob"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

const MongoDBCollectionName = "blobs"

var _ blob.Storage = (*MongoDBStorage)(nil)

func NewMongoDBStorage(ctx context.Context, uri string, opts ...*options.ClientOptions) (*MongoDBStorage, error) {
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}

	if cs.Database == "" {
		return nil, fmt.Errorf("must specify mongodb database")
	}

	mongoClient, err := mongo.Connect(ctx, append(opts, options.Client().ApplyURI(uri))...)
	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDBStorage{
		collection: mongoClient.Database(cs.Database).Collection(MongoDBCollectionName),
	}, nil
}

type MongoDBStorage struct {
	collection *mongo.Collection
}

func (s *MongoDBStorage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	var blobs []string

	cursor, err := s.collection.Find(ctx, bson.M{
		"name": bson.M{
			"$regex": "^" + cleanBlobPrefix(path),
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mblobs []mongodbBlob

	err = cursor.All(ctx, &mblobs)
	if err != nil {
		return nil, err
	}

	for _, b := range mblobs {
		if strings.HasPrefix(b.name, path) {
			blobs = append(blobs, b.name)
		}
	}

	slices.Sort(blobs)

	return blobs, nil
}

func (s *MongoDBStorage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	result := s.collection.FindOne(ctx, bson.M{
		"name": name,
	})

	err := result.Err()
	if err != nil {
		return nil, err
	}

	var mblob mongodbBlob

	err = result.Decode(&mblob)
	if err != nil {
		return nil, err
	}

	b := &blobv1.Blob{
		Name:      mblob.name,
		Content:   mblob.content,
		UpdatedAt: mblob.updatedAt,
	}

	return b, nil
}

func (s *MongoDBStorage) WriteBlob(ctx context.Context, name string, content []byte) error {
	b := mongodbBlob{
		name:      name,
		content:   content,
		updatedAt: time.Now().Unix(),
	}

	_, err := s.collection.InsertOne(ctx, b)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoDBStorage) RemoveBlob(ctx context.Context, name string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{
		"name": name,
	})
	if err != nil {
		return err
	}

	return nil
}

type mongodbBlob struct {
	name      string `bson:"name"`
	content   []byte `bson:"content"`
	updatedAt int64  `bson:"updated_at"`
}
