package utils

import (
	"dijan/models"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var GlobalSystemConfig models.GlobalConfig

func InitGlobalSystemConfig() {

	yamlFile, err := ioutil.ReadFile("/etc/dijan/config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %s", err.Error()))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalSystemConfig)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal configuration: %s", err.Error()))
	}

	if os.Getenv("HEADLESS_SERVICE") == "" {
		GlobalSystemConfig.Server.Address = os.Getenv("HOSTNAME")
		GlobalSystemConfig.Server.ClusterAddress = os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0"
	} else {
		GlobalSystemConfig.Server.Address = fmt.Sprintf("%s.%s", os.Getenv("HOSTNAME"), os.Getenv("HEADLESS_SERVICE"))
		GlobalSystemConfig.Server.ClusterAddress = fmt.Sprintf("%s.%s", os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0", os.Getenv("HEADLESS_SERVICE"))
	}
}
