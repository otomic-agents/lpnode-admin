package service

import (
	globalvar "admin-panel/global_var"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"admin-panel/utils"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
)

func MD5(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

type LpMonitService struct {
}

func NewLpMonitService() *LpMonitService {
	return &LpMonitService{}
}
func (*LpMonitService) DeployMonitor(task types.DBMonitorListRow) (err error) {
	v := task
	templatePath := fmt.Sprintf("./setup/task/%s/user_setup.yaml", globalvar.SystemEnv)
	outputPath := fmt.Sprintf("./setup/task/%s/setup_%s.yaml", globalvar.SystemEnv, MD5(v.Name))
	log.Println("deploy:", outputPath)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	tmpWriter := &types.TemplateWriter{}
	setupSet := types.MonitorSetupConfig{
		Namespace:       os.Getenv("POD_NAMESPACE"),
		Name:            v.Name,
		Cron:            v.Cron,
		ScriptPath:      v.ScriptPath,
		MongoDbHost:     os.Getenv("MONGODB_HOST"),
		MongoDbPort:     os.Getenv("MONGODB_PORT"),
		MongoDbUser:     os.Getenv("MONGODB_USER"),
		MongoDbPass:     os.Getenv("MONGODBPASS"),
		MongoDbLpStore:  os.Getenv("MONGODB_DBNAME_LP_STORE"),
		MongoDbBusiness: os.Getenv("MONGODB_DBNAME_BUSINESS_HISTORY"),
		RedisHost:       os.Getenv("REDIS_HOST"),
		RedisPort:       os.Getenv("REDIS_PORT"),
		RedisPass:       os.Getenv("REDIS_PASS"),
		RedisDb:         0,
		UserSpacePath:   os.Getenv("USERSPACE_DATA_PATH"),
	}
	err = tmpl.Execute(tmpWriter, setupSet)
	if err != nil {
		logger.System.Errorf("err: %v\n", err)
		return err
	}
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)
	logger.System.Debug(v)

	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"apply", "-f", outputPath})
	if err != nil {
		return err
	}
	filter := bson.M{
		"name": v.Name,
	}
	set := bson.M{
		"$set": bson.M{
			"deploy_message":     cmdRes.Stdout,
			"deploy_err_message": cmdRes.Stderr,
		},
	}
	err = database.Update("main", "monitor_list", filter, set)
	if err != nil {
		return err
	}
	return nil
}

func (*LpMonitService) UnDeployMonitor(task types.DBMonitorListRow) (err error) {
	outputPath := fmt.Sprintf("./setup/task/%s/setup_%s.yaml", globalvar.SystemEnv, MD5(task.Name))
	log.Println(outputPath)
	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"delete", "-f", outputPath})
	log.Println(cmdRes.Stdout)
	log.Println(cmdRes.Stderr)
	return
}
func (*LpMonitService) DeployMonitorRun(scriptName string) (err error) {
	templatePath := fmt.Sprintf("./setup/task/%s/user_execute.yaml", globalvar.SystemEnv)
	outputPath := fmt.Sprintf("./setup/task/%s/user_execute_%s.yaml", globalvar.SystemEnv, scriptName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	tmpWriter := &types.TemplateWriter{}
	filenameall := path.Base(scriptName)
	filesuffix := path.Ext(scriptName)
	fileprefix := filenameall[0 : len(filenameall)-len(filesuffix)]
	v := struct {
		Namespace       string
		ScriptName      string
		JobName         string
		MongoDbHost     string
		MongoDbPort     string
		MongoDbUser     string
		MongoDbPass     string
		MongoDbLpStore  string
		MongoDbBusiness string
		RedisHost       string
		RedisPort       string
		RedisPass       string
		RedisDb         int
		UserSpacePath   string
	}{
		Namespace:       os.Getenv("POD_NAMESPACE"),
		ScriptName:      scriptName,
		JobName:         fileprefix,
		MongoDbHost:     os.Getenv("MONGODB_HOST"),
		MongoDbPort:     os.Getenv("MONGODB_PORT"),
		MongoDbUser:     os.Getenv("MONGODB_USER"),
		MongoDbPass:     os.Getenv("MONGODBPASS"),
		MongoDbLpStore:  os.Getenv("MONGODB_DBNAME_LP_STORE"),
		MongoDbBusiness: os.Getenv("MONGODB_DBNAME_BUSINESS_HISTORY"),
		RedisHost:       os.Getenv("REDIS_HOST"),
		RedisPort:       os.Getenv("REDIS_PORT"),
		RedisPass:       os.Getenv("REDIS_PASS"),
		RedisDb:         0,
		UserSpacePath:   os.Getenv("USERSPACE_DATA_PATH"),
	}

	err = tmpl.Execute(tmpWriter, v)
	if err != nil {
		logger.System.Errorf("err: %v\n", err)
		return err
	}
	os.WriteFile(outputPath, tmpWriter.ByteBuffer, 0755)

	cmdRes, err := utils.ExecuteCMD("kubectl", []string{"apply", "-f", outputPath})
	if err != nil {
		log.Println("ExecuteCMD err:", err.Error())
		return err
	}
	log.Println("DeployMonitorRun", cmdRes.Stdout)
	log.Println("DeployMonitorRun", cmdRes.Stderr)
	return nil
}
