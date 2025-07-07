package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	XFS      XFSConfig      `mapstructure:"xfs"`
	Monitor  MonitorConfig  `mapstructure:"monitor"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string    `mapstructure:"host"`
	Port int       `mapstructure:"port"`
	Mode string    `mapstructure:"mode"` // debug, release, test
	TLS  TLSConfig `mapstructure:"tls"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // sqlite, mysql, postgres
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level"`    // debug, info, warn, error
	Format     string `mapstructure:"format"`   // json, text
	Output     string `mapstructure:"output"`   // stdout, file
	File       string `mapstructure:"file"`     // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"` // MB
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"` // days
	Compress   bool   `mapstructure:"compress"`
}

// XFSConfig XFS配置
type XFSConfig struct {
	DefaultPath   string           `mapstructure:"default_path"`
	ProjectsFile  string           `mapstructure:"projects_file"`
	ProjidFile    string           `mapstructure:"projid_file"`
	DefaultLimits DefaultLimits    `mapstructure:"default_limits"`
	AutoCreate    bool             `mapstructure:"auto_create"`
	BackupEnabled bool             `mapstructure:"backup_enabled"`
	BackupPath    string           `mapstructure:"backup_path"`
	Filesystems   []FilesystemInfo `mapstructure:"filesystems"`
}

// DefaultLimits 默认配额限制
type DefaultLimits struct {
	UserBlockSoft  string `mapstructure:"user_block_soft"` // e.g., "1GB"
	UserBlockHard  string `mapstructure:"user_block_hard"` // e.g., "2GB"
	UserInodeSoft  uint64 `mapstructure:"user_inode_soft"`
	UserInodeHard  uint64 `mapstructure:"user_inode_hard"`
	GroupBlockSoft string `mapstructure:"group_block_soft"`
	GroupBlockHard string `mapstructure:"group_block_hard"`
	GroupInodeSoft uint64 `mapstructure:"group_inode_soft"`
	GroupInodeHard uint64 `mapstructure:"group_inode_hard"`
}

// FilesystemInfo 文件系统信息
type FilesystemInfo struct {
	Name       string `mapstructure:"name"`
	MountPoint string `mapstructure:"mount_point"`
	Device     string `mapstructure:"device"`
	Options    string `mapstructure:"options"`
	Enabled    bool   `mapstructure:"enabled"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	Interval          string `mapstructure:"interval"`        // e.g., "5m"
	AlertThreshold    int    `mapstructure:"alert_threshold"` // 使用率百分比
	ReportPath        string `mapstructure:"report_path"`
	ReportInterval    string `mapstructure:"report_interval"` // e.g., "1h"
	EmailNotification bool   `mapstructure:"email_notification"`
	WebhookURL        string `mapstructure:"webhook_url"`
}

// Load 加载配置
func Load(configFile string) (*Config, error) {
	config := &Config{}

	// 设置默认值
	setDefaults()

	// 设置配置文件
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("/etc/xfs-quota-kit")
		viper.AddConfigPath("$HOME/.xfs-quota-kit")
		viper.AddConfigPath(".")
	}

	// 支持环境变量
	viper.SetEnvPrefix("XFS_QUOTA")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认配置
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 解析到结构体
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("server.tls.enabled", false)

	// 数据库默认配置
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.database", "xfs_quota.db")

	// 日志默认配置
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_backups", 3)
	viper.SetDefault("logging.max_age", 28)
	viper.SetDefault("logging.compress", true)

	// XFS默认配置
	viper.SetDefault("xfs.default_path", "/mnt/xfs")
	viper.SetDefault("xfs.projects_file", "/etc/projects")
	viper.SetDefault("xfs.projid_file", "/etc/projid")
	viper.SetDefault("xfs.auto_create", true)
	viper.SetDefault("xfs.backup_enabled", true)
	viper.SetDefault("xfs.backup_path", "/var/backups/xfs-quota-kit")

	// 默认配额限制
	viper.SetDefault("xfs.default_limits.user_block_soft", "1GB")
	viper.SetDefault("xfs.default_limits.user_block_hard", "2GB")
	viper.SetDefault("xfs.default_limits.user_inode_soft", 100000)
	viper.SetDefault("xfs.default_limits.user_inode_hard", 200000)
	viper.SetDefault("xfs.default_limits.group_block_soft", "10GB")
	viper.SetDefault("xfs.default_limits.group_block_hard", "20GB")
	viper.SetDefault("xfs.default_limits.group_inode_soft", 1000000)
	viper.SetDefault("xfs.default_limits.group_inode_hard", 2000000)

	// 监控默认配置
	viper.SetDefault("monitor.enabled", true)
	viper.SetDefault("monitor.interval", "5m")
	viper.SetDefault("monitor.alert_threshold", 80)
	viper.SetDefault("monitor.report_path", "/var/log/xfs-quota-kit/reports")
	viper.SetDefault("monitor.report_interval", "1h")
	viper.SetDefault("monitor.email_notification", false)
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.Mode != "debug" && c.Server.Mode != "release" && c.Server.Mode != "test" {
		return fmt.Errorf("invalid server mode: %s", c.Server.Mode)
	}

	// 验证TLS配置
	if c.Server.TLS.Enabled {
		if c.Server.TLS.CertFile == "" || c.Server.TLS.KeyFile == "" {
			return fmt.Errorf("TLS enabled but cert_file or key_file not specified")
		}
	}

	// 验证日志配置
	validLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLevels, c.Logging.Level) {
		return fmt.Errorf("invalid logging level: %s", c.Logging.Level)
	}

	validFormats := []string{"json", "text"}
	if !contains(validFormats, c.Logging.Format) {
		return fmt.Errorf("invalid logging format: %s", c.Logging.Format)
	}

	validOutputs := []string{"stdout", "file"}
	if !contains(validOutputs, c.Logging.Output) {
		return fmt.Errorf("invalid logging output: %s", c.Logging.Output)
	}

	if c.Logging.Output == "file" && c.Logging.File == "" {
		return fmt.Errorf("logging output set to file but no file specified")
	}

	return nil
}

// GetAddress 获取服务器地址
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDebugMode 判断是否为调试模式
func (c *Config) IsDebugMode() bool {
	return c.Server.Mode == "debug"
}

// 辅助函数
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
