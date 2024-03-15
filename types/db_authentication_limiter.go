package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthenticationLimiterRow struct {
	ID   primitive.ObjectID `bson:"_id"`
	Data string             `bson:"data"`
}
