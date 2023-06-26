package adminapiservice

import (
	lpmonit "admin-panel/gen/lpmonit"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// lpmonit service example implementation.
// The example methods log the requests and return zero values.
type lpmonitsrvc struct {
	logger *log.Logger
}

// NewLpmonit returns the lpmonit service implementation.
func NewLpmonit(logger *log.Logger) lpmonit.Service {
	return &lpmonitsrvc{logger}
}

// upload user script
func (s *lpmonitsrvc) AddScript(ctx context.Context, p *lpmonit.AddScriptPayload) (res *lpmonit.AddScriptResult, err error) {
	s.logger.Println(ptr.ToString(p.Name))
	s.logger.Println(ptr.ToString(p.ScriptBody))
	errList := utils.IsDNS1123Subdomain(ptr.ToString(p.Name))
	if len(errList) != 0 {
		err = errors.WithMessage(errors.New("name格式不正确"), "00:")
		return
	}
	res = &lpmonit.AddScriptResult{}
	baseScriptName := fmt.Sprintf("user_script_%d.js", time.Now().UnixNano())
	userScriptFile := fmt.Sprintf("/user-script/%s", baseScriptName)
	rawDecodedText, err := base64.StdEncoding.DecodeString(ptr.ToString(p.ScriptBody))
	if err != nil {
		err = errors.WithMessage(err, "解码base64发生了错误")
		return
	}
	var v = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "monitor_list", bson.M{"name": p.Name}, &v)
	if err != nil {
		err = errors.WithMessage(err, "查询数据库发生了错误")
		return
	}
	if v.Id.Hex() != types.MongoEmptyIdHex {
		err = errors.WithMessage(errors.New("exist"), "已经存在的任务")
		return
	}
	err = os.WriteFile(userScriptFile, []byte(rawDecodedText), 0644)
	if err != nil {
		err = errors.WithMessage(err, "无法写入脚本文件")
		return
	}

	log.Println(rawDecodedText, userScriptFile)
	ret, err := database.FindOneAndUpdate("main", "monitor_list", bson.M{"name": p.Name}, bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"cron":        p.Cron,
			"script_path": baseScriptName,
			"task_type":   "user",
			"createAt":    time.Now().UnixNano() / 1e6,
		},
	})
	if err != nil {
		return
	}
	if v, ok := ret.UpsertedID.(primitive.ObjectID); ok {
		res.TaskID = ptr.String(v.Hex())
		res.Result = ptr.String(v.Hex())
	} else {
		res.TaskID = ptr.String("错误的格式")
	}

	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	taskRow := types.DBMonitorListRow{}
	err = database.FindOne("main", "monitor_list", bson.M{"name": p.Name}, &taskRow)
	if err != nil {
		return
	}
	err = service.NewLpMonitService().DeployMonitor(taskRow)
	if err != nil {
		err = errors.WithMessage(err, "An error occurred when deploying to k8s")
		return
	}
	s.logger.Print("lpmonit.add_script")
	return
}
func (s *lpmonitsrvc) ListScript(ctx context.Context) (res *lpmonit.ListScriptResult, err error) {
	res = &lpmonit.ListScriptResult{
		Result: []*lpmonit.LpMointTaskItem{}}
	var results []types.DBMonitorListRow
	filters := bson.M{}
	err, cursor := database.FindAll("main", "monitor_list", filters)
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	for _, v := range results {
		res.Result = append(res.Result, &lpmonit.LpMointTaskItem{
			ID:         ptr.String(v.ID.Hex()),
			Name:       v.Name,
			Cron:       v.Cron,
			CreatedAt:  v.CreatedAt,
			ScriptPath: ptr.String(v.ScriptPath),
			TaskType:   v.TaskType,
		})
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	log.Println("list")
	return
}
func (s *lpmonitsrvc) DeleteScript(ctx context.Context, p *lpmonit.DeleteScriptPayload) (res *lpmonit.DeleteScriptResult, err error) {
	res = &lpmonit.DeleteScriptResult{
		Code:    ptr.Int64(0),
		Message: ptr.String(""),
	}
	log.Println("DeleteScript")
	objectId, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		err = errors.WithMessage(err, "生成MongodbId错误")
		return
	}
	dbr := types.DBMonitorListRow{}
	err = database.FindOne("main", "monitor_list", bson.M{
		"_id": objectId,
	}, &dbr)
	if err != nil {
		err = errors.WithMessage(err, "没有找到要删除的task")
		return
	}
	err = service.NewLpMonitService().UnDeployMonitor(dbr)
	if err != nil {
		err = errors.WithMessage(err, "An error occurred when deploying to k8s")
		log.Println(err)
		res.Message = ptr.String(err.Error())
		err = nil
	}
	if err != nil {
		err = errors.WithMessage(err, "查询Tasklist 发生了错误")
		return
	}

	delCount, err := database.DeleteOne("main", "monitor_list", bson.M{
		"_id": objectId,
	})
	if delCount <= 0 {
		err = errors.New("没有找到删除的记录，无操作.")
		return
	}
	res.Result = ptr.Int64(delCount)
	log.Println(objectId)

	return
}
func (s *lpmonitsrvc) RunScript(ctx context.Context, p *lpmonit.RunScriptPayload) (res *lpmonit.RunScriptResult, err error) {
	log.Println("RunScript")
	baseScriptName := fmt.Sprintf("tmp-script-%d.js", time.Now().UnixNano())
	dir := "/user-script/run-script/"
	os.MkdirAll(dir, 0777)
	savePath := fmt.Sprintf("%s%s", dir, baseScriptName)
	log.Println("savePath is:", savePath)

	rawDecodedText, err := base64.StdEncoding.DecodeString(ptr.ToString(p.ScriptContent))
	if err != nil {
		err = errors.WithMessage(err, "无法解码base64,请检查格式")
		return
	}
	// fmt.Printf("Decoded text: %s\n", rawDecodedText)
	os.WriteFile(savePath, rawDecodedText, 0644)
	service.NewLpMonitService().DeployMonitorRun(baseScriptName)
	res = &lpmonit.RunScriptResult{}
	res.Code = ptr.Int64(0)
	res.Result = ptr.String(baseScriptName)
	res.Message = ptr.String("ok")
	return
}
func (s *lpmonitsrvc) RunResult(ctx context.Context, p *lpmonit.RunResultPayload) (res *lpmonit.RunResultResult, err error) {
	scriptName := ptr.ToString(p.ScriptName)
	log.Println("scriptName:", scriptName)
	res = &lpmonit.RunResultResult{}
	v := struct {
		Id         primitive.ObjectID `bson:"_id" json:"_id"`
		ScriptName string             `bson:"scriptName" json:"scriptName"`
		Stdout     string             `bson:"stdout" json:"stdout"`
		Stderr     string             `bson:"stderr" json:"stderr"`
		ExecResult string             `bson:"execResult" json:"execResult"`
	}{}
	log.Println("filter", scriptName)
	err = database.FindOne("main", "monitor_historys", bson.M{
		"scriptName": scriptName,
	}, &v)
	log.Println(v)
	if err != nil {
		err = errors.WithMessage(err, "读取历史发生了错误")
		return
	}
	if v.Id.Hex() == types.MongoEmptyIdHex {
		res.Code = ptr.Int64(0)
		res.Message = ptr.String("not found")
		res.Result = ptr.String("not fount")
		return
	}
	json_bytes, err := json.Marshal(v)
	if err != nil {
		err = errors.WithMessage(err, "序列化执行结果集发生了错误")
		return
	}
	log.Println(v.Id.Hex())

	res.Code = ptr.Int64(0)
	res.Result = ptr.String(string(json_bytes))
	res.Message = ptr.String(v.ScriptName)
	return
}
