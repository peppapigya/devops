package dto

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// HostPageRequest 主机分页请求参数
type HostPageRequest struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
}

// CreateHostDTO 创建主机请求参数
type CreateHostDTO struct {
	HostName     string `json:"hostName" validate:"required"`
	Address      string `json:"address" validate:"required"`
	HostPort     int32  `json:"hostPort" validate:"required"`
	Username     string `json:"username" validate:"required"`
	HostPassword string `json:"hostPassword" validate:"required"`
	Remark       string `json:"remark"`
}

// UpdateHostDTO 更新主机请求参数
type UpdateHostDTO struct {
	ID           uint32 `json:"id" validate:"required"`
	HostName     string `json:"hostName" validate:"required"`
	Address      string `json:"address" validate:"required"`
	HostPort     int32  `json:"hostPort" validate:"required"`
	Username     string `json:"username" validate:"required"`
	HostPassword string `json:"hostPassword"`
	Remark       string `json:"remark"`
}

// DiskUsage 磁盘使用情况统计
type DiskUsage struct {
	Device      string  `json:"device" common:"设备"`
	Path        string  `json:"path" common:"挂载点"`
	Fstype      string  `json:"fstype" common:"文件系统类型"`
	Total       uint64  `json:"total" common:"总容量"`
	Free        uint64  `json:"free" common:"空闲空间"`
	Used        uint64  `json:"used" common:"已用空间"`
	UsedPercent float64 `json:"usedPercent" common:"使用率"`
}

// DiskIOCounters 磁盘IO计数器统计
type DiskIOCounters struct {
	Name       string `json:"name" common:"磁盘名称"`
	ReadCount  uint64 `json:"readCount" common:"读取次数"`
	WriteCount uint64 `json:"writeCount" common:"写入次数"`
	ReadBytes  uint64 `json:"readBytes" common:"读取字节数"`
	WriteBytes uint64 `json:"writeBytes" common:"写入字节数"`
	ReadTime   uint64 `json:"readTime" common:"读取时间"`
	WriteTime  uint64 `json:"writeTime" common:"写入时间"`
	IoTime     uint64 `json:"ioTime" common:"IO时间"`
	WeightedIO uint64 `json:"weightedIo" common:"加权IO"`
}

// NetInterface 网络接口信息
type NetInterface struct {
	Name         string   `json:"name" common:"接口名称"`
	MTU          int      `json:"mtu" common:"最大传输单元"`
	HardwareAddr string   `json:"hardwareAddr" common:"MAC地址"`
	Flags        []string `json:"flags" common:"接口标志"`
	Addresses    []string `json:"addresses" common:"IP地址列表"`
}

// NetIOCounters 网络IO计数器统计
type NetIOCounters struct {
	Name        string `json:"name" common:"接口名称"`
	BytesSent   uint64 `json:"bytesSent" common:"发送字节数"`
	BytesRecv   uint64 `json:"bytesRecv" common:"接收字节数"`
	PacketsSent uint64 `json:"packetsSent" common:"发送包数"`
	PacketsRecv uint64 `json:"packetsRecv" common:"接收包数"`
	Errin       uint64 `json:"errin" common:"接收错误数"`
	Errout      uint64 `json:"errout" common:"发送错误数"`
	Dropin      uint64 `json:"dropin" common:"接收丢包数"`
	Dropout     uint64 `json:"dropout" common:"发送丢包数"`
}

// ProcessStat 进程状态信息
type ProcessStat struct {
	Pid           int32   `json:"pid" common:"进程ID"`
	Name          string  `json:"name" common:"进程名称"`
	Username      string  `json:"username" common:"用户名"`
	CPUPercent    float64 `json:"cpuPercent" common:"CPU使用率"`
	MemoryPercent float32 `json:"memoryPercent" common:"内存使用率"`
	MemoryRSS     uint64  `json:"memoryRss" common:"物理内存使用量"`
	CommandLine   string  `json:"cmdline" common:"命令行"`
}

// SystemMetrics 系统指标综合结构体
type SystemMetrics struct {
	Host            *host.InfoStat  `json:"host" common:"主机信息"`
	OS              string          `json:"os" common:"操作系统"`
	Platform        string          `json:"platform" common:"平台"`
	PlatformVersion string          `json:"platformVersion" common:"平台版本"`
	Virtualization  string          `json:"virtualization" common:"虚拟化"`
	Users           []host.UserStat `json:"users" common:"在线用户"`

	LoadAvg  *load.AvgStat  `json:"loadAvg" common:"平均负载"`
	LoadMisc *load.MiscStat `json:"loadMisc,omitempty" common:"负载杂项"`

	CPUInfo          []cpu.InfoStat  `json:"cpuInfo" common:"CPU信息"`
	CPUPercentTotal  float64         `json:"cpuPercentTotal" common:"总CPU使用率"`
	CPUPercentPerCPU []float64       `json:"cpuPercentPerCpu" common:"单核CPU使用率"`
	CPUTimes         []cpu.TimesStat `json:"cpuTimes" common:"CPU时间"`

	Memory *mem.VirtualMemoryStat `json:"memory" common:"内存信息"`
	Swap   *mem.SwapMemoryStat    `json:"swap" common:"交换分区"`

	// 磁盘信息
	Disks  []DiskUsage      `json:"disks" common:"磁盘使用情况"`
	DiskIO []DiskIOCounters `json:"diskIo" common:"磁盘IO"`

	// 网络信息
	NetInterfaces []NetInterface  `json:"netInterfaces" common:"网络接口"`
	NetIOPerNIC   []NetIOCounters `json:"netIoPerNic" common:"网络接口IO"`
	NetIOTotal    NetIOCounters   `json:"netIoTotal" common:"网络IO总计"`

	// 进程信息
	ProcessCount int           `json:"processCount" common:"进程总数"`
	TopProcesses []ProcessStat `json:"topProcesses" common:"重点进程"`

	// 时间戳
	Timestamp time.Time `json:"timestamp" common:"采集时间"`
}

// ssh拿到的数据美化
type BeautifiedMetrics struct {
	Status    string      `json:"status"`
	Timestamp string      `json:"timestamp"`
	Hostname  string      `json:"hostname"`
	System    SystemInfo  `json:"system"`
	CPU       CPUInfo     `json:"cpu"`
	Memory    MemoryInfo  `json:"memory"`
	Disk      []DiskInfo  `json:"disk"`
	Network   NetworkInfo `json:"network"`
	Processes ProcessInfo `json:"processes"`
	Summary   SummaryInfo `json:"summary"`
}

type SystemInfo struct {
	OS           string `json:"os"`
	Kernel       string `json:"kernel"`
	Architecture string `json:"architecture"`
	Uptime       string `json:"uptime"`
	Users        int    `json:"users"`
}

type CPUInfo struct {
	Model       string   `json:"model"`
	Cores       int      `json:"cores"`
	Threads     int      `json:"threads"`
	Frequency   string   `json:"frequency"`
	Cache       string   `json:"cache"`
	LoadAverage LoadAvg  `json:"load_average"`
	Usage       CPUUsage `json:"usage"`
}

type LoadAvg struct {
	Min1  float64 `json:"1min"`
	Min5  float64 `json:"5min"`
	Min15 float64 `json:"15min"`
}

type CPUUsage struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
	Wait   float64 `json:"wait"`
}

type MemoryInfo struct {
	Total        string  `json:"total"`
	Used         string  `json:"used"`
	Free         string  `json:"free"`
	Available    string  `json:"available"`
	UsagePercent float64 `json:"usage_percent"`
	SwapTotal    string  `json:"swap_total"`
	SwapUsed     string  `json:"swap_used"`
	SwapFree     string  `json:"swap_free"`
}

type DiskInfo struct {
	Filesystem   string `json:"filesystem"`
	Size         string `json:"size"`
	Used         string `json:"used"`
	Available    string `json:"available"`
	UsagePercent int    `json:"usage_percent"`
	MountedOn    string `json:"mounted_on"`
}

type NetworkInfo struct {
	Interfaces  []NetworkInterface `json:"interfaces"`
	Connections ConnectionInfo     `json:"connections"`
}

type NetworkInterface struct {
	Name string   `json:"name"`
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
	Type string   `json:"type"`
}

type ConnectionInfo struct {
	Total       int `json:"total"`
	Listening   int `json:"listening"`
	Established int `json:"established"`
}

type ProcessInfo struct {
	Total        int       `json:"total"`
	Running      int       `json:"running"`
	Sleeping     int       `json:"sleeping"`
	Stopped      int       `json:"stopped"`
	Zombie       int       `json:"zombie"`
	TopProcesses []Process `json:"top_processes"`
}

type Process struct {
	PID           int     `json:"pid"`
	Name          string  `json:"name"`
	User          string  `json:"user"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryUsage   string  `json:"memory_usage"`
	State         string  `json:"state"`
}

type SummaryInfo struct {
	Health          string   `json:"health"`
	Alerts          []string `json:"alerts"`
	Recommendations []string `json:"recommendations"`
}
