package adminapiservice

import (
	relayaccount "admin-panel/gen/relay_account"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"admin-panel/types"
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// relayAccount service example implementation.
// The example methods log the requests and return zero values.
type relayAccountsrvc struct {
	logger *log.Logger
}

// NewRelayAccount returns the relayAccount service implementation.
func NewRelayAccount(logger *log.Logger) relayaccount.Service {
	return &relayAccountsrvc{logger}
}

// List implements list.
func (s *relayAccountsrvc) ListAccount(ctx context.Context) (res *relayaccount.ListAccountResult, err error) {
	res = &relayaccount.ListAccountResult{}
	err, cursor := database.FindAll("main", "relayAccounts", bson.M{})
	if err != nil {
		err = errors.WithMessage(err, "query database error occur")
		return
	}
	var results []types.DBRelayAccount
	retList := make([]*relayaccount.RelayAccountItem, 0)
	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "handle cursor error occur")
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		retList = append(retList, &relayaccount.RelayAccountItem{
			ID:           ptr.String(result.Id.Hex()),
			Name:         ptr.String(result.Name),
			Profile:      ptr.String(result.Profile),
			LpIDFake:     ptr.String(result.LpIdFake),
			LpNodeAPIKey: ptr.String(result.LpnodeApiKey),
			RelayAPIKey:  ptr.String(result.RelayApiKey),
		})
	}
	res.Result = retList
	s.logger.Print("relayAccount.list")
	return
}
func (s *relayAccountsrvc) RegisterAccount(ctx context.Context, p *relayaccount.RegisterAccountPayload) (res *relayaccount.RegisterAccountResult, err error) {
	res = &relayaccount.RegisterAccountResult{Result: &relayaccount.RelayAccountItem{}}
	rrs := service.NewRelayRequestService()
	rowData := struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "relayAccounts", bson.M{}, &rowData)
	if err != nil {
		err = errors.WithMessage(err, "query database error occur")
		return
	}
	if rowData.Id.Hex() != types.MongoEmptyIdHex {
		err = errors.New("account already exist")
		return
	}
	lpName := os.Getenv("LP_NAME")
	if lpName == "" {
		err = errors.WithMessage(errors.New("value error"), "cannot get env variable lpname")
		return
	}
	registerResult, err := rrs.RegisterAccount(lpName, ptr.ToString(p.Profile))
	if err != nil {
		err = errors.WithMessage(err, "register account to backend error occur")
		return
	}
	log.Println(registerResult)

	ret, err := database.FindOneAndUpdate("main", "relayAccounts", bson.M{
		"name": lpName,
	}, bson.M{
		"$set": bson.M{
			"profile":      p.Profile,
			"lpnodeApiKey": registerResult.LpnodeApiKey,
			"relayApiKey":  registerResult.RelayApiKey,
			"responseName": registerResult.Name,
			"lpIdFake":     registerResult.LpIdFake,
			"registerAt":   time.Now().UnixNano() / 1e6,
		},
	})
	if err != nil {
		err = errors.WithMessage(err, "update database record error occur:")
		return
	}

	// bcls := service.NewBridgeConfigLogicService()
	// _, err = bcls.ConfigLp()
	// if err != nil {
	// 	err = errors.WithMessage(err, "ConfigLp error occur.")
	// 	return
	// }

	_id := ""
	if ret.UpsertedID != nil {
		_id = ret.UpsertedID.(primitive.ObjectID).Hex()
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.ID = ptr.String(_id)
	res.Result.Name = ptr.String(lpName)
	res.Result.RelayAPIKey = ptr.String(registerResult.RelayApiKey)
	res.Result.LpIDFake = ptr.String(registerResult.LpIdFake)
	log.Println(ret)
	return
}
func (s *relayAccountsrvc) DeleteAccount(ctx context.Context, p *relayaccount.DeleteAccountPayload) (res *relayaccount.DeleteAccountResult, err error) {
	res = &relayaccount.DeleteAccountResult{}
	mongoId, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		err = errors.WithMessage(err, "id format incorrect cannot convert to Mongoid")
		return
	}
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	} = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "relayAccounts", bson.M{"_id": mongoId}, &v)
	if err != nil {
		err = errors.WithMessage(err, "query database error occur")
		return
	}
	if v.Id.Hex() == types.MongoEmptyIdHex {
		err = errors.New("relayaccount not found")
		return
	}
	delCount, err := database.DeleteOne("main", "relayAccounts", bson.M{"_id": mongoId})
	logger.System.Debug("deleted %d accounts", delCount)
	if err != nil {
		return
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("ok")
	res.Result = ptr.String("ok")
	return
}
