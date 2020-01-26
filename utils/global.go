package utils

import (
	"bytes"
	"dijan/models"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
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
		GlobalSystemConfig.Server.Address = strings.Replace(out.String(), "\n", "", -1)
PING:
		clusterIp, err := getClusterIp()
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

					clusterIp, err := getClusterIp()
					if err == nil && clusterIp != GlobalSystemConfig.Server.ClusterAddress {
						fmt.Println("[!] 集群ip发生变更")
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

func getClusterIp() (string, error){

	out := bytes.Buffer{}
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("ping -c 1 %s.%s | grep PING | awk '{print $3}'", os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0", os.Getenv("HEADLESS_SERVICE")))
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	pingIp := regexp.MustCompile("\\((.*)\\):")
	r := pingIp.FindSubmatch(out.Bytes())
	if len(r) != 2 {
		return "", errors.New("获取集群ip失败")
	} else {
		return strings.Replace(string(r[1]), "\n", "", -1), nil
	}
}
