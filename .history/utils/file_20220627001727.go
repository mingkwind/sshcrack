package utils

import (
	"bufio"
	"os"
	"sshpwd_crack/logger"
	"sshpwd_crack/models"
	"strconv"
	"strings"
)

func ReadIpList(filename string) (ipList []models.IpAddr) {
	ipListFile, err := os.Open(filename)
	if err != nil {
		logger.Log.Fatal("Open ip List file err, %v", err)
	}

	// 记得关闭文件
	defer ipListFile.Close()
	//先逐行读取存储到scanner
	scanner := bufio.NewScanner(ipListFile)
	// 然后用Split进行分割
	scanner.Split(bufio.ScanLines)
	// 遍历Split后的数组
	for scanner.Scan() {
		// 格式为：
		// 127.0.0.1:22
		// 8.8.8.8
		ipPort := strings.TrimSpace(scanner.Text())
		var ip string
		var port int
		if strings.Contains(ipPort, ":") {
			t := strings.Split(ipPort, ":")
			ip = t[0]
			port, _ = strconv.Atoi(t[1])
		} else {
			ip = ipPort
			port = 22
		}
		addr := models.IpAddr{Ip: ip, Port: port}
		ipList = append(ipList, addr)
	}
	return ipList
}

func ReadUserDict(userDict string) (users []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		logger.Log.Fatalf("Open user dict file err,%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users, err
}

func ReadPasswordDict(passDict string) (passwords []string, err error) {
	file, err := os.Open(passDict)
	if err != nil {
		logger.Log.Fatalf("Open password dict file err,%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if password != "" {
			passwords = append(passwords, password)
		}
	}
	return passwords, err
}

func ReadReScanStrs(filename string) (reScanStrs []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		logger.Log.Fatalf("Open reScanStrs file err,%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		reScanStr := strings.TrimSpace(scanner.Text())
		if reScanStr != "" {
			reScanStrs = append(reScanStrs, reScanStr)
		}
	}
	return reScanStrs, err
}
