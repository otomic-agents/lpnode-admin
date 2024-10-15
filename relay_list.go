package adminapiservice

import (
	relaylist "admin-panel/gen/relay_list"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// relayList service example implementation.
// The example methods log the requests and return zero values.
type relayListsrvc struct {
	logger *log.Logger
}

// NewRelayList returns the relayList service implementation.
func NewRelayList(logger *log.Logger) relaylist.Service {
	return &relayListsrvc{logger}
}

// ListRelay implements listRelay.
func (s *relayListsrvc) ListRelay(ctx context.Context) (res *relaylist.ListRelayResult, err error) {

	res = &relaylist.ListRelayResult{}
	err, cursor := database.FindAll("main", "relayAccounts", bson.M{})
	if err != nil {
		err = errors.WithMessage(err, "query database error occur")
		return
	}
	var results []types.DBRelayAccount
	retList := make([]*relaylist.RelayListRelayItem, 0)
	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "handle cursor error occur")
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		retList = append(retList, &relaylist.RelayListRelayItem{
			ID:           ptr.String(result.Id.Hex()),
			Name:         ptr.String(result.Name),
			Profile:      ptr.String(result.Profile),
			LpIDFake:     ptr.String(result.LpIdFake),
			LpNodeAPIKey: ptr.String(result.LpnodeApiKey),
			RelayAPIKey:  ptr.String(result.RelayApiKey),
			RelayURI:     ptr.String(result.RelayUri),
		})
	}
	res.Result = retList
	s.logger.Print("relayAccount.list")
	return
}
