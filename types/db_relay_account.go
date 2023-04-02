package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBRelayAccount struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	LpIdFake     string             `bson:"lpIdFake"`
	LpnodeApiKey string             `bson:"lpnodeApiKey"`
	Profile      string             `bson:"profile"`
	RelayApiKey  string             `bson:"relayApiKey"`
	ResponseName string             `bson:"responseName"`
}
