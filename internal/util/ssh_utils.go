package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
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

type ExecutorResult struct {
	Success  bool
	Output   string
	Error    error
	Host     string
	Command  string
	ExitCode int
	Duration string
}

// BatchConfig 批量执行配置
type BatchConfig struct {
	// 最大并发数
	MaxConcurrent int
	GlobalTimeout time.Duration
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

// Execute 单个主机执行命令
func (host *HostInfo) Execute(command string, isSave bool) (*ExecutorResult, error) {
	// 1. 连接主机
	connection, err := host.Connection()
	if err != nil {
		return &ExecutorResult{
			Success:  err == nil,
			Output:   fmt.Sprintf("连接失败: %v", err),
			Error:    err,
			Host:     host.Address,
			Command:  command,
			ExitCode: 1,
			Duration: "0ms",
		}, err
	}
	defer connection.Close()

	// 2. 建立session链接
	session, err := connection.NewSession()
	if err != nil {
		return &ExecutorResult{
			Success:  err == nil,
			Output:   fmt.Sprintf("建立会话失败: %v", err),
			Error:    err,
			Host:     host.Address,
			Command:  command,
			ExitCode: 1,
			Duration: "0ms",
		}, err
	}
	defer session.Close()

	// 3. 执行命令
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	startTime := time.Now()
	err = session.Run(command)

	exitCode := 0
	fmt.Println("stdout:", stdout.String())
	fmt.Println("stderr:", stderr.String())
	fmt.Println("err:", err)

	result := &ExecutorResult{
		Success:  err == nil,
		Output:   fmt.Sprintf("正在执行命令: %s\n执行成功：%s", command, stdout.String()),
		Error:    err,
		Host:     host.Address,
		Command:  command,
		ExitCode: exitCode,
		Duration: formatDuration(time.Since(startTime)),
	}
	if err != nil {
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitStatus()
		}
		result.Output = stderr.String()
		result.ExitCode = exitCode
	}
	return result, nil
}

// ExecuteCommands 单机机器批量执行命令
func (host *HostInfo) ExecuteCommands(commands []string) ([]*ExecutorResult, error) {
	var results []*ExecutorResult
	for _, command := range commands {
		result, err := host.Execute(command, true)
		if err != nil {
			results = append(results, result)
			continue
		}
		results = append(results, result)
	}
	return results, nil
}

// BatchExecuteCommands 多台机器批量执行多个命令
func BatchExecuteCommands(hosts []*HostInfo, commands []string, config *BatchConfig) (map[string][]*ExecutorResult, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make(map[string][]*ExecutorResult)
	)

	semaphore := make(chan struct{}, config.MaxConcurrent)
	ctx, cancel := context.WithTimeout(context.Background(), config.GlobalTimeout*time.Second)
	defer cancel()
	for _, host := range hosts {
		wg.Add(1)
		go func(host *HostInfo) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			// 执行命令
			executeCommands, err := host.ExecuteCommands(commands)
			if err != nil {
				fmt.Printf("执行命令失败: %v", err)
			}
			mu.Lock()
			results[host.Address] = executeCommands
			mu.Unlock()
		}(host)

	}
	wg.Wait()
	return results, nil
}
