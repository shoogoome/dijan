package nodeUtils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// 获取集群主节点ip地址
func GetClusterIp() (string, error){

	out := bytes.Buffer{}
	fmt.Println(fmt.Sprintf("%s %s %s", "/bin/sh", "-c", fmt.Sprintf("ping -c 1 %s.%s | grep PING | awk '{print $3}'", os.Getenv("HOSTNAME")[:len(os.Getenv("HOSTNAME")) - 1] + "0", os.Getenv("HEADLESS_SERVICE"))))
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
