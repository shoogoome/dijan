package utils

import (
	"bytes"
	"dijan/models"
	"dijan/utils/node"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

var GlobalSystemConfig models.GlobalConfig

func InitGlobalSystemConfig() chan string {

	yamlFile, err := ioutil.ReadFile("/etc/dijan/config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %s", err.Error()))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalSystemConfig)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal configuration: %s", err.Error()))
	}

	if os.Getenv("HEADLESS_SERVICE") == "" {
		GlobalSystemConfig.Server.HostName = os.Getenv("HOSTNAME")
		GlobalSystemConfig.Server.Address = os.Getenv("HOSTNAME")
		GlobalSystemConfig.Server.ClusterAddress = os.Getenv("HOSTNAME")
	} else {
		GlobalSystemConfig.Server.HostName = os.Getenv("HOSTNAME")
		var out bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cat /etc/hosts | grep %s | awk '{print $1}'", os.Getenv("HOSTNAME")))
		cmd.Stdout = &out
		cmd.Stderr = os.Stdout
		err := cmd.Run()
		fmt.Println("?", out.String(), err)
		GlobalSystemConfig.Server.Address = strings.Replace(out.String(), "\n", "", -1)
PING:
		clusterIp, err := nodeUtils.GetClusterIp()
		if err != nil {
			fmt.Println("[!] 集群ip获取出错， 5s后重试")
			time.Sleep(time.Second * 5)
			goto PING
		} else {
			GlobalSystemConfig.Server.ClusterAddress = clusterIp

			// 实时检测0号机情况
			ipChan := make(chan string)
			go func() {
				t := time.NewTimer(time.Second * 5)
				for {
					<-t.C

					clusterIp, err := nodeUtils.GetClusterIp()
					if err == nil && clusterIp != GlobalSystemConfig.Server.ClusterAddress {
						ipChan <- clusterIp
					}
					t.Reset(time.Second * 5)
				}
			}()
			return ipChan
		}
	}
	return make(chan string)
}
