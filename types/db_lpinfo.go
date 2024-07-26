package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DbLpInfoRow struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Profile string             `bson:"Profile"`
}
