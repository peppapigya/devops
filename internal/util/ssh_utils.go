package util

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type HostInfo struct {
	ID       int
	Address  string
	Port     int
	Username string
	Password string
	Timeout  time.Duration
}

type ExecutorResult struct {
	HostId   int
	Success  bool
	Output   string
	Error    error
	Host     string
	Command  string
	ExitCode int
	Duration int64
}

// BatchConfig 批量执行配置
type BatchConfig struct {
	// 最大并发数
	MaxConcurrent int
	GlobalTimeout int
	PerCmdTimeout int
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

// TestSSHConnection 测试ssh链接
func (host *HostInfo) TestSSHConnection() (bool, error) {

	client, err := host.Connection()
	if err != nil {
		log.Printf("连接失败: %v", err)
		return false, err
	}
	defer func() { _ = client.Close() }()
	// 执行简单的命令
	session, err := client.NewSession()
	if err != nil {
		return false, err
	}
	defer func() { _ = session.Close() }()

	output, err := session.Output("whoami")
	if err != nil {
		return false, err
	}
	log.Printf("连接成功，当前登录的用户是: %s", output)
	return true, nil
}

// Execute 单个主机执行命令
func (host *HostInfo) Execute(command string, isSave bool) (*ExecutorResult, error) {
	startTime := time.Now()
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
			Duration: time.Since(startTime).Milliseconds(),
		}, err
	}
	defer func() { _ = connection.Close() }()

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
			Duration: time.Since(startTime).Milliseconds(),
		}, err
	}
	defer func() { _ = session.Close() }()

	// 3. 执行命令
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)

	exitCode := 0
	fmt.Println("stdout:", stdout.String())
	fmt.Println("stderr:", stderr.String())
	fmt.Println("err:", err)

	result := &ExecutorResult{
		Success:  err == nil,
		HostId:   host.ID,
		Output:   fmt.Sprintf("正在执行命令: %s\n执行成功：%s", command, stdout.String()),
		Error:    err,
		Host:     host.Address,
		Command:  command,
		ExitCode: exitCode,
		Duration: time.Since(startTime).Milliseconds(),
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
		if err != nil || result.Success == false || result.ExitCode != 0 {
			results = append(results, result)
			break
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
	ctx := context.Background()
	var cancel context.CancelFunc
	// 只有全局超时时间时才启用全局超时限制
	if config.GlobalTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(config.GlobalTimeout)*time.Second)
		defer cancel()
	}

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
	if ctx.Err() != nil {
		return nil, fmt.Errorf("全局调度超时")
	}
	return results, nil
}

// ---------------------- 流式执行相关类型与方法 ----------------------

// StreamEvent 是流式输出的事件结构体（方案 3）
type StreamEvent struct {
	Host    string
	Command string
	Line    string
	IsError bool
	Time    time.Time
}

// ExecuteStream 在单台主机上执行单个命令，并实时通过回调返回 stdout/stderr 的每一行。
// ctx 控制单个命令的超时或取消。回调不能阻塞过久（建议快速异步处理）。
func (host *HostInfo) ExecuteStream(ctx context.Context, command string, onEvent func(StreamEvent)) *ExecutorResult {
	res := &ExecutorResult{
		Success:  false,
		Output:   "",
		Error:    nil,
		Host:     host.Address,
		Command:  command,
		ExitCode: 1,
		Duration: 0,
	}

	start := time.Now()
	conn, err := host.Connection()
	if err != nil {
		res.Output = fmt.Sprintf("连接失败: %v", err)
		res.Error = err
		res.Duration = time.Since(start).Milliseconds()
		return res
	}
	defer func() { _ = conn.Close() }()

	sess, err := conn.NewSession()
	if err != nil {
		res.Output = fmt.Sprintf("建立会话失败: %v", err)
		res.Error = err
		res.Duration = time.Since(start).Milliseconds()
		return res
	}
	defer func() { _ = sess.Close() }()

	stdoutPipe, err := sess.StdoutPipe()
	if err != nil {
		res.Output = fmt.Sprintf("获取 stdout 管道失败: %v", err)
		res.Error = err
		res.Duration = time.Since(start).Milliseconds()
		return res
	}
	stderrPipe, err := sess.StderrPipe()
	if err != nil {
		res.Output = fmt.Sprintf("获取 stderr 管道失败: %v", err)
		res.Error = err
		res.Duration = time.Since(start).Milliseconds()
		return res
	}

	// 合并 buf 保存最终输出
	var outBuf, errBuf bytes.Buffer

	// 读取流的工人
	readLine := func(r io.Reader, isErr bool) {
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadString(' ')
			if len(line) > 0 {
				if isErr {
					errBuf.WriteString(line)
				} else {
					outBuf.WriteString(line)
				}
				e := StreamEvent{Host: host.Address, Command: command, Line: strings.TrimRight(line, ""), IsError: isErr, Time: time.Now()}
				// 回调尽量非阻塞
				go func(ev StreamEvent) {
					defer func() {
						_ = recover()
					}()
					onEvent(ev)
				}(e)
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				// 非 EOF 错误也退出
				break
			}
		}
	}

	// 开始命令执行
	if err := sess.Start(command); err != nil {
		res.Output = fmt.Sprintf("命令启动失败: %v", err)
		res.Error = err
		res.Duration = time.Since(start).Milliseconds()
		return res
	}

	// 启动读取 goroutine
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		readLine(stdoutPipe, false)
	}()
	go func() {
		defer wg.Done()
		readLine(stderrPipe, true)
	}()

	// 等待命令完成或 ctx 取消
	done := make(chan error, 1)
	go func() {
		done <- sess.Wait()
	}()

	select {
	case <-ctx.Done():
		// 尝试发送信号终止（若支持）
		_ = sess.Signal(ssh.SIGKILL)
		res.Output = fmt.Sprintf("命令被取消或超时: %v", ctx.Err())
		res.Error = ctx.Err()
		res.Duration = time.Since(start).Milliseconds()
		// 等待读取 goroutine 结束
		wg.Wait()
		return res
	case err := <-done:
		// 命令正常/异常结束
		wg.Wait()
		res.Duration = time.Since(start).Milliseconds()
		res.Output = fmt.Sprintf("stdout:%sstderr:%s", outBuf.String(), errBuf.String())
		if err == nil {
			res.Success = true
			res.ExitCode = 0
			return res
		}
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			res.ExitCode = exitErr.ExitStatus()
		} else {
			res.ExitCode = 1
		}
		res.Error = err
		return res
	}
}

// ExecuteCommandsStream 在单台主机上顺序执行一系列命令，并为每条命令触发 onEvent
func (host *HostInfo) ExecuteCommandsStream(commands []string, perCmdTimeout time.Duration, onEvent func(StreamEvent)) ([]*ExecutorResult, error) {
	var results []*ExecutorResult
	for _, cmd := range commands {
		ctx := context.Background()
		var cancel context.CancelFunc
		if perCmdTimeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), perCmdTimeout)
		}
		res := host.ExecuteStream(ctx, cmd, onEvent)
		if cancel != nil {
			cancel()
		}
		results = append(results, res)
	}
	return results, nil
}

// BatchExecuteCommandsStream 实现方案 B：
// - GlobalTimeout 仅控制是否继续启动新的主机任务（到期后不再启动新的），
// - 已经启动的主机不会被中断，按各自 perCmdTimeout 执行
// - onEvent 回调负责接收所有实时输出（来自任意主机/任意命令）
func BatchExecuteCommandsStream(hosts []*HostInfo, commands []string, cfg *BatchConfig, onEvent func(StreamEvent)) (map[string][]*ExecutorResult, error) {
	if cfg == nil {
		cfg = &BatchConfig{MaxConcurrent: 10, GlobalTimeout: 0, PerCmdTimeout: 30}
	}

	results := make(map[string][]*ExecutorResult)
	var mu sync.Mutex

	sem := make(chan struct{}, cfg.MaxConcurrent)
	// 全局调度超时，仅用于决定是否继续启动新的主机
	schedCtx := context.Background()
	var schedCancel context.CancelFunc
	if cfg.GlobalTimeout > 0 {
		schedCtx, schedCancel = context.WithTimeout(context.Background(), time.Duration(cfg.GlobalTimeout)*time.Second)
		defer schedCancel()
	}

	var wg sync.WaitGroup

	for _, h := range hosts {

		wg.Add(1)
		host := h
		go func() {
			defer wg.Done()

			// 并发控制
			sem <- struct{}{}
			defer func() { <-sem }()

			select {
			case <-schedCtx.Done():
				return
			default:
			}

			// 每台主机使用独立的 per-host 执行，不受 schedCtx 影响
			perCmdTimeout := time.Duration(cfg.PerCmdTimeout) * time.Second
			resSlice, _ := host.ExecuteCommandsStream(commands, perCmdTimeout, onEvent)

			mu.Lock()
			results[host.Address] = resSlice
			mu.Unlock()
		}()
	}

	wg.Wait()
	if schedCtx.Err() != nil {
		return nil, fmt.Errorf("全局调度超时")
	}
	return results, nil
}
