package util

import (
	"errors"
	"fmt"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate 解析参数并将参数绑定到obj
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var errs validator.ValidationErrors
		if ok := errors.As(err, &errs); ok {
			log.Printf("参数校验失败: %v", errs)
			common.Fail(c, common.NewErrorCode(400, errs[0].Translate(common.GetTranslator())))
			c.Abort()
			return false
		}
		log.Printf("解析参数失败: %v", err)
		common.Fail(c, common.ServerError)
		return false
	}
	return true
}

func BindQueryParam(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		var errs validator.ValidationErrors
		if ok := errors.As(err, &errs); ok {
			log.Printf("参数校验失败: %v", errs)
			common.FailWithMsg(c, errs[0].Translate(trans))
			return false
		}
		log.Printf("解析参数失败: %v", err)
		common.Fail(c, common.ServerError)
		return false
	}
	return true
}

// GetParam 获取路径参数以及参数校验
func GetParam(c *gin.Context, key string, param interface{}, validate func(param interface{})) {
	var value string
	value = c.Query(key)
	if value == "" {
		value = c.Param(key)
	}
	if strParam, ok := param.(*string); ok {
		*strParam = value
	}
	if int64Param, ok := param.(*int64); ok {
		*int64Param, _ = strconv.ParseInt(value, 10, 64)
	}

	if validate != nil {
		validate(param)
	}
	return
}
func BeautifyRawData(rawData map[string]string) (*dto.BeautifiedMetrics, error) {
	metrics := &dto.BeautifiedMetrics{
		Status:    "success",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 解析主机名
	metrics.Hostname = strings.TrimSpace(rawData["主机名："])

	// 解析系统信息
	if err := parseSystemInfo(metrics, rawData); err != nil {
		return nil, err
	}

	// 解析CPU信息
	if err := parseCPUInfo(metrics, rawData); err != nil {
		return nil, err
	}

	// 解析内存信息
	if err := parseMemoryInfo(metrics, rawData); err != nil {
		return nil, err
	}

	// 解析磁盘信息
	if err := parseDiskInfo(metrics, rawData); err != nil {
		return nil, err
	}

	// 解析网络信息
	if err := parseNetworkInfo(metrics, rawData); err != nil {
		return nil, err
	}

	// 解析进程信息
	if err := parseProcessInfo(metrics, rawData); err != nil {
		return nil, err
	}

	// 生成总结信息
	generateSummary(metrics)

	return metrics, nil
}

func parseSystemInfo(metrics *dto.BeautifiedMetrics, raw map[string]string) error {
	// 解析操作系统信息
	if osInfo, ok := raw["操作系统信息："]; ok {
		lines := strings.Split(osInfo, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				metrics.System.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			} else if strings.HasPrefix(line, "NAME=") && metrics.System.OS == "" {
				name := strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
				version := ""
				if versionLine := findLine(lines, "VERSION="); versionLine != "" {
					version = strings.Trim(strings.TrimPrefix(versionLine, "VERSION="), "\"")
				}
				metrics.System.OS = name + " " + version
			}
		}
	}

	// 解析内核信息
	if kernelInfo, ok := raw["内核信息："]; ok {
		metrics.System.Kernel = strings.TrimSpace(kernelInfo)
		// 提取架构
		if strings.Contains(kernelInfo, "x86_64") {
			metrics.System.Architecture = "x86_64"
		} else if strings.Contains(kernelInfo, "aarch64") {
			metrics.System.Architecture = "arm64"
		} else if strings.Contains(kernelInfo, "arm") {
			metrics.System.Architecture = "arm"
		} else {
			metrics.System.Architecture = "unknown"
		}
	}

	// 解析uptime
	if uptime, ok := raw["uptime"]; ok {
		metrics.System.Uptime = parseUptime(uptime)
		metrics.System.Users = parseUserCount(uptime)
	}

	return nil
}

func parseCPUInfo(metrics *dto.BeautifiedMetrics, raw map[string]string) error {
	if cpuInfo, ok := raw["cpu信息："]; ok {
		lines := strings.Split(cpuInfo, "\n")
		processorCount := 0
		for _, line := range lines {
			if strings.HasPrefix(line, "processor") {
				processorCount++
			}
			if strings.HasPrefix(line, "model name") && metrics.CPU.Model == "" {
				metrics.CPU.Model = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			}
			if strings.HasPrefix(line, "cpu MHz") && metrics.CPU.Frequency == "" {
				freq := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				if f, err := strconv.ParseFloat(freq, 64); err == nil {
					metrics.CPU.Frequency = fmt.Sprintf("%.2f GHz", f/1000)
				}
			}
			if strings.HasPrefix(line, "cache size") && metrics.CPU.Cache == "" {
				metrics.CPU.Cache = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			}
			if strings.HasPrefix(line, "cpu cores") {
				if cores, err := strconv.Atoi(strings.TrimSpace(strings.SplitN(line, ":", 2)[1])); err == nil {
					metrics.CPU.Cores = cores
				}
			}
		}
		if metrics.CPU.Cores == 0 {
			metrics.CPU.Cores = processorCount
		}
		metrics.CPU.Threads = processorCount
	}

	// 解析负载信息
	if loadInfo, ok := raw["负载信息："]; ok {
		parts := strings.Fields(loadInfo)
		if len(parts) >= 3 {
			metrics.CPU.LoadAverage.Min1, _ = strconv.ParseFloat(parts[0], 64)
			metrics.CPU.LoadAverage.Min5, _ = strconv.ParseFloat(parts[1], 64)
			metrics.CPU.LoadAverage.Min15, _ = strconv.ParseFloat(parts[2], 64)
		}
	}

	// 解析CPU使用率（从top命令）
	if topInfo, ok := raw["top"]; ok {
		lines := strings.Split(topInfo, "\n")
		for _, line := range lines {
			if strings.Contains(line, "%Cpu(s):") {
				parseCPUUsage(line, &metrics.CPU.Usage)
				break
			}
		}
	}

	return nil
}

func parseCPUUsage(line string, usage *dto.CPUUsage) {
	re := regexp.MustCompile(`(\d+\.\d+)\s+us,\s+(\d+\.\d+)\s+sy,\s+(\d+\.\d+)\s+ni,\s+(\d+\.\d+)\s+id,\s+(\d+\.\d+)\s+wa`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 6 {
		usage.User, _ = strconv.ParseFloat(matches[1], 64)
		usage.System, _ = strconv.ParseFloat(matches[2], 64)
		usage.Idle, _ = strconv.ParseFloat(matches[4], 64)
		usage.Wait, _ = strconv.ParseFloat(matches[5], 64)
	}
}

func parseMemoryInfo(metrics *dto.BeautifiedMetrics, raw map[string]string) error {
	if memInfo, ok := raw["内存："]; ok {
		lines := strings.Split(memInfo, "\n")
		if len(lines) >= 2 {
			memFields := strings.Fields(lines[1])
			if len(memFields) >= 6 {
				totalMB, _ := strconv.ParseUint(memFields[1], 10, 64)
				usedMB, _ := strconv.ParseUint(memFields[2], 10, 64)
				freeMB, _ := strconv.ParseUint(memFields[3], 10, 64)
				availableMB, _ := strconv.ParseUint(memFields[6], 10, 64)

				metrics.Memory.Total = formatBytes(totalMB * 1024 * 1024)
				metrics.Memory.Used = formatBytes(usedMB * 1024 * 1024)
				metrics.Memory.Free = formatBytes(freeMB * 1024 * 1024)
				metrics.Memory.Available = formatBytes(availableMB * 1024 * 1024)
				metrics.Memory.UsagePercent = float64(usedMB) / float64(totalMB) * 100
			}

			// 解析swap行: "Swap:             0           0           0"
			if len(lines) >= 3 {
				swapFields := strings.Fields(lines[2])
				if len(swapFields) >= 4 {
					swapTotal, _ := strconv.ParseUint(swapFields[1], 10, 64)
					swapUsed, _ := strconv.ParseUint(swapFields[2], 10, 64)
					swapFree, _ := strconv.ParseUint(swapFields[3], 10, 64)

					metrics.Memory.SwapTotal = formatBytes(swapTotal * 1024 * 1024)
					metrics.Memory.SwapUsed = formatBytes(swapUsed * 1024 * 1024)
					metrics.Memory.SwapFree = formatBytes(swapFree * 1024 * 1024)
				}
			}
		}
	}
	return nil
}

func parseDiskInfo(metrics *dto.BeautifiedMetrics, raw map[string]string) error {
	if diskInfo, ok := raw["磁盘信息："]; ok {
		lines := strings.Split(diskInfo, "\n")
		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue // 跳过标题行和空行
			}

			// 示例: "/dev/sda3        96G  8.2G   88G    9% /"
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				disk := dto.DiskInfo{
					Filesystem: fields[0],
					Size:       fields[1],
					Used:       fields[2],
					Available:  fields[3],
					MountedOn:  fields[5],
				}

				// 解析使用百分比
				if percentStr := strings.TrimSuffix(fields[4], "%"); percentStr != "" {
					disk.UsagePercent, _ = strconv.Atoi(percentStr)
				}

				metrics.Disk = append(metrics.Disk, disk)
			}
		}
	}
	return nil
}

func parseNetworkInfo(metrics *dto.BeautifiedMetrics, raw map[string]string) error {
	// 解析网络接口
	if netIf, ok := raw["net_if"]; ok {
		lines := strings.Split(netIf, "\n")
		currentInterface := ""

		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}

			interfaceName := strings.TrimSuffix(fields[1], ":")
			if interfaceName != currentInterface {
				netIf := dto.NetworkInterface{
					Name: interfaceName,
					Type: getInterfaceType(interfaceName),
				}
				metrics.Network.Interfaces = append(metrics.Network.Interfaces, netIf)
				currentInterface = interfaceName
			}

			// 添加IP地址
			idx := len(metrics.Network.Interfaces) - 1
			if strings.Contains(fields[2], "inet6") {
				metrics.Network.Interfaces[idx].IPv6 = append(metrics.Network.Interfaces[idx].IPv6, fields[3])
			} else if strings.Contains(fields[2], "inet") {
				metrics.Network.Interfaces[idx].IPv4 = append(metrics.Network.Interfaces[idx].IPv4, fields[3])
			}
		}
	}

	// 解析网络连接统计
	if netstat, ok := raw["netstat"]; ok {
		lines := strings.Split(netstat, "\n")
		for _, line := range lines {
			if strings.Contains(line, "LISTEN") {
				metrics.Network.Connections.Listening++
			} else if strings.Contains(line, "ESTAB") {
				metrics.Network.Connections.Established++
			}
		}
		metrics.Network.Connections.Total = len(lines) - 1 // 减去标题行
	}

	return nil
}

func parseProcessInfo(metrics *dto.BeautifiedMetrics, raw map[string]string) error {
	// 解析进程总数
	if procCount, ok := raw["proc_count"]; ok {
		metrics.Processes.Total, _ = strconv.Atoi(strings.TrimSpace(procCount))
	}

	// 解析top命令中的进程状态
	if topInfo, ok := raw["top"]; ok {
		lines := strings.Split(topInfo, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Tasks:") {
				// 示例: "Tasks: 163 total,   7 running, 156 sleeping,   0 stopped,   0 zombie"
				re := regexp.MustCompile(`(\d+)\s+total,\s+(\d+)\s+running,\s+(\d+)\s+sleeping,\s+(\d+)\s+stopped,\s+(\d+)\s+zombie`)
				matches := re.FindStringSubmatch(line)
				if len(matches) == 6 {
					metrics.Processes.Running, _ = strconv.Atoi(matches[2])
					metrics.Processes.Sleeping, _ = strconv.Atoi(matches[3])
					metrics.Processes.Stopped, _ = strconv.Atoi(matches[4])
					metrics.Processes.Zombie, _ = strconv.Atoi(matches[5])
				}
				break
			}
		}
	}

	// 解析进程列表
	if procInfo, ok := raw["进程信息："]; ok {
		lines := strings.Split(procInfo, "\n")
		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}

			fields := strings.Fields(line)
			n := len(fields)
			if n >= 6 {
				process := dto.Process{
					PID:   parseInt(fields[0]),
					Name:  fields[6],
					User:  fields[2],
					State: fields[5],
				}
				process.CPUPercent, _ = strconv.ParseFloat(fields[3], 64)
				process.MemoryPercent, _ = strconv.ParseFloat(fields[4], 64)

				metrics.Processes.TopProcesses = append(metrics.Processes.TopProcesses, process)

				// 只取前10个进程
				if len(metrics.Processes.TopProcesses) >= 10 {
					break
				}
			}
		}
	}

	return nil
}

// ==================== 辅助函数 =======================
func parseUptime(uptime string) string {
	re := regexp.MustCompile(`up\s+([^,]+),`)
	matches := re.FindStringSubmatch(uptime)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return "unknown"
}

func parseUserCount(uptime string) int {
	re := regexp.MustCompile(`(\d+)\s+user`)
	matches := re.FindStringSubmatch(uptime)
	if len(matches) > 1 {
		count, _ := strconv.Atoi(matches[1])
		return count
	}
	return 0
}

func getInterfaceType(name string) string {
	if name == "lo" {
		return "loopback"
	} else if strings.HasPrefix(name, "eth") || strings.HasPrefix(name, "ens") {
		return "ethernet"
	} else if strings.HasPrefix(name, "wlan") || strings.HasPrefix(name, "wlp") {
		return "wireless"
	} else if strings.HasPrefix(name, "docker") || strings.HasPrefix(name, "br-") {
		return "bridge"
	} else {
		return "unknown"
	}
}

func findLine(lines []string, prefix string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// 判断是否报警并获取报警信息
func generateSummary(metrics *dto.BeautifiedMetrics) {
	var alerts []string
	var recommendations []string

	if metrics.Memory.UsagePercent > 85 {
		alerts = append(alerts, "内存使用率过高")
		recommendations = append(recommendations, "建议检查内存使用情况")
	}

	if metrics.CPU.LoadAverage.Min1 > float64(metrics.CPU.Cores) {
		alerts = append(alerts, "系统负载较高")
		recommendations = append(recommendations, "建议检查CPU密集型进程")
	}

	for _, disk := range metrics.Disk {
		if disk.UsagePercent > 90 {
			alerts = append(alerts, fmt.Sprintf("磁盘 %s 使用率过高", disk.Filesystem))
			recommendations = append(recommendations, "建议清理磁盘空间")
		}
	}

	if len(alerts) == 0 {
		metrics.Summary.Health = "healthy"
		metrics.Summary.Recommendations = []string{"系统运行正常"}
	} else {
		metrics.Summary.Health = "warning"
		metrics.Summary.Alerts = alerts
		metrics.Summary.Recommendations = recommendations
	}
}
