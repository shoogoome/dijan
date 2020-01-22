package utils

import (
	"bytes"
	"dijan/models"
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

func InitGlobalSystemConfig() {

	yamlFile, err := ioutil.ReadFile("/etc/dijan/config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %s", err.Error()))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalSystemConfig)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal configuration: %s", err.Error()))
	}
	//fmt.Println("这里")
	//var out bytes.Buffer
	//cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cat /etc/hosts"))
	//cmd.Stdout = &out
	//cmd.Stderr = os.Stdout
	//err = cmd.Run()
	//fmt.Println("?", out.String(), err)
	//
	//out = bytes.Buffer{}
	//fmt.Println(fmt.Sprintf("\"ping -c 1 %s.%s | grep PING | awk '{print $3}'\"", os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0", os.Getenv("HEADLESS_SERVICE")))
	//cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("ping -c 1 %s | grep PING | awk '{print $3}'", os.Getenv("HOSTNAME")))
	//cmd.Stdout = &out
	//cmd.Stderr = os.Stdout
	//err = cmd.Run()
	//fmt.Println("?", out.String(), err)
	//fmt.Println("结束")


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
		GlobalSystemConfig.Server.Address = out.String()
PING:
		out = bytes.Buffer{}
		fmt.Println(fmt.Sprintf("%s %s %s", "/bin/sh", "-c", fmt.Sprintf("ping -c 1 %s.%s | grep PING | awk '{print $3}'", os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0", os.Getenv("HEADLESS_SERVICE"))))
		cmd = exec.Command("/bin/sh", "-c", fmt.Sprintf("ping -c 1 %s.%s | grep PING | awk '{print $3}'", os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0", os.Getenv("HEADLESS_SERVICE")))
		cmd.Stdout = &out
		cmd.Stderr = os.Stdout
		err = cmd.Run()
		fmt.Println("?", out.String(), err)
		pingIp := regexp.MustCompile("\\((.*)\\):")
		r := pingIp.FindSubmatch(out.Bytes())
		if len(r) != 2 {
			fmt.Println(len(r))
			fmt.Println(out.String())
			for v, i := range r {
				fmt.Println("!!", v, string(i))
			}
			fmt.Println("[!] 集群ip获取出错， 5s后重试")
			time.Sleep(time.Second * 5)
			goto PING
		} else {
			s1 := strings.Replace(string(r[1]), "\n", "", -1)
			s2 := strings.Replace(s1, " ", "", -1)
			GlobalSystemConfig.Server.ClusterAddress = s2
		}
		fmt.Println("flag", GlobalSystemConfig.Server.HostName, GlobalSystemConfig.Server.Address, GlobalSystemConfig.Server.ClusterAddress)
	}
}
