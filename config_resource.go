package adminapiservice

import (
	configresource "admin-panel/gen/config_resource"
	database "admin-panel/mongo_database"
	redisbus "admin-panel/redis_bus"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// configResource service example implementation.
// The example methods log the requests and return zero values.
type configResourcesrvc struct {
	logger *log.Logger
}

// NewConfigResource returns the configResource service implementation.
func NewConfigResource(logger *log.Logger) configresource.Service {
	return &configResourcesrvc{logger}
}

// CreateResource implements createResource.
func (s *configResourcesrvc) CreateResource(ctx context.Context, p *configresource.CreateResourcePayload) (res *configresource.CreateResourceResult, err error) {
	res = &configresource.CreateResourceResult{Result: &configresource.ConfigResultIDItem{}}
	sourceStr := fmt.Sprintf("%d", (int64(time.Now().UnixNano())))
	hash := md5.Sum([]byte(sourceStr))
	md5Str := fmt.Sprintf("%x", hash)
	v := struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "configResources", bson.M{
		"clientId": p.ClientID,
	}, &v)
	if err != nil {
		err = errors.WithMessage(err, "query database failed")
		return
	}
	if v.Id.Hex() != types.MongoEmptyIdHex {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "already exist config")
		return
	}
	updateResult, err := database.FindOneAndUpdate("main", "configResources", bson.M{
		"clientId": p.ClientID,
	}, bson.M{
		"$set": bson.M{
			"template":       p.Template,
			"version":        p.Version,
			"appName":        p.AppName,
			"templateResult": "",
			"versionHash":    md5Str,
		},
	})
	if err != nil {
		errors.WithMessage(err, "update database occur error")
		return
	}
	if updateResult.UpsertedID == nil {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "no update document")
		return
	}
	_id := updateResult.UpsertedID.(primitive.ObjectID).Hex()
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.ID = _id
	res.Result.ClientID = ptr.String(p.ClientID)
	s.logger.Print("configResource.createResource")
	return
}

// GetResource implements getResource.
func (s *configResourcesrvc) GetResource(ctx context.Context, p *configresource.GetResourcePayload) (res *configresource.GetResourceResult, err error) {
	res = &configresource.GetResourceResult{Result: &configresource.ResourceConfigResult{}}
	var result = struct {
		ID             primitive.ObjectID `bson:"_id"`
		ClientID       string             `bson:"clientId"`
		AppName        string             `bson:"appName"`
		Template       string             `bson:"template"`
		TemplateResult string             `bson:"templateResult"`
		Version        string             `bson:"version"`
		VersionHash    string             `bson:"versionHash"`
	}{}
	err = database.FindOne("main", "configResources", bson.M{
		"clientId": p.ClientID,
	}, &result)
	if err != nil {
		err = errors.WithMessage(err, "reading database occur error")
		return
	}
	if result.ID.Hex() == types.MongoEmptyIdHex {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "configId is not exist")
		return
	}
	res.Result.ClientID = result.ClientID
	res.Result.AppName = ptr.String(result.AppName)
	res.Result.Template = ptr.String(result.Template)
	res.Result.TemplateResult = ptr.String(result.TemplateResult)
	res.Result.Version = ptr.String(result.Version)
	res.Result.VersionHash = ptr.String(result.VersionHash)
	s.logger.Print("configResource.getResource")
	return
}

// ListResource implements listResource.
func (s *configResourcesrvc) ListResource(ctx context.Context) (res *configresource.ListResourceResult, err error) {
	res = &configresource.ListResourceResult{
		Result: make([]*configresource.ResourceConfigResult, 0)}
	var results []struct {
		ID             primitive.ObjectID `bson:"_id"`
		ClientID       string             `bson:"clientId"`
		AppName        string             `bson:"appName"`
		Template       string             `bson:"template"`
		TemplateResult string             `bson:"templateResult"`
		Version        string             `bson:"version"`
		VersionHash    string             `bson:"versionHash"`
	}
	err, cursor := database.FindAll("main", "configResources", bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		res.Result = append(res.Result, &configresource.ResourceConfigResult{
			ID:             ptr.String(result.ID.Hex()),
			TemplateResult: ptr.String(result.TemplateResult),
			Template:       ptr.String(result.Template),
			ClientID:       result.ClientID,
			AppName:        ptr.String(result.AppName),
			Version:        ptr.String(result.Version),
			VersionHash:    ptr.String(result.VersionHash),
		})
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	s.logger.Print("configResource.listResource")
	return
}

// DeleteResult implements deleteResult.
func (s *configResourcesrvc) DeleteResult(ctx context.Context) (res *configresource.DeleteResultResult, err error) {
	res = &configresource.DeleteResultResult{}
	s.logger.Print("configResource.deleteResult")
	return
}

// EditResult implements editResult.
func (s *configResourcesrvc) EditResult(ctx context.Context, p *configresource.EditResultPayload) (res *configresource.EditResultResult, err error) {
	res = &configresource.EditResultResult{}
	sourceStr := fmt.Sprintf("%d", (int64(time.Now().UnixNano())))
	hash := md5.Sum([]byte(sourceStr))
	md5Str := fmt.Sprintf("%x", hash)
	v := struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "configResources", bson.M{
		"clientId": p.ClientID,
	}, &v)
	if err != nil {
		err = errors.WithMessage(err, "query database failed")
		return
	}
	if v.Id.Hex() == types.MongoEmptyIdHex {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "no update document")
		return
	}
	cps := service.NewCtrlPanelLogicService()
	ammInstallResult, err := cps.GetInstallRow("amm", ptr.ToString(p.AppName))
	if err != nil {
		err = errors.WithMessage(err, "did not find the corresponding installed service, unable to configure resources")
		return
	}
	updateResult, err := database.FindOneAndUpdate("main", "configResources", bson.M{
		"clientId": p.ClientID,
	}, bson.M{
		"$set": bson.M{
			"template":       p.Template,
			"version":        p.Version,
			"appName":        p.AppName,
			"templateResult": p.TemplateResult,
			"versionHash":    md5Str,
		},
	})
	if err != nil {
		errors.WithMessage(err, "updating database occur error")
		return
	}
	if updateResult.UpsertedID == nil && updateResult.ModifiedCount <= 0 {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "no update document")
		return
	}
	go func() {
		service.NewLpCluster().RestartPod(ammInstallResult.Namespace, "amm-", ptr.ToString(p.AppName))
	}()

	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result = ptr.String(v.Id.Hex())
	msg := fmt.Sprintf(`{"type":"configResourceUpdate","payload":{"clientId":"%s","appName":"%s"}}`, v.Id.Hex(), ptr.ToString(p.AppName))
	redisbus.GetRedisBus().PublishEvent("LP_SYSTEM_Notice", msg)
	s.logger.Print("configResource.editResult")
	return
}
