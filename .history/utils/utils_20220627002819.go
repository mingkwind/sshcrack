package utils

import (
	"fmt"
	"net"
	"sshpwd_crack/logger"
	"sshpwd_crack/models"
	"sshpwd_crack/vars"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	AliveAddr []models.IpAddr
	mutex     sync.Mutex
)

func init() {
	AliveAddr = make([]models.IpAddr, 0)
}

func Scan(ctx *cli.Context) (err error) {
	if ctx.IsSet("debug") {
		vars.DebugMode = ctx.Bool("debug")
	}
	if vars.DebugMode {
		logger.Log.Level = logrus.DebugLevel
	}
	if ctx.IsSet("timeout") {
		vars.Timeout = time.Duration(ctx.Int("timeout")) * time.Second
	}
	if ctx.IsSet("thread_num") {
		vars.ThreadNum = ctx.Int("thread_num")
	}
	if ctx.IsSet("ip_list") {
		vars.IpList = ctx.String("ip_list")
	}
	if ctx.IsSet("user_dict") {
		vars.UserDict = ctx.String("user_dict")
	}
	if ctx.IsSet("pass_dict") {
		vars.PassDict = ctx.String("pass_dict")
	}
	if ctx.IsSet("outfile") {
		vars.ResultFile = ctx.String("outfile")
	}

	vars.StartTime = time.Now()
	userDict, uErr := ReadUserDict(vars.UserDict)
	passDict, pErr := ReadPasswordDict(vars.PassDict)
	ipList := ReadIpList(vars.IpList)
	vars.ReScanStrs, _ = ReadReScanStrs("reScanStrs.conf")
	aliveIpList := CheckAlive(ipList)
	if uErr == nil && pErr == nil {
		tasks, _ := GenerateTask(aliveIpList, userDict, passDict)
		RunTask(tasks)
		// 对因为线程数过多造成的断连进行重新测试
		for _, tmptasks := models.GetRestTaskList(); len(tmptasks) != 0; {
			if vars.DebugMode {
				logger.Log.Info(tmptasks)
			}
			vars.TmpService = cache.New(cache.NoExpiration, cache.DefaultExpiration)
			RunTask(tmptasks)
			_, tmptasks = models.GetRestTaskList()
		}

		{
			// 将cache对象以序列化方式存到文件中
			_ = models.SaveResultToFile()
			// 打印结果个数和停止进度条
			models.ResultTotal()
			// 自定义保存结果到txt，方便查看
			_ = models.DumpToFile(vars.ResultFile)
		}
	}
	return err
}

func CheckAlive(ipList []models.IpAddr) []models.IpAddr {
	logger.Log.Infoln("check ip active")

	vars.ProcessBarActive = pb.StartNew(len(ipList))
	vars.ProcessBarActive.SetTemplate(`{{ rndcolor "Checking progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor}}  {{rtime . | rndcolor }}`)

	var wg sync.WaitGroup
	wg.Add(len(ipList))

	for _, addr := range ipList {
		go func(addr models.IpAddr) {
			defer wg.Done()
			SaveAddr(check(addr))
		}(addr)
	}

	wg.Wait()
	vars.ProcessBarActive.Finish()

	return AliveAddr

}

func check(ipAddr models.IpAddr) (bool, models.IpAddr) {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), vars.Timeout)
	if err == nil {
		alive = true
	}
	// 进度条增加
	vars.ProcessBarActive.Increment()
	return alive, ipAddr
}

func SaveAddr(alive bool, ipAddr models.IpAddr) {
	if alive {
		mutex.Lock()
		AliveAddr = append(AliveAddr, ipAddr)
		mutex.Unlock()
	}
}
