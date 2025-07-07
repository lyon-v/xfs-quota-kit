package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Host: "0.0.0.0",
					Port: 8080,
					Mode: "release",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
					Output: "stdout",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			config: Config{
				Server: ServerConfig{
					Port: 0,
					Mode: "release",
				},
			},
			wantErr: true,
			errMsg:  "invalid server port",
		},
		{
			name: "invalid port - too high",
			config: Config{
				Server: ServerConfig{
					Port: 70000,
					Mode: "release",
				},
			},
			wantErr: true,
			errMsg:  "invalid server port",
		},
		{
			name: "invalid server mode",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid server mode",
		},
		{
			name: "TLS enabled but no cert file",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "release",
					TLS: TLSConfig{
						Enabled:  true,
						CertFile: "",
						KeyFile:  "key.pem",
					},
				},
			},
			wantErr: true,
			errMsg:  "TLS enabled but cert_file or key_file not specified",
		},
		{
			name: "invalid logging level",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "release",
				},
				Logging: LoggingConfig{
					Level: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid logging level",
		},
		{
			name: "invalid logging format",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "release",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid logging format",
		},
		{
			name: "invalid logging output",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "release",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
					Output: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid logging output",
		},
		{
			name: "file output but no file specified",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "release",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "json",
					Output: "file",
					File:   "",
				},
			},
			wantErr: true,
			errMsg:  "logging output set to file but no file specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_GetAddress(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 9090,
		},
	}

	assert.Equal(t, "127.0.0.1:9090", config.GetAddress())
}

func TestConfig_IsDebugMode(t *testing.T) {
	tests := []struct {
		mode     string
		expected bool
	}{
		{"debug", true},
		{"release", false},
		{"test", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			config := &Config{
				Server: ServerConfig{
					Mode: tt.mode,
				},
			}
			assert.Equal(t, tt.expected, config.IsDebugMode())
		})
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")

	configContent := `
server:
  host: "127.0.0.1"
  port: 9000
  mode: "debug"

logging:
  level: "debug"
  format: "text"
  output: "stdout"

xfs:
  default_path: "/tmp/xfs-test"
  projects_file: "/tmp/projects"
  projid_file: "/tmp/projid"

monitor:
  enabled: false
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 加载配置
	config, err := Load(configFile)
	require.NoError(t, err)

	// 验证配置
	assert.Equal(t, "127.0.0.1", config.Server.Host)
	assert.Equal(t, 9000, config.Server.Port)
	assert.Equal(t, "debug", config.Server.Mode)
	assert.Equal(t, "debug", config.Logging.Level)
	assert.Equal(t, "text", config.Logging.Format)
	assert.Equal(t, "/tmp/xfs-test", config.XFS.DefaultPath)
	assert.False(t, config.Monitor.Enabled)
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// 测试加载不存在的配置文件（应该使用默认值）
	config, err := Load("/non/existent/config.yaml")
	require.NoError(t, err)

	// 验证默认值
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, "release", config.Server.Mode)
	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "json", config.Logging.Format)
	assert.Equal(t, "stdout", config.Logging.Output)
	assert.Equal(t, "/mnt/xfs", config.XFS.DefaultPath)
	assert.True(t, config.Monitor.Enabled)
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// 设置环境变量
	oldEnv := os.Getenv("XFS_QUOTA_SERVER_PORT")
	defer func() {
		if oldEnv != "" {
			os.Setenv("XFS_QUOTA_SERVER_PORT", oldEnv)
		} else {
			os.Unsetenv("XFS_QUOTA_SERVER_PORT")
		}
	}()

	os.Setenv("XFS_QUOTA_SERVER_PORT", "9999")

	config, err := Load("")
	require.NoError(t, err)

	// 环境变量应该覆盖默认值
	assert.Equal(t, 9999, config.Server.Port)
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	// 创建无效的YAML文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-config.yaml")

	invalidYAML := `
server:
  host: "127.0.0.1"
  port: invalid_port
  invalid_yaml: [
`

	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	// 加载配置应该失败
	_, err = Load(configFile)
	assert.Error(t, err)
}

func TestContainsHelper(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	assert.True(t, contains(slice, "apple"))
	assert.True(t, contains(slice, "banana"))
	assert.True(t, contains(slice, "cherry"))
	assert.False(t, contains(slice, "orange"))
	assert.False(t, contains(slice, ""))
	assert.False(t, contains([]string{}, "apple"))
}

func TestDefaultLimitsConfig(t *testing.T) {
	config, err := Load("")
	require.NoError(t, err)

	// 验证默认配额限制
	assert.Equal(t, "1GB", config.XFS.DefaultLimits.UserBlockSoft)
	assert.Equal(t, "2GB", config.XFS.DefaultLimits.UserBlockHard)
	assert.Equal(t, uint64(100000), config.XFS.DefaultLimits.UserInodeSoft)
	assert.Equal(t, uint64(200000), config.XFS.DefaultLimits.UserInodeHard)
	assert.Equal(t, "10GB", config.XFS.DefaultLimits.GroupBlockSoft)
	assert.Equal(t, "20GB", config.XFS.DefaultLimits.GroupBlockHard)
}

func TestMonitorConfig(t *testing.T) {
	config, err := Load("")
	require.NoError(t, err)

	// 验证监控配置默认值
	assert.True(t, config.Monitor.Enabled)
	assert.Equal(t, "5m", config.Monitor.Interval)
	assert.Equal(t, 80, config.Monitor.AlertThreshold)
	assert.Equal(t, "/var/log/xfs-quota-kit/reports", config.Monitor.ReportPath)
	assert.Equal(t, "1h", config.Monitor.ReportInterval)
	assert.False(t, config.Monitor.EmailNotification)
}
