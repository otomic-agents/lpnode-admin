package adminapiservice

import (
	relayaccount "admin-panel/gen/relay_account"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"admin-panel/types"
	"context"
	"log"
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
		err = errors.WithMessage(err, "查询数据库发生了错误")
		return
	}
	var results []types.DBRelayAccount
	retList := make([]*relayaccount.RelayAccountItem, 0)
	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "处理cusor发生了错误")
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
		err = errors.WithMessage(err, "查询数据库发生了错误")
		return
	}
	if rowData.Id.Hex() != types.MongoEmptyIdHex {
		err = errors.New("已经存在的账号")
		return
	}
	registerResult, err := rrs.RegisterAccount(p.Name, ptr.ToString(p.Profile)) // 向relay 发起注册请求
	if err != nil {
		err = errors.WithMessage(err, "向后端注册账号发生了错误")
		return
	}
	log.Println(registerResult)

	ret, err := database.FindOneAndUpdate("main", "relayAccounts", bson.M{
		"name": p.Name,
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
		err = errors.WithMessage(err, "更新数据库记录发生了错误:")
		return
	}
	_id := ""
	if ret.UpsertedID != nil {
		_id = ret.UpsertedID.(primitive.ObjectID).Hex()
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.ID = ptr.String(_id)
	res.Result.Name = ptr.String(p.Name)
	res.Result.RelayAPIKey = ptr.String(registerResult.RelayApiKey)
	res.Result.LpIDFake = ptr.String(registerResult.LpIdFake)
	log.Println(ret)
	return
}
