package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type GlobalConfig struct {
	DataBase *DatabaseDO
	Server   *ServerConfig
	Jwt      *JwtProperties
	Redis    *RedisProperties
}

// DatabaseDO 数据库配置
type DatabaseDO struct {
	Dsn  string
	Opts *gorm.Config
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port  string
	Debug bool
}

// RedisProperties Redis 配置
type RedisProperties struct {
	Host     string
	Port     string
	Password string
	Username string
	DB       int
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
	v.SetConfigName("config")
	// 读取配置

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read file error : %v\n", err)
	}

	if err := v.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("unmarshal error : %v\n", err)
	}
	return nil
}

// GetGlobalConfig 提供给外部获取全局配置
func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}
