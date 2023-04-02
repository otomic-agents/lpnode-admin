package adminapiservice

import (
	installctrlpanel "admin-panel/gen/install_ctrl_panel"
	globalvar "admin-panel/global_var"
	"admin-panel/logger"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
)

// installCtrlPanel service example implementation.
// The example methods log the requests and return zero values.
type installCtrlPanelsrvc struct {
	logger *log.Logger
}

// NewInstallCtrlPanel returns the installCtrlPanel service implementation.
func NewInstallCtrlPanel(logger *log.Logger) installctrlpanel.Service {
	return &installCtrlPanelsrvc{logger}
}

// InstallLpClient implements installLpClient.
func (s *installCtrlPanelsrvc) InstallLpClient(ctx context.Context, p *installctrlpanel.InstallLpClientPayload) (res *installctrlpanel.InstallLpClientResult, err error) {
	res = &installctrlpanel.InstallLpClientResult{Result: &struct {
		Template  *string
		CmdStdout *string
		CmdStderr *string
	}{}}
	cps := service.NewCtrlPanelLogicService()
	bds := service.NewBaseDataLogicService()
	installed := cps.Installed("ammClient", p.SetupConfig.Type)
	if installed {
		err = errors.New("已经安装过了，不能再次重复安装....")
		return
	}
	// return
	templatePath := fmt.Sprintf("./setup/client/%s/%s_out.yaml", globalvar.SystemEnv, p.SetupConfig.Type)
	outputPath := fmt.Sprintf("./setup/client/%s/%s_out_install.yaml", globalvar.SystemEnv, p.SetupConfig.Type)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	customEnv := make([]types.AmmSetupConfigDeploymentCustomEnv, 0)
	if len(p.SetupConfig.CustomEnv) > 0 {
		for _, v := range p.SetupConfig.CustomEnv {
			customEnv = append(customEnv, types.AmmSetupConfigDeploymentCustomEnv{Key: ptr.ToString(v.Key), Value: ptr.ToString(v.Value)})
		}
	}
	// 模板合并输出
	setupConfig := types.SetupConfig{
		Service: types.ClientSetupConfigService{},
		Deployment: types.ClientSetupConfigDeployment{
			Name:                  p.SetupConfig.Type,
			RunEnv:                globalvar.SystemEnv, // 运行时环境变量
			CustomEnv:             customEnv,
			Namespace:             os.Getenv("POD_NAMESPACE"),
			OsSystemServer:        os.Getenv("OS_SYSTEM_SERVER"),
			OsApiSecret:           os.Getenv("OS_API_SECRET"),
			OsApiKey:              os.Getenv("OS_API_KEY"),
			Image:                 p.SetupConfig.ImageRepository,
			RpcUrl:                ptr.ToString(p.SetupConfig.RPCURL),
			StartBlock:            ptr.ToString(p.SetupConfig.StartBlock),
			ConnectionNodeurl:     ptr.ToString(p.SetupConfig.ConnectionNodeurl),
			ConnectionWalleturl:   ptr.ToString(p.SetupConfig.ConnectionWalleturl),
			ConnectionHelperurl:   ptr.ToString(p.SetupConfig.ConnectionHelperurl),
			ConnectionExplorerurl: ptr.ToString(p.SetupConfig.ConnectionExplorerurl),
			AwsAccessKeyId:        ptr.ToString(p.SetupConfig.AwsAccessKeyID),
			AwsSecretAccessKey:    ptr.ToString(p.SetupConfig.AwsSecretAccessKey),
		},
	}
	tmpWriter := &types.TemplateWriter{}
	err = tmpl.Execute(tmpWriter, setupConfig)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	res.Code = ptr.Int64(0)

	res.Message = ptr.String("")
	res.Result.Template = ptr.String(string(tmpWriter.ByteBuffer))

	log.Println("当前的Install flag", p.SetupConfig.Install)
	baseChainRow, err := bds.GetChainRowByName(p.SetupConfig.Type)
	if err != nil {
		err = errors.New("没有查询到正确的Chain数据")
		return
	}
	chainId := baseChainRow.ChainId
	// 之后进入安装动作
	if !p.SetupConfig.Install { //无需安装的时候
		return
	}
	logger.System.Debug("开始生成安装的Yaml文件..")
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)
	deployService := service.NewDeploymentService()
	chainType, err := deployService.GetEnv(outputPath, "CHAIN_TYPE")
	if err != nil {
		err = errors.WithMessage(err, "配置中缺少链类型的配置")
		return
	}
	envList, err := deployService.GetEnvList(outputPath)
	if err != nil {
		err = errors.WithMessage(err, "没有获取到envList")
		return
	}
	if chainType == "" {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "配置中缺少链类型的配置:")
		return
	}
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"apply", "-f", outputPath})
	if err != nil {
		return
	}
	res.Result.CmdStdout = &cmdRes.Stdout
	res.Result.CmdStderr = &cmdRes.Stderr
	serviceName := utils.ParseServiceNameFromDeplayMessage(cmdRes.Stdout)
	installContext, err := json.Marshal(setupConfig)
	if err != nil {
		return
	}

	cps.MarkInstalled("ammClient", p.SetupConfig.Type, tmpWriter.ByteBuffer, []byte(cmdRes.Stdout), []byte(cmdRes.Stderr))
	cps.UpdateInstallRow("ammClient", p.SetupConfig.Type, bson.M{"$set": bson.M{
		"serviceName":    serviceName,
		"envList":        envList,
		"chainType":      chainType,
		"namespace":      setupConfig.Deployment.Namespace,
		"installContext": string(installContext),
		"lastInstall":    int64(time.Now().UnixNano() / 1e6),
		"chainId":        chainId,
	}})
	s.logger.Print("ctrlPanel.installLpClient")
	return
}

func (s *installCtrlPanelsrvc) UninstallLpClient(ctx context.Context, p *installctrlpanel.UninstallLpClientPayload) (res *installctrlpanel.UninstallLpClientResult, err error) {
	res = &installctrlpanel.UninstallLpClientResult{Result: &struct {
		Template  *string
		CmdStdout *string
		CmdStderr *string
	}{}}
	cps := service.NewCtrlPanelLogicService()
	installRow, queryErr := cps.GetInstallRow("ammClient", ptr.ToString(p.SetupConfig.Type))
	if queryErr != nil {
		err = queryErr
		return
	}
	if installRow.Status <= 0 {
		err = errors.New("没有已经安装过的服务")
		return
	}

	logger.System.Debug("yaml", installRow.Yaml)
	outputPath := fmt.Sprintf("./setup/client/%s/%s_out_uninstall.yaml", globalvar.SystemEnv, ptr.ToString(p.SetupConfig.Type))
	os.WriteFile(outputPath, []byte(installRow.Yaml), 0766)
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"delete", "-f", outputPath})
	if err != nil {
		logger.System.Errorf("卸载service 时发生了错误:%s", err)
		//return
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.CmdStdout = ptr.String(cmdRes.Stdout)
	res.Result.CmdStderr = ptr.String(cmdRes.Stderr)
	cps.UpdateInstallRow("ammClient", ptr.ToString(p.SetupConfig.Type), bson.M{"$set": bson.M{
		"installContext": "{}",
	}})

	cps.MarkUninstalled("ammClient", ptr.ToString(p.SetupConfig.Type), []byte(installRow.Yaml), []byte(cmdRes.Stdout), []byte(cmdRes.Stderr))
	logger.System.Debug("uninstall")
	return
}

func (s *installCtrlPanelsrvc) ListInstall(ctx context.Context, p *installctrlpanel.ListInstallPayload) (res *installctrlpanel.ListInstallResult, err error) {
	res = &installctrlpanel.ListInstallResult{}
	cps := service.NewCtrlPanelLogicService()
	list, err := cps.ListInstallByInstallType(p.InstallType)
	if err != nil {
		return
	}
	for _, v := range list {
		res.Result = append(res.Result, &installctrlpanel.CtrlDeploayItem{
			Name:           ptr.String(v.Name),
			Yaml:           ptr.String(v.Yaml),
			Status:         ptr.Int64(v.Status),
			InstallType:    ptr.String(v.InstallType),
			InstallContext: ptr.String(v.InstallContext),
		})
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	logger.System.Debug("list installed")
	return
}
func (s *installCtrlPanelsrvc) UpdateDeployment(ctx context.Context, p *installctrlpanel.UpdateDeploymentPayload) (res *installctrlpanel.UpdateDeploymentResult, err error) {
	res = &installctrlpanel.UpdateDeploymentResult{Result: &struct {
		CmdStdout *string
		CmdStderr *string
		Template  *string
	}{}}
	cps := service.NewCtrlPanelLogicService()

	installRow, queryErr := cps.GetInstallRow(p.SetupConfig.InstallType, p.SetupConfig.Name)
	if queryErr != nil {
		err = queryErr
		return
	}
	if installRow.Status <= 0 {
		err = errors.New("没有已经安装过的服务")
		return
	}
	templatePath := ""
	outputPath := ""
	chainType := ""
	if p.SetupConfig.InstallType == "ammClient" {
		templatePath = fmt.Sprintf("./setup/client/%s/%s_out.yaml", globalvar.SystemEnv, p.SetupConfig.Name)
		outputPath = fmt.Sprintf("./setup/client/%s/%s_out_update.yaml", globalvar.SystemEnv, p.SetupConfig.Name)
	}
	if p.SetupConfig.InstallType == "amm" {
		templatePath = fmt.Sprintf("./setup/amm/%s/deployment.yaml", globalvar.SystemEnv)
		os.MkdirAll(fmt.Sprintf("./setup/amm/%s/%s", globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
		outputPath = fmt.Sprintf("./setup/amm/%s/%s/deployment_update.yaml", globalvar.SystemEnv, p.SetupConfig.Name)
	}
	if p.SetupConfig.InstallType == "market" {
		templatePath = fmt.Sprintf("./setup/market/%s/deployment.yaml", globalvar.SystemEnv)
		os.MkdirAll(fmt.Sprintf("./setup/%s/%s/%s", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
		outputPath = fmt.Sprintf("./setup/market/%s/%s/deployment_update.yaml", globalvar.SystemEnv, p.SetupConfig.Name)
	}
	if p.SetupConfig.InstallType == "userApp" {
		templatePath = fmt.Sprintf("./setup/userApp/%s/deployment.yaml", globalvar.SystemEnv)
		os.MkdirAll(fmt.Sprintf("./setup/%s/%s/%s", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
		outputPath = fmt.Sprintf("./setup/userApp/%s/%s/deployment_update.yaml", globalvar.SystemEnv, p.SetupConfig.Name)
	}
	if templatePath == "" {
		err = fmt.Errorf("没有适配的安装程序")
		return
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.System.Error(err)
		return
	}
	setupConfig := types.SetupConfig{}
	setupConfig.Deployment.Name = p.SetupConfig.Name
	updateConfig := types.SetupConfig{}

	installContext := ptr.ToString(p.SetupConfig.InstallContext)
	if installContext != "" {
		err = json.Unmarshal([]byte(installContext), &updateConfig)
		if err != nil {
			return
		}
	}
	err = json.Unmarshal([]byte(installRow.InstallContext), &setupConfig)
	if err != nil {
		return
	}
	log.Println(setupConfig)
	setupConfig = *cps.MergeSetupConfig(&setupConfig, &updateConfig)
	logger.System.Debug("🤼🤼🤼🤼", len(setupConfig.Deployment.CustomEnv))

	tmpWriter := &types.TemplateWriter{}

	log.Println(setupConfig, "0000000000")
	err = tmpl.Execute(tmpWriter, setupConfig)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	logger.System.Debug("UpdateDeployment ")
	if !p.SetupConfig.Update { //无需安装的时候
		return
	}
	logger.System.Debug("开始生成更新Yaml文件..", outputPath)
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)
	if p.SetupConfig.InstallType == "ammClient" {
		chainType, err = service.NewDeploymentService().GetEnv(outputPath, "CHAIN_TYPE")
		if err != nil {
			err = errors.WithMessage(err, "配置中缺少链类型的配置")
			return
		}
	}

	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"apply", "-f", outputPath})
	if err != nil {
		return
	}
	res.Result.Template = ptr.String(string(tmpWriter.ByteBuffer))
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.CmdStdout = &cmdRes.Stdout
	res.Result.CmdStderr = &cmdRes.Stderr
	installContextByte, encodeErr := json.Marshal(setupConfig)
	if encodeErr != nil {
		err = encodeErr
		return
	}

	cps.UpdateInstallRow(p.SetupConfig.InstallType, p.SetupConfig.Name, bson.M{"$set": bson.M{
		"installContext": string(installContextByte),
		"chainType":      chainType,
		"updateResult": bson.M{
			"stdout":     res.Result.CmdStdout,
			"stderr":     res.Result.CmdStderr,
			"updateYaml": string(tmpWriter.ByteBuffer),
			"lastupdate": int64(time.Now().UnixNano() / 1e6),
		},
	}})
	return
}
func (s *installCtrlPanelsrvc) InstallDeployment(ctx context.Context, p *installctrlpanel.InstallDeploymentPayload) (res *installctrlpanel.InstallDeploymentResult, err error) {
	res = &installctrlpanel.InstallDeploymentResult{
		Result: &installctrlpanel.InstallDeploymentDataResult{},
	}
	cps := service.NewCtrlPanelLogicService()
	installed := cps.Installed(p.SetupConfig.InstallType, p.SetupConfig.Name)
	if installed {
		err = errors.New("已经安装过了，不能再次重复安装")
		return
	}
	if p.SetupConfig.InstallType == "market" {
		marketInstalled := cps.InstalledByType(p.SetupConfig.InstallType)
		if marketInstalled {
			err = errors.WithMessage(utils.GetNoEmptyError(err), "market应用只能安装一个")
			return
		}
	}
	templatePath := fmt.Sprintf("./setup/%s/%s/deployment.yaml", p.SetupConfig.InstallType, globalvar.SystemEnv)
	os.MkdirAll(fmt.Sprintf("./setup/%s/%s/%s", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
	outputPath := fmt.Sprintf("./setup/%s/%s/%s/deployment_install.yaml", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name)
	tmpl, err := template.ParseFiles(templatePath) // 创建模版对象
	if err != nil {
		return
	}
	// 处理自定义的环境变量
	customEnv := make([]types.AmmSetupConfigDeploymentCustomEnv, 0)
	if len(p.SetupConfig.CustomEnv) > 0 {
		for _, v := range p.SetupConfig.CustomEnv {
			customEnv = append(customEnv, types.AmmSetupConfigDeploymentCustomEnv{Key: ptr.ToString(v.Key), Value: ptr.ToString(v.Value)})
		}
	}
	// 初始化配置
	setupConfig := types.AmmSetupConfig{
		Service: types.AmmSetupConfigService{},
		Deployment: types.AmmSetupConfigDeployment{
			Namespace:     os.Getenv("POD_NAMESPACE"),
			Name:          p.SetupConfig.Name,
			CustomEnv:     customEnv,
			Image:         p.SetupConfig.ImageRepository,
			ContainerPort: ptr.ToString(p.SetupConfig.ContainerPort),
		},
	}
	tmpWriter := &types.TemplateWriter{}
	// 渲染模版
	err = tmpl.Execute(tmpWriter, setupConfig)
	if err != nil {
		logger.System.Errorf("err: %v\n", err)
		return
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.Template = ptr.String(string(tmpWriter.ByteBuffer))
	if !p.SetupConfig.Install { //无需安装的时候
		return
	}
	logger.System.Debug("开始安装Amm")
	logger.System.Debug("开始生成安装的Yaml文件..")
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)
	envList, err := service.NewDeploymentService().GetEnvList(outputPath)
	if err != nil {
		return
	}
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"apply", "-f", outputPath})
	if err != nil {
		return
	}
	installContext, err := json.Marshal(setupConfig)
	if err != nil {
		return
	}
	res.Result.CmdStdout = &cmdRes.Stdout
	res.Result.CmdStderr = &cmdRes.Stderr
	cps.MarkInstalled(p.SetupConfig.InstallType, p.SetupConfig.Name, tmpWriter.ByteBuffer, []byte(cmdRes.Stdout), []byte(cmdRes.Stderr))
	// 保存安装上下文
	cps.UpdateInstallRow(p.SetupConfig.InstallType, p.SetupConfig.Name, bson.M{
		"$set": bson.M{
			"envList":        envList,
			"namespace":      setupConfig.Deployment.Namespace,
			"installContext": string(installContext),
			"lastinstall":    time.Now().UnixNano() / 1e6,
		},
	})
	return
}
func (s *installCtrlPanelsrvc) UninstallDeployment(ctx context.Context, p *installctrlpanel.UninstallDeploymentPayload) (res *installctrlpanel.UninstallDeploymentResult, err error) {
	res = &installctrlpanel.UninstallDeploymentResult{
		Result: &installctrlpanel.UnInstallDeploymentDataResult{},
	}
	cps := service.NewCtrlPanelLogicService()
	installRow, queryErr := cps.GetInstallRow(p.SetupConfig.InstallType, p.SetupConfig.Name)
	if queryErr != nil {
		err = queryErr
		return
	}
	if installRow.Status <= 0 {
		err = errors.New("没有已经安装过的服务")
		return
	}
	logger.System.Debugf("卸载%s", p.SetupConfig.InstallType)
	outputPath := fmt.Sprintf("./setup/%s/%s/%s/deployment_uninstall.yaml", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name)
	os.MkdirAll(fmt.Sprintf("./setup/%s/%s/%s", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
	os.WriteFile(outputPath, []byte(installRow.Yaml), 0766)
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"delete", "-f", outputPath})
	if err != nil {
		logger.System.Errorf("卸载服务发生了错误:%s", err)
		//return
	}
	res.Result.CmdStdout = ptr.String(cmdRes.Stdout)
	res.Result.CmdStderr = ptr.String(cmdRes.Stderr)
	res.Code = ptr.Int64(0)
	cps.UpdateInstallRow(p.SetupConfig.InstallType, p.SetupConfig.InstallType, bson.M{
		"$set": bson.M{
			"installContext": "{}",
			"last_uninstall": time.Now().UnixNano() / 1e6,
		},
	})
	cps.MarkUninstalled(p.SetupConfig.InstallType, p.SetupConfig.Name, []byte(installRow.Yaml), []byte(cmdRes.Stdout), []byte(cmdRes.Stderr))
	return
}
