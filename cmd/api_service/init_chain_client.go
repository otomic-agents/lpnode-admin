package main

import (
	globalvar "admin-panel/global_var"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitChainClientInstall initializes chain client installation after a 3-minute delay
// Returns:
//   - error: Any error that occurred during initialization
func InitChainClientInstall() (err error) {
	go func() {
		// Wait for 3 minutes before proceeding
		logger.System.Info("Waiting for 3 minutes before initializing chain client installation...", "Delay")
		time.Sleep(3 * time.Minute)
		logger.System.Info("Starting chain client installation initialization", "Process")

		// Query installation records
		var results []types.InstallRow
		opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
		err, cursor := database.FindAllOpt("main", "install", bson.M{
			"installType": "ammClient",
		}, opts)
		if err != nil {
			logger.System.Error("Failed to find installation records: "+err.Error(), "Database")
			return
		}

		// Decode all results
		if err = cursor.All(context.TODO(), &results); err != nil {
			err = errors.WithMessage(err, "cursor all error")
			logger.System.Error("Failed to decode installation records: "+err.Error(), "Database")
			return
		}

		// Process each installation record
		logger.System.Info(fmt.Sprintf("Processing %d chain client installation records", len(results)), "Process")
		for _, result := range results {
			deployChainClient(result)
		}

		logger.System.Info("Chain client installation initialization completed", "Success")
		return
	}()
	return nil
}
func deployChainClient(row types.InstallRow) {
	if row.UpdateResult.LatestConfig.Deployment.Image != "" {
		logger.System.Info("🚀 Deployment Image: "+row.UpdateResult.LatestConfig.Deployment.Image, "📦")
		spew.Dump(row.UpdateResult.LatestConfig.Deployment)
		keysToCheck := []string{"RpcUrl", "StartBlock"}
		// If user has updated rpc, startblock, etc., there will be records in the database
		// Admin will process the latest record during first startup
		if checkAnyDeploymentKeyUpdated(row.UpdateResult.LatestConfig.Deployment, keysToCheck) {
			logger.System.Info("✅ Configuration update detected, continuing deployment process", "🔄")
			doDeployChainClient(row)
		} else {
			logger.System.Warn("⚠️ No configuration update detected, no need to continue deployment", "🛑")
		}
	}
}
func doDeployChainClient(row types.InstallRow) (err error) {
	templatePath := ""
	outputPath := ""
	chainType := ""
	templatePath = fmt.Sprintf("./setup/client/%s/%s_out.yaml", globalvar.SystemEnv, row.Name)
	outputPath = fmt.Sprintf("./setup/client/%s/%s_out_update.yaml", globalvar.SystemEnv, row.Name)
	setupConfig := types.SetupConfig{}
	setupConfig.Deployment.Name = row.UpdateResult.LatestConfig.Deployment.Name
	setupConfig.Deployment.Namespace = row.UpdateResult.LatestConfig.Deployment.Namespace

	updateConfig := types.SetupConfig{}
	updateConfig.Deployment = row.UpdateResult.LatestConfig.Deployment
	updateConfig.Service = row.UpdateResult.LatestConfig.Service

	// Retrieve install context from DB, which contains the latest installed or updated image
	// When users update or version upgrades occur, the install context image will be updated accordingly
	// This ensures we always have the most recent configuration
	dbInstallContext := types.SetupConfig{}
	installContext := row.InstallContext
	if installContext != "" {
		err = json.Unmarshal([]byte(installContext), &dbInstallContext)
		if err != nil {
			return
		}
	}
	setupConfig.Deployment.Image = dbInstallContext.Deployment.Image
	if row.ChainId == 501 {
		chainType = "solana"
	} else {
		chainType = "evm"
	}
	findValueByName := func(envList []struct {
		Name  string
		Value string
	}, name string) string {
		for _, env := range envList {
			if env.Name == name {
				return env.Value
			}
		}
		return ""
	}
	envList, env_err := service.NewLpCluster().DescPodEnv(os.Getenv("NAMESPACE"), fmt.Sprintf("chain-client-%s-%s-%d", chainType, row.Name, row.ChainId))
	if env_err != nil {
		err = errors.WithMessage(errors.New(""), env_err.Error())
		return
	}
	// setupConfig.Deployment.Image = findValueByName(envList,"")
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
	if updateConfig.Deployment.RpcUrl != "" {
		setupConfig.Deployment.RpcUrl = updateConfig.Deployment.RpcUrl
	}
	if updateConfig.Deployment.StartBlock != "" {
		setupConfig.Deployment.StartBlock = updateConfig.Deployment.StartBlock
	}
	logger.System.Info("Template path: " + templatePath)
	logger.System.Info("Output path: " + outputPath)
	spew.Dump(setupConfig)

	if setupConfig.Deployment.RedisHost == "" || setupConfig.Deployment.MongodbAccount == "" {
		err = errors.WithMessage(errors.New(""), "Installation variables are incomplete")
		return
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return
	}
	tmpWriter := &types.TemplateWriter{}
	err = tmpl.Execute(tmpWriter, setupConfig)
	if err != nil {
		return
	}
	logger.System.Debug("Starting to generate update YAML file", outputPath)
	content := string(tmpWriter.ByteBuffer)
	log.Printf("File content to be written:\n%s", content)
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)

	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"apply", "-f", outputPath})
	if err != nil {
		return
	}
	file := ptr.String(string(tmpWriter.ByteBuffer))
	log.Println(file)

	spew.Dump(cmdRes)
	return
}

// CheckAnyDeploymentKeyUpdated checks if any of the specified keys has a non-empty value in the Deployment structure
// Parameters:
//   - deployment: The Deployment structure to check
//   - keys: Array of key names to check
//
// Returns:
//   - bool: Returns true if any of the specified keys has a non-empty value; otherwise returns false
func checkAnyDeploymentKeyUpdated(deployment types.ClientSetupConfigDeployment, keys []string) bool {
	// Use reflection to get the fields of the structure
	val := reflect.ValueOf(deployment)
	typ := val.Type()

	// Create a mapping from field names to field values
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < val.NumField(); i++ {
		fieldMap[typ.Field(i).Name] = val.Field(i)
	}

	// Check each specified key
	for _, key := range keys {
		// Check if the key exists in the structure
		field, exists := fieldMap[key]
		if !exists {
			logger.System.Warn("🔍 Field does not exist: "+key, "⚠️")
			continue
		}

		// Check if the field value is non-empty
		if field.Kind() == reflect.String && field.String() != "" {
			logger.System.Info("🔍 Field "+key+" has been updated: "+field.String(), "✅")
			return true
		} else if field.Kind() == reflect.Slice && !field.IsNil() && field.Len() > 0 {
			logger.System.Info("🔍 Field "+key+" has been updated", "✅")
			return true
		}
	}

	// None of the specified keys have non-empty values
	logger.System.Info("🔍 None of the specified fields have been updated", "❌")
	return false
}
