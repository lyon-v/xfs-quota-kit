package commands

import (
	"context"

	"github.com/xfs-quota-kit/pkg/config"
)

type contextKey string

const (
	configKey contextKey = "config"
)

// WithConfig 将配置添加到context
func WithConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// GetConfig 从context获取配置
func GetConfig(ctx context.Context) *config.Config {
	if cfg, ok := ctx.Value(configKey).(*config.Config); ok {
		return cfg
	}
	return nil
}
