package main

import (
	globalvar "admin-panel/global_var"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func InitMonitor() (err error) {
	log.Println("InitMonitor")
	filters := bson.M{}
	err, cursor := database.FindAll("main", "monitor_list", filters)
	if err != nil {
		return
	}
	var results []types.DBMonitorListRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	log.Println(results)
	err = deployMonitorList(results)
	if err != nil {
		return
	}
	return nil
}

// MD5 hashes using md5 algorithm
func MD5(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}
func deployMonitorList(listData []types.DBMonitorListRow) (err error) {

	for _, v := range listData {
		templatePath := fmt.Sprintf("./setup/task/%s/setup.yaml", globalvar.SystemEnv)
		outputPath := fmt.Sprintf("./setup/task/%s/setup_%s.yaml", globalvar.SystemEnv, MD5(v.Name))
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
		}
		tmpWriter := &types.TemplateWriter{}
		setupSet := types.MonitorSetupConfig{
			Namespace:  os.Getenv("POD_NAMESPACE"),
			Name:       v.Name,
			Corn:       v.Corn,
			ScriptPath: v.ScriptPath,
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
	}

	return
}
