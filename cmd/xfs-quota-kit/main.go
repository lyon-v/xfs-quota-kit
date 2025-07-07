package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xfs-quota-kit/cmd/xfs-quota-kit/commands"
	"github.com/xfs-quota-kit/pkg/config"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "xfs-quota-kit",
		Short: "Advanced XFS quota management toolkit",
		Long: `XFS Quota Kit is a comprehensive command-line tool for managing XFS filesystem quotas.
It provides advanced features for user, group, and project quota management,
monitoring, reporting, and automation.

Features:
  - User, Group, and Project quota management
  - Batch operations and automation
  - Real-time monitoring and alerting
  - Comprehensive reporting
  - REST API server
  - Configuration file support`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 加载配置文件
			cfg, err := config.Load(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// 验证配置
			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			// 将配置添加到context
			ctx := commands.WithConfig(cmd.Context(), cfg)
			cmd.SetContext(ctx)

			return nil
		},
	}

	// 全局标志
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")

	// 添加子命令
	cmd.AddCommand(
		commands.NewQuotaCommand(),
		commands.NewProjectCommand(),
		commands.NewReportCommand(),
		commands.NewMonitorCommand(),
		commands.NewServerCommand(),
		commands.NewCompletionCommand(),
		newVersionCommand(),
	)

	return cmd
}

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("XFS Quota Kit %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Built: %s\n", date)
		},
	}

	return cmd
}
