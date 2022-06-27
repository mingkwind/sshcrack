package models

import (
	"encoding/gob"
	"fmt"
	"os"
	"sshcrack/logger"
	"sshcrack/utils/hash"
	"sshcrack/vars"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

func init() {
	// 要加这个，不加的话cache无法被反序列化，因为cache存储的是interface
	gob.Register(Service{})
	gob.Register(ScanResult{})
}

func SaveResult(result ScanResult, err error) {
	if err == nil && result.Result {

		k := fmt.Sprintf("%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Username)

		h := hash.MakeTaskHash(k)
		hash.SetTaskHash(h)

		if _, found := vars.CacheService.Get(k); !found {
			logger.Log.Infof("Ip: %v, Port: %v, Username: %v, Password: %v", result.Service.Ip,
				result.Service.Port, result.Service.Username, result.Service.Password)
		}
		vars.CacheService.Set(k, result, cache.NoExpiration)
	}

	// 对于非密码错误性质的ssh连接错误，进行重新爆破
	if err != nil {
		if vars.DebugMode {
			logger.Log.Debugf("%v:%v,%v:%v,Error:%v", result.Service.Ip,
				result.Service.Port, result.Service.Username, result.Service.Password, err)
		}
		shouldReScan := false
		for _, str := range vars.ReScanStrs {
			if strings.Contains(err.Error(), str) {
				shouldReScan = true
				break
			}
		}
		if shouldReScan {
			k := fmt.Sprintf("%v-%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Username, result.Service.Password)
			vars.TmpService.Set(k, result, cache.NoExpiration)

		}
	}
}

// 获取非密码错误的任务集，后续进行重新扫描爆破
func GetRestTaskList() (count int, tasks []Service) {
	count = vars.TmpService.ItemCount()
	if count == 0 {
		return count, nil
	} else {
		tmpTasks := []Service{}
		items := vars.TmpService.Items()
		for _, v := range items {
			result := v.Object.(ScanResult)
			// 防止重新扫描
			k := fmt.Sprintf("%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Username)
			h := hash.MakeTaskHash(k)
			if hash.CheckTaskHash(h) {
				continue
			}
			tmpTasks = append(tmpTasks, result.Service)
		}
		return count, tmpTasks
	}
}

func CacheStatus() (count int, items map[string]cache.Item) {
	count = vars.CacheService.ItemCount()
	items = vars.CacheService.Items()
	return count, items
}

func SaveResultToFile() error {
	return vars.CacheService.SaveFile("password_crack.db")
}

//SavaFile读取
// LoadDbFile("password_crack.db")
func LoadDbFile(filename string) (err error) {
	err = vars.CacheService.LoadFile(filename)
	if err == nil {
		items := vars.CacheService.Items()
		for _, v := range items {
			result := v.Object.(ScanResult)
			fmt.Printf("%v:%v|%v:%v\n",
				result.Service.Ip,
				result.Service.Port,
				result.Service.Username,
				result.Service.Password)
		}
	}
	return err
}

func ResultTotal() {
	//vars.ProgressBar.Finish()
	logger.Log.Info(fmt.Sprintf("Finshed scan, total result: %v, used time: %v",
		vars.CacheService.ItemCount(),
		time.Since(vars.StartTime)))
}

func DumpToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, items := CacheStatus()
	for _, v := range items {
		result := v.Object.(ScanResult)
		_, _ = file.WriteString(fmt.Sprintf("%v:%v|%v:%v\n",
			result.Service.Ip,
			result.Service.Port,
			result.Service.Username,
			result.Service.Password),
		)
	}
	return err
}
