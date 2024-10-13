package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blob struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Content   []byte             `bson:"content,omitempty"`
	UpdatedAt int64              `bson:"updated_at,omitempty"`
}
