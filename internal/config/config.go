package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Email    EmailConfig    `mapstructure:"email"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path    string `mapstructure:"path"`
	MaxSize int64  `mapstructure:"max_size"` // 字节
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
	File  string `mapstructure:"file"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Server.Mode != "debug" && c.Server.Mode != "release" {
		return fmt.Errorf("invalid server mode: %s (must be 'debug' or 'release')", c.Server.Mode)
	}

	// 验证MySQL配置
	if c.Database.MySQL.Host == "" {
		return fmt.Errorf("mysql host is required")
	}
	if c.Database.MySQL.User == "" {
		return fmt.Errorf("mysql user is required")
	}
	if c.Database.MySQL.DBName == "" {
		return fmt.Errorf("mysql dbname is required")
	}

	// 验证Redis配置
	if c.Database.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	// 验证JWT配置
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt secret is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("jwt secret must be at least 32 characters")
	}
	if c.JWT.ExpireHours <= 0 {
		return fmt.Errorf("jwt expire_hours must be positive")
	}

	// 验证上传配置
	if c.Upload.Path == "" {
		return fmt.Errorf("upload path is required")
	}
	if c.Upload.MaxSize <= 0 {
		return fmt.Errorf("upload max_size must be positive")
	}

	// 验证日志配置
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("invalid log level: %s (must be 'debug', 'info', 'warn', or 'error')", c.Log.Level)
	}

	return nil
}

// GetJWTExpireDuration 获取JWT过期时间
func (c *Config) GetJWTExpireDuration() time.Duration {
	return time.Duration(c.JWT.ExpireHours) * time.Hour
}

// GetMySQLDSN 获取MySQL连接字符串
func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.MySQL.User,
		c.Database.MySQL.Password,
		c.Database.MySQL.Host,
		c.Database.MySQL.Port,
		c.Database.MySQL.DBName,
	)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Database.Redis.Host, c.Database.Redis.Port)
}
