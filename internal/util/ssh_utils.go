package util

import (
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestSSHConnection(host string, port int, username, password string, privateKey string, timeout time.Duration) (bool, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
	if password != "" {
		config.Auth = append(config.Auth, ssh.Password(password))
	}
	// 建立连接
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return false, err
	}
	defer client.Close()
	// 执行简单的命令
	session, err := client.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	output, err := session.Output("whoami")
	if err != nil {
		return false, err
	}
	log.Printf("连接成功，当前登录的用户是: %s", output)
	return true, nil
}
