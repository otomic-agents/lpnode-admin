package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBMonitorListRow struct {
	ID            primitive.ObjectID `bson:"_id"`
	Corn          string             `bson:"corn"`
	Name          string             `bson:"name"`
	ScriptPath    string             `bson:"script_path"`
	TaskType      string             `bson:"task_type"`
	DeployMessage string             `bson:"deploy_message"`
}
