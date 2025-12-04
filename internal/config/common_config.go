package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type GlobalConfig struct {
	DataBase *DatabaseDO
	Server   *ServerConfig
	Jwt      *JwtProperties
	Redis    *RedisProperties
	Job      *JobConfig `yaml:"job"`
	Log      *LogProperties
}

// DatabaseDO 数据库配置
type DatabaseDO struct {
	Dsn  string
	Opts *gorm.Config
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port           string
	Debug          bool
	TrustedProxies []string `mapstructure:"trusted-proxies"`
}

// RedisProperties Redis 配置
type RedisProperties struct {
	Host     string
	Port     string
	Password string
	Username string
	DB       int
}

// LogProperties 日志配置
type LogProperties struct {
	Enable bool
	Path   string
	Level  string
	Stdout bool
}

// JwtProperties JWT 配置
type JwtProperties struct {
	Secret            string
	ExpireTime        int64    `mapstructure:"expire-time"`
	RefreshExpireTime int64    `mapstructure:"refresh-expire-time"`
	ExcludePaths      []string `mapstructure:"exclude-paths"` // 不需要进行参数校验的路径
}

var globalConfig *GlobalConfig

func LoadConfig() error {
	v := viper.New()
	v.AddConfigPath(".") // 添加文件搜索路径
	v.AddConfigPath("./configs")

	v.SetConfigType("yaml") // 匹配的文件后缀
	configName := getConfigNameByMode()
	v.SetConfigName(configName)
	// 读取配置

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read file error : %v\n", err)
	}

	if err := v.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("unmarshal global config error : %v\n", err)
	}
	return nil
}

// 根据启动模式去获取读取哪个配置文件
func getConfigNameByMode() string {
	switch gin.Mode() {
	case gin.DebugMode:
		return "config-dev"
	case gin.ReleaseMode:
		return "config-prod"
	case gin.TestMode:
		return "config-test"
	default:
		return "config"
	}
}

// GetGlobalConfig 提供给外部获取全局配置
func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}
