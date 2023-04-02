package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type CtrlPanelLogicService struct {
}

func NewCtrlPanelLogicService() *CtrlPanelLogicService {
	return &CtrlPanelLogicService{}
}
func (cps *CtrlPanelLogicService) ListInstallByInstallType(installType string) ([]types.InstallRow, error) {
	filter := bson.M{
		"installType": installType,
		"status": bson.M{
			"$gte": 0,
		},
	}
	emptyList := []types.InstallRow{}
	err, cursor := database.FindAll("main", "install", filter)
	if err != nil {
		return emptyList, err
	}
	var results []types.InstallRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		return emptyList, err
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	return results, nil
}
func (cps *CtrlPanelLogicService) Installed(installType string, name string) bool {
	filter := bson.M{
		"installType": installType,
		"name":        name,
	}
	var ir types.InstallRow
	database.FindOne("main", "install", filter, &ir)
	if ir.ID.Hex() == types.MongoEmptyIdHex {
		return false
	}
	if ir.Status >= 0 {
		return true
	}
	return false
}
func (cps *CtrlPanelLogicService) InstalledByType(installType string) bool {
	filter := bson.M{
		"installType": installType,
	}
	var ir types.InstallRow
	database.FindOne("main", "install", filter, &ir)
	if ir.ID.Hex() == types.MongoEmptyIdHex {
		return false
	}
	if ir.Status >= 0 {
		return true
	}
	return false
}
func (cps *CtrlPanelLogicService) UpdateInstallRow(installType string, name string, setData bson.M) error {
	filter := bson.M{
		"installType": installType,
		"name":        name,
	}
	return database.Update("main", "install", filter, setData)
}
func (cps *CtrlPanelLogicService) MarkInstalled(installType string, name string, source []byte, stdout []byte, stderr []byte) bool {
	filter := bson.M{
		"installType": installType,
		"name":        name,
	}
	set := bson.M{
		"$set": bson.M{
			"lastinstall":  time.Now().UnixNano() / 1e6,
			"status":       1,
			"yaml":         string(source),
			"stdout":       string(stdout),
			"stderr":       string(stderr),
			"configStatus": 0,
		},
	}
	var ir types.InstallRow
	database.FindOneAndUpdate("main", "install", filter, set)
	return ir.ID.Hex() != types.MongoEmptyIdHex
}
func (cps *CtrlPanelLogicService) MarkUninstalled(installType string, name string, source []byte, stdout []byte, stderr []byte) bool {
	filter := bson.M{
		"installType": installType,
		"name":        name,
	}
	set := bson.M{
		"$set": bson.M{
			"status":    -1,
			"yaml":      string(source),
			"un_stdout": string(stdout),
			"un_stderr": string(stderr),
		},
	}
	var ir types.InstallRow
	database.FindOneAndUpdate("main", "install", filter, set)
	return ir.ID.Hex() != types.MongoEmptyIdHex
}

func (cps *CtrlPanelLogicService) GetInstallRow(installType string, name string) (*types.InstallRow, error) {
	filter := bson.M{
		"installType": installType,
		"name":        name,
	}
	log.Println(filter)
	var ir types.InstallRow
	database.FindOne("main", "install", filter, &ir)
	if ir.ID.Hex() == types.MongoEmptyIdHex {
		return nil, errors.New("没有找到安装记录")
	}
	return &ir, nil
}

func (cps *CtrlPanelLogicService) GetInstallRowByInstallType(installType string) (ret []types.InstallRow, err error) {
	filter := bson.M{"installType": installType}
	log.Println(filter)
	var results []types.InstallRow
	err, cur := database.FindAll("main", "install", filter)
	if err != nil {
		return
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cur.Decode(&result)
	}
	ret = results
	return
}

func (cps *CtrlPanelLogicService) MergeSetupConfig(source *types.SetupConfig, setValue *types.SetupConfig) (data *types.SetupConfig) {
	log.Println(setValue, "👥👥")
	if setValue.Deployment.StartBlock != "" {
		source.Deployment.StartBlock = setValue.Deployment.StartBlock
	}
	if setValue.Deployment.RpcUrl != "" {
		source.Deployment.RpcUrl = setValue.Deployment.RpcUrl
	}
	if setValue.Deployment.ConnectionNodeurl != "" {
		source.Deployment.ConnectionNodeurl = setValue.Deployment.ConnectionNodeurl
	}
	if setValue.Deployment.ConnectionWalleturl != "" {
		source.Deployment.ConnectionWalleturl = setValue.Deployment.ConnectionWalleturl
	}
	if setValue.Deployment.ConnectionHelperurl != "" {
		source.Deployment.ConnectionHelperurl = setValue.Deployment.ConnectionHelperurl
	}
	if setValue.Deployment.ConnectionExplorerurl != "" {
		source.Deployment.ConnectionExplorerurl = setValue.Deployment.ConnectionExplorerurl
	}
	if setValue.Deployment.AwsAccessKeyId != "" {
		source.Deployment.AwsAccessKeyId = setValue.Deployment.AwsAccessKeyId
	}
	if setValue.Deployment.AwsSecretAccessKey != "" {
		source.Deployment.AwsSecretAccessKey = setValue.Deployment.AwsSecretAccessKey
	}
	if setValue.Deployment.ContainerPort != "" {
		source.Deployment.ContainerPort = setValue.Deployment.ContainerPort
	}
	if setValue.Deployment.Image != "" {
		source.Deployment.Image = setValue.Deployment.Image
	}
	if len(setValue.Deployment.CustomEnv) > 0 {
		source.Deployment.CustomEnv = setValue.Deployment.CustomEnv
	}
	return source
}
