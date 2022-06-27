package utils

import (
	"fmt"
	"runtime"
	"sshcrack/crack"
	"sshcrack/logger"
	"sshcrack/models"
	"sshcrack/utils/hash"
	"sshcrack/vars"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func GenerateTask(ipList []models.IpAddr, users []string, passwords []string) (tasks []models.Service, taskName int) {
	//tasks = make([]models.Service,0)
	tasks = []models.Service{}

	for _, user := range users {
		for _, password := range passwords {
			for _, addr := range ipList {
				service := models.Service{Ip: addr.Ip, Port: addr.Port, Username: user, Password: password}
				tasks = append(tasks, service)
			}
		}
	}

	return tasks, len(tasks)
}

func RunTask(tasks []models.Service) {
	totalTask := len(tasks)
	vars.ProgressBar = pb.StartNew(totalTask)
	// 反引号表示不支持转义
	vars.ProgressBar.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)

	wg := &sync.WaitGroup{}
	//创建一个buffer为vars.ThreadNum的channel
	taskChan := make(chan models.Service, vars.ThreadNum)

	// 创建vars.ThreadNum个协程
	for i := 0; i < vars.ThreadNum; i++ {
		go crackPassword(taskChan, wg)
	}

	// 生产者，不断往taskChan channel发送数据，直到channel阻塞
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
	// 记得关闭
	close(taskChan)

	// 等待所有协程执行结束
	wg.Wait()

	// 停止进度条
	vars.ProgressBar.Finish()

}

func crackPassword(taskChan chan models.Service, wg *sync.WaitGroup) {
	for task := range taskChan {
		if vars.DebugMode {
			logger.Log.Debugf("checking: Ip: %v, Port: %v, UserName: %v, Password: %v, goroutineNum: %v",
				task.Ip, task.Port, task.Username, task.Password, runtime.NumGoroutine())
		}

		k := fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Username)
		// 用于防止密码破解后继续提交破解请求，浪费时间资源
		h := hash.MakeTaskHash(k)
		if hash.CheckTaskHash(h) {
			wg.Done()
			continue
		}

		scanResultChan := make(chan models.ScanResult, 1)
		errChan := make(chan error, 1)
		var scanResult models.ScanResult
		var err error
		// 超时控制
		go func() {
			defer close(scanResultChan)
			defer close(errChan)
			scanResult, err := crack.ScanSsh(task)
			scanResultChan <- scanResult
			errChan <- err
		}()
		select {
		case err = <-errChan:
			scanResult = <-scanResultChan
		case <-time.After(vars.Timeout):
			err = fmt.Errorf("Timeout")
			scanResult = models.ScanResult{
				Service: task,
				Result:  false,
			}
		}

		models.SaveResult(scanResult, err)

		// 增加进度条进度
		vars.ProgressBar.Increment()
		wg.Done()
	}
}
