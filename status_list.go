package adminapiservice

import (
	statuslist "admin-panel/gen/status_list"
	database "admin-panel/mongo_database"
	"admin-panel/redis_database"
	"admin-panel/types"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
)

// statusList service example implementation.
// The example methods log the requests and return zero values.
type statusListsrvc struct {
	logger *log.Logger
}

// NewStatusList returns the statusList service implementation.
func NewStatusList(logger *log.Logger) statuslist.Service {
	return &statusListsrvc{logger}
}

// StatList implements statList.
func (s *statusListsrvc) StatList(ctx context.Context) (res *statuslist.StatListResult, err error) {
	res = &statuslist.StatListResult{}
	res.Result = make([]*statuslist.StatusListItem, 0)
	err, cursor := database.FindAll("main", "install", bson.M{
		"status": bson.M{"$gt": 0},
	})
	if err != nil {
		return
	}
	var results []types.InstallRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		statusKey := ""
		statusBody := ""
		errMessage := ""
		log.Println(len(result.EnvList), "🧊🧊🧊🧊🧊")
		for _, v := range result.EnvList {
			if v.Name == "STATUS_KEY" {
				statusKey = v.Value
			}
		}
		if statusKey != "" {
			redisBody, readErr := redis_database.GetStatusDb().GetString(statusKey)
			statusBody = redisBody
			if readErr != nil {
				errMessage = readErr.Error()
			}
		}

		res.Result = append(res.Result, &statuslist.StatusListItem{
			Name:        ptr.String(result.Name),
			InstallType: ptr.String(result.InstallType),
			StatusKey:   ptr.String(statusKey),
			StatusBody:  ptr.String(statusBody),
			ErrMessage:  ptr.String(errMessage),
		})
	}

	s.logger.Print("statusList.statList")
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	return
}
