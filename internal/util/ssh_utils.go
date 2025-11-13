package util

import (
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type HostInfo struct {
	Address  string
	Port     int
	Username string
	Password string
	Timeout  time.Duration
}

// Connection 通过ssh连接远程主机
func (host *HostInfo) Connection() (*ssh.Client, error) {
	address := net.JoinHostPort(host.Address, strconv.Itoa(host.Port))
	config := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         host.Timeout,
	}
	if host.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(host.Password))
	}
	// 建立连接
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (host *HostInfo) TestSSHConnection() (bool, error) {

	client, err := host.Connection()
	if err != nil {
		log.Printf("连接失败: %v", err)
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
