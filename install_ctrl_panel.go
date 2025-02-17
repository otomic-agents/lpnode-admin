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
		err = errors.New("install client already exist")
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
	// template merge output
	setupConfig := types.SetupConfig{
		Service: types.ClientSetupConfigService{},
		Deployment: types.ClientSetupConfigDeployment{
			Name:                  p.SetupConfig.Type,
			RunEnv:                globalvar.SystemEnv, // runEnv
			CustomEnv:             customEnv,
			Namespace:             os.Getenv("POD_NAMESPACE"),
			OsSystemServer:        os.Getenv("OS_SYSTEM_SERVER"),
			OsApiSecret:           os.Getenv("OS_API_SECRET"),
			OsApiKey:              os.Getenv("OS_API_KEY"),
			RedisHost:             os.Getenv("REDIS_HOST"),
			RedisPort:             os.Getenv("REDIS_PORT"),
			RedisPass:             os.Getenv("REDIS_PASS"),
			MongodbHost:           os.Getenv("MONGODB_HOST"),
			MongodbPort:           os.Getenv("MONGODB_PORT"),
			MongodbAccount:        os.Getenv("MONGODB_USER"),
			MongodbPass:           os.Getenv("MONGODB_PASS"),
			MongodbDbnameLpStore:  os.Getenv("MONGODB_DBNAME_LP_STORE"),
			MongodbDbnameHistory:  os.Getenv("MONGODB_DBNAME_BUSINESS_HISTORY"),
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

	log.Println("current install flag", p.SetupConfig.Install)
	baseChainRow, err := bds.GetChainRowByName(p.SetupConfig.Type)
	if err != nil {
		err = errors.New("no chain row")
		return
	}
	chainId := baseChainRow.ChainId
	// then go into install action
	if !p.SetupConfig.Install { //when no need to install
		return
	}
	logger.System.Debug("start generating install yaml file..")
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)
	deployService := service.NewDeploymentService()
	chainType, err := deployService.GetEnv(outputPath, "CHAIN_TYPE")
	if err != nil {
		err = errors.WithMessage(err, "configuration missing chain type config")
		return
	}
	envList, err := deployService.GetEnvList(outputPath)
	if err != nil {
		err = errors.WithMessage(err, "did not get envlist")
		return
	}
	if chainType == "" {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "configuration missing chain type config:")
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
		err = errors.New("no service installed before")
		return
	}

	logger.System.Debug("yaml", installRow.Yaml)
	outputPath := fmt.Sprintf("./setup/client/%s/%s_out_uninstall.yaml", globalvar.SystemEnv, ptr.ToString(p.SetupConfig.Type))
	os.WriteFile(outputPath, []byte(installRow.Yaml), 0766)
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"delete", "-f", outputPath})
	if err != nil {
		logger.System.Errorf("uninstalling service error: %s", err)
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
		err = errors.New("no service installed before")
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
		err = fmt.Errorf("no adapted installer")
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
	findValueByName := func(envList []struct {
		Name  string
		Value string
	}, name string) string {
		for _, env := range envList {
			if env.Name == name {
				log.Println("find .... a ", env.Value)
				return env.Value
			}
		}
		// If not found, return empty string or error value
		return ""
	}
	if p.SetupConfig.InstallType == "ammClient" {
		var chainType string
		if installRow.ChainId == 501 {
			chainType = "solana"
		} else {
			chainType = "evm"
		}

		envList, env_err := service.NewLpCluster().DescPodEnv(os.Getenv("NAMESPACE"), fmt.Sprintf("chain-client-%s-%s-%d", chainType, installRow.Name, installRow.ChainId))
		fmt.Println(envList, env_err, fmt.Sprintf("|%s|", os.Getenv("NAMESPACE")))
		setupConfig.Deployment.RedisHost = findValueByName(envList, "REDIS_HOST")
		setupConfig.Deployment.RedisPort = findValueByName(envList, "REDIS_PORT")
		setupConfig.Deployment.RedisPass = findValueByName(envList, "REDIS_PASSWORD")
		setupConfig.Deployment.OsApiKey = findValueByName(envList, "OS_API_KEY")
		setupConfig.Deployment.OsApiSecret = findValueByName(envList, "OS_API_SECRET")
		setupConfig.Deployment.OsSystemServer = findValueByName(envList, "OS_SYSTEM_SERVER")
		setupConfig.Deployment.RpcUrl = findValueByName(envList, "RPC_URL")
		setupConfig.Deployment.StartBlock = findValueByName(envList, "START_BLOCK")
	}
	if p.SetupConfig.InstallType == "amm" {
		envList, env_err := service.NewLpCluster().DescPodEnv(os.Getenv("NAMESPACE"), fmt.Sprintf("amm-%s", installRow.Name))
		fmt.Println(envList, env_err, fmt.Sprintf("|%s|", os.Getenv("NAMESPACE")))
		setupConfig.Deployment.RedisHost = findValueByName(envList, "REDIS_HOST")
		setupConfig.Deployment.RedisPort = findValueByName(envList, "REDIS_PORT")
		setupConfig.Deployment.RedisPass = findValueByName(envList, "REDIS_PASSWORD")
		setupConfig.Deployment.MongodbHost = findValueByName(envList, "MONGODB_HOST")
		setupConfig.Deployment.MongodbPort = findValueByName(envList, "MONGODB_PORT")
		setupConfig.Deployment.MongodbAccount = findValueByName(envList, "MONGODB_ACCOUNT")
		setupConfig.Deployment.MongodbPass = findValueByName(envList, "MONGODB_PASSWORD")
		setupConfig.Deployment.MongodbDbnameLpStore = findValueByName(envList, "MONGODB_DBNAME_LP_STORE")
		setupConfig.Deployment.MongodbDbnameHistory = findValueByName(envList, "MONGODB_DBNAME_HISTORY")
		setupConfig.Deployment.OsApiKey = findValueByName(envList, "OS_API_KEY")
		setupConfig.Deployment.OsApiSecret = findValueByName(envList, "OS_API_SECRET")
		setupConfig.Deployment.OsSystemServer = findValueByName(envList, "OS_SYSTEM_SERVER")
	}
	if p.SetupConfig.InstallType == "market" {
		envList, env_err := service.NewLpCluster().DescPodEnv(os.Getenv("NAMESPACE"), fmt.Sprintf("amm-market-%s", installRow.Name))
		fmt.Println(envList, env_err, fmt.Sprintf("|%s|", os.Getenv("NAMESPACE")))
		setupConfig.Deployment.RedisHost = findValueByName(envList, "REDIS_HOST")
		setupConfig.Deployment.RedisPort = findValueByName(envList, "REDIS_PORT")
		setupConfig.Deployment.RedisPass = findValueByName(envList, "REDIS_PASSWORD")
		setupConfig.Deployment.MongodbHost = findValueByName(envList, "MONGODB_HOST")
		setupConfig.Deployment.MongodbPort = findValueByName(envList, "MONGODB_PORT")
		setupConfig.Deployment.MongodbAccount = findValueByName(envList, "MONGODB_ACCOUNT")
		setupConfig.Deployment.MongodbPass = findValueByName(envList, "MONGODB_PASS")
		setupConfig.Deployment.MongodbDbnameLpStore = findValueByName(envList, "MONGODB_DBNAME_LP_STORE")
	}
	if p.SetupConfig.InstallType == "userApp" {
		envList, env_err := service.NewLpCluster().DescPodEnv(os.Getenv("NAMESPACE"), fmt.Sprintf("user-app-%s", installRow.Name))
		fmt.Println(envList, env_err, fmt.Sprintf("|%s|", os.Getenv("NAMESPACE")))
		setupConfig.Deployment.RedisHost = findValueByName(envList, "REDIS_HOST")
		setupConfig.Deployment.RedisPort = findValueByName(envList, "REDIS_PORT")
		setupConfig.Deployment.RedisPass = findValueByName(envList, "REDIS_PASSWORD")
		setupConfig.Deployment.MongodbHost = findValueByName(envList, "MONGODB_HOST")
		setupConfig.Deployment.MongodbPort = findValueByName(envList, "MONGODB_PORT")
		setupConfig.Deployment.MongodbAccount = findValueByName(envList, "MONGODB_ACCOUNT")
		setupConfig.Deployment.MongodbPass = findValueByName(envList, "MONGODB_PASS")
		setupConfig.Deployment.MongodbDbnameLpStore = findValueByName(envList, "MONGODB_DBNAME_LP_STORE")
	}

	log.Println("000______", updateConfig.Deployment)
	setupConfig = *cps.MergeSetupConfig(&setupConfig, &updateConfig)
	logger.System.Debug("🤼🤼🤼🤼", len(setupConfig.Deployment.CustomEnv))

	tmpWriter := &types.TemplateWriter{}

	log.Println(setupConfig)
	err = tmpl.Execute(tmpWriter, setupConfig)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	logger.System.Debug("UpdateDeployment ")
	if !p.SetupConfig.Update { //when no need to install
		return
	}
	logger.System.Debug("start generating update yaml file", outputPath)
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)
	if p.SetupConfig.InstallType == "ammClient" {
		chainType, err = service.NewDeploymentService().GetEnv(outputPath, "CHAIN_TYPE")
		if err != nil {
			err = errors.WithMessage(err, "configuration missing chain type config")
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
		err = errors.New("already installed, cannot install repeatedly")
		return
	}
	if p.SetupConfig.InstallType == "market" {
		marketInstalled := cps.InstalledByType(p.SetupConfig.InstallType)
		if marketInstalled {
			err = errors.WithMessage(utils.GetNoEmptyError(err), "market app can only install one")
			return
		}
	}
	templatePath := fmt.Sprintf("./setup/%s/%s/deployment.yaml", p.SetupConfig.InstallType, globalvar.SystemEnv)
	os.MkdirAll(fmt.Sprintf("./setup/%s/%s/%s", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
	outputPath := fmt.Sprintf("./setup/%s/%s/%s/deployment_install.yaml", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name)
	tmpl, err := template.ParseFiles(templatePath) // create template object
	if err != nil {
		return
	}
	// handle custom env variables
	customEnv := make([]types.AmmSetupConfigDeploymentCustomEnv, 0)
	if len(p.SetupConfig.CustomEnv) > 0 {
		for _, v := range p.SetupConfig.CustomEnv {
			customEnv = append(customEnv, types.AmmSetupConfigDeploymentCustomEnv{Key: ptr.ToString(v.Key), Value: ptr.ToString(v.Value)})
		}
	}
	// initialize config
	setupConfig := types.AmmSetupConfig{
		Service: types.AmmSetupConfigService{},
		Deployment: types.AmmSetupConfigDeployment{
			OsSystemServer:       os.Getenv("OS_SYSTEM_SERVER"),
			OsApiSecret:          os.Getenv("OS_API_SECRET"),
			OsApiKey:             os.Getenv("OS_API_KEY"),
			RedisHost:            os.Getenv("REDIS_HOST"),
			MongodbHost:          os.Getenv("MONGODB_HOST"),
			MongodbPass:          os.Getenv("MONGODB_PASS"),
			RedisPass:            os.Getenv("REDIS_PASS"),
			RedisPort:            os.Getenv("REDIS_PORT"),
			MongodbPort:          os.Getenv("MONGODB_PORT"),
			Namespace:            os.Getenv("POD_NAMESPACE"),
			MongodbAccount:       os.Getenv("MONGODB_USER"),
			MongodbDbnameLpStore: os.Getenv("MONGODB_DBNAME_LP_STORE"),
			MongodbDbnameHistory: os.Getenv("MONGODB_DBNAME_BUSINESS_HISTORY"),
			Name:                 p.SetupConfig.Name,
			CustomEnv:            customEnv,
			Image:                p.SetupConfig.ImageRepository,
			ContainerPort:        ptr.ToString(p.SetupConfig.ContainerPort),
		},
	}
	tmpWriter := &types.TemplateWriter{}
	// render template
	err = tmpl.Execute(tmpWriter, setupConfig)
	if err != nil {
		logger.System.Errorf("err: %v\n", err)
		return
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.Template = ptr.String(string(tmpWriter.ByteBuffer))
	if !p.SetupConfig.Install {
		return
	}
	logger.System.Debug("start installing amm")
	logger.System.Debug("start generating install yaml file")
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
	// save install context
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
		err = errors.New("no service installed before")
		return
	}
	logger.System.Debugf("uninstalling %s", p.SetupConfig.InstallType)
	outputPath := fmt.Sprintf("./setup/%s/%s/%s/deployment_uninstall.yaml", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name)
	os.MkdirAll(fmt.Sprintf("./setup/%s/%s/%s", p.SetupConfig.InstallType, globalvar.SystemEnv, p.SetupConfig.Name), os.ModePerm)
	os.WriteFile(outputPath, []byte(installRow.Yaml), 0766)
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"delete", "-f", outputPath})
	if err != nil {
		logger.System.Errorf("uninstall service occur error:%s", err)
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
