package crack

import (
	"fmt"
	"net"

	"sshcrack/models"
	"sshcrack/vars"

	"golang.org/x/crypto/ssh"
)

func ScanSsh(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout: vars.Timeout,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		if err == nil {
			errRet := session.Run("echo myssh")
			if errRet == nil {
				defer session.Close()
				result.Result = true
			}
		}
	}

	return result, err
}
