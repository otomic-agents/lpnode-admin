package service

import (
	"admin-panel/types"
	"errors"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type DeploymentService struct {
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
}

func (*DeploymentService) GetEnv(yamlPath string, envName string) (val string, err error) {
	contentByte, err := os.ReadFile(yamlPath)
	if err != nil {
		return
	}
	content := strings.Split(string(contentByte), "---")

	// Find the Deployment document
	var deploymentIndex int
	for i, doc := range content {
		if strings.Contains(doc, "kind: Deployment") {
			deploymentIndex = i
			break
		}
	}

	dts := types.DeploymentStruct{}
	err = yaml.Unmarshal([]byte(content[deploymentIndex]), &dts)
	if err != nil {
		return
	}

	if len(dts.Spec.Template.Spec.Containers) <= 0 {
		err = errors.New("cannot find container")
		return
	}

	val = ""
	// Loop through all containers to find the environment variable
	for _, container := range dts.Spec.Template.Spec.Containers {
		for _, v := range container.Env {
			log.Println(v.Name, v.Value)
			if v.Name == envName {
				val = v.Value
				return
			}
		}
	}
	return
}

func (*DeploymentService) GetAllEnv(yamlPath string, envName string) (val string, err error) {
	contentByte, err := os.ReadFile(yamlPath)
	if err != nil {
		return
	}
	content := strings.Split(string(contentByte), "---")
	maxIndex := len(content) - 1
	dts := types.DeploymentStruct{}
	dtsByte := []byte(content[maxIndex])
	err = yaml.Unmarshal(dtsByte, &dts)

	if len(dts.Spec.Template.Spec.Containers) <= 0 {
		err = errors.New("cannot find container")
		return
	}

	val = ""
	for _, container := range dts.Spec.Template.Spec.Containers {
		for _, v := range container.Env {
			log.Println(v.Name, v.Value)
			if v.Name == envName {
				val = v.Value
				return
			}
		}
	}
	return
}

func (*DeploymentService) GetEnvList(yamlPath string) (list []struct {
	Name  string
	Value string
}, err error) {
	list = make([]struct {
		Name  string
		Value string
	}, 0)
	contentByte, err := os.ReadFile(yamlPath)
	if err != nil {
		return
	}
	content := strings.Split(string(contentByte), "---")
	maxIndex := len(content) - 1
	dts := types.DeploymentStruct{}
	dtsByte := []byte(content[maxIndex])
	log.Println(string(dtsByte))
	err = yaml.Unmarshal(dtsByte, &dts)
	if len(dts.Spec.Template.Spec.Containers) <= 0 {
		err = errors.New("cannot find container")
		return
	}

	for _, v := range dts.Spec.Template.Spec.Containers[0].Env {
		list = append(list, struct {
			Name  string
			Value string
		}{
			Name: v.Name, Value: v.Value})
	}
	return
}
