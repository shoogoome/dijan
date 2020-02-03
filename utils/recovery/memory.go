package recoveryUtils

import (
	"bytes"
	"dijan/core/cache"
	"dijan/models"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func MemoryRecovery() {
	go func() {
		for {
			// 每天凌晨3点开始
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 3,0,0,0,next.Location())
			t := time.NewTimer(next.Sub(now))
			<- t.C

			fmt.Println("[!] 开始清理过期缓存...")
			tu := false
			scan := make(chan int, 1)
			cpu := make(chan int, 1)
			go getCpuUser(cpu)
			s := cache.Conn.NewScanner()
			for s.Scan() {
				scan <- 1
				select {
				case <-scan:
					var storageValue models.StorageValue
					if err := json.Unmarshal(s.Value(), &storageValue); err != nil {
						continue
					}
					if storageValue.TTL != -1 && time.Now().Unix() >= storageValue.TTL  {
						cache.Conn.Del(s.Key())
					}
				case <-cpu:
					tu = true
					break
				}
				if tu {
					break
				}
			}
			s.Close()
		}
	}()
}

func getCpuUser(cpu chan int) {
	t := time.NewTimer(time.Second*5)
	cpuUser := regexp.MustCompile("^%CPU(.*)$")
	for {
		<-t.C

		var out bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("ps u %d | awk '{print $3}'", os.Getpid()))
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err == nil {
			s := strings.Replace(out.String(), "\n", "", -1)
			s = strings.Replace(s, " ", "", -1)
			r := cpuUser.FindSubmatch([]byte(s))
			if len(r) == 2 {
				if user, err := strconv.ParseFloat(string(r[1]), 32); err != nil || user > 20 {
					cpu <- 1
				}
			} else {
				cpu <- 1
			}

		} else {
			cpu <- 1
			break
		}
		t.Reset(time.Second * 5)
	}
}
