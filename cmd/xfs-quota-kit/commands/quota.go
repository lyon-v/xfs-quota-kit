package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xfs-quota-kit/pkg/xfs"
)

// NewQuotaCommand 创建配额管理命令
func NewQuotaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage XFS quotas",
		Long:  `Manage user, group, and project quotas on XFS filesystems.`,
	}

	cmd.AddCommand(
		newQuotaGetCommand(),
		newQuotaSetCommand(),
		newQuotaRemoveCommand(),
		newQuotaListCommand(),
	)

	return cmd
}

func newQuotaGetCommand() *cobra.Command {
	var quotaType string
	var id uint32

	cmd := &cobra.Command{
		Use:   "get [path]",
		Short: "Get quota information",
		Long:  `Get quota information for a user, group, or project.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manager := xfs.NewQuotaManager()

			qType, err := parseQuotaType(quotaType)
			if err != nil {
				return err
			}

			quota, err := manager.GetQuota(qType, id, path)
			if err != nil {
				return fmt.Errorf("failed to get quota: %w", err)
			}

			printQuotaInfo(quota)
			return nil
		},
	}

	cmd.Flags().StringVarP(&quotaType, "type", "t", "user", "quota type (user, group, project)")
	cmd.Flags().Uint32VarP(&id, "id", "i", 0, "user/group/project ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newQuotaSetCommand() *cobra.Command {
	var quotaType string
	var id uint32
	var blockSoft, blockHard string
	var inodeSoft, inodeHard uint64

	cmd := &cobra.Command{
		Use:   "set [path]",
		Short: "Set quota limits",
		Long:  `Set quota limits for a user, group, or project.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manager := xfs.NewQuotaManager()

			qType, err := parseQuotaType(quotaType)
			if err != nil {
				return err
			}

			limits := xfs.QuotaLimits{
				InodeSoft: inodeSoft,
				InodeHard: inodeHard,
			}

			// 解析块大小
			if blockSoft != "" {
				limits.BlockSoft, err = parseSize(blockSoft)
				if err != nil {
					return fmt.Errorf("invalid block soft limit: %w", err)
				}
			}

			if blockHard != "" {
				limits.BlockHard, err = parseSize(blockHard)
				if err != nil {
					return fmt.Errorf("invalid block hard limit: %w", err)
				}
			}

			err = manager.SetQuota(qType, id, path, limits)
			if err != nil {
				return fmt.Errorf("failed to set quota: %w", err)
			}

			fmt.Printf("Quota set successfully for %s ID %d\n", qType, id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&quotaType, "type", "t", "user", "quota type (user, group, project)")
	cmd.Flags().Uint32VarP(&id, "id", "i", 0, "user/group/project ID")
	cmd.Flags().StringVar(&blockSoft, "block-soft", "", "block soft limit (e.g., 1GB, 500MB)")
	cmd.Flags().StringVar(&blockHard, "block-hard", "", "block hard limit (e.g., 2GB, 1000MB)")
	cmd.Flags().Uint64Var(&inodeSoft, "inode-soft", 0, "inode soft limit")
	cmd.Flags().Uint64Var(&inodeHard, "inode-hard", 0, "inode hard limit")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newQuotaRemoveCommand() *cobra.Command {
	var quotaType string
	var id uint32

	cmd := &cobra.Command{
		Use:   "remove [path]",
		Short: "Remove quota limits",
		Long:  `Remove quota limits for a user, group, or project.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manager := xfs.NewQuotaManager()

			qType, err := parseQuotaType(quotaType)
			if err != nil {
				return err
			}

			err = manager.RemoveQuota(qType, id, path)
			if err != nil {
				return fmt.Errorf("failed to remove quota: %w", err)
			}

			fmt.Printf("Quota removed successfully for %s ID %d\n", qType, id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&quotaType, "type", "t", "user", "quota type (user, group, project)")
	cmd.Flags().Uint32VarP(&id, "id", "i", 0, "user/group/project ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newQuotaListCommand() *cobra.Command {
	var quotaType string
	var format string

	cmd := &cobra.Command{
		Use:   "list [path]",
		Short: "List all quotas",
		Long:  `List all quotas for a specific type on a filesystem.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manager := xfs.NewQuotaManager()

			qType, err := parseQuotaType(quotaType)
			if err != nil {
				return err
			}

			quotas, err := manager.GetAllQuotas(qType, path)
			if err != nil {
				return fmt.Errorf("failed to list quotas: %w", err)
			}

			switch format {
			case "table":
				printQuotasTable(quotas)
			case "json":
				printQuotasJSON(quotas)
			default:
				printQuotasTable(quotas)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&quotaType, "type", "t", "user", "quota type (user, group, project)")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json)")

	return cmd
}

// 辅助函数
func parseQuotaType(typeStr string) (xfs.QuotaType, error) {
	switch strings.ToLower(typeStr) {
	case "user", "u":
		return xfs.UserQuota, nil
	case "group", "g":
		return xfs.GroupQuota, nil
	case "project", "p":
		return xfs.ProjectQuota, nil
	default:
		return 0, fmt.Errorf("invalid quota type: %s", typeStr)
	}
}

func parseSize(sizeStr string) (uint64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	var multiplier uint64 = 1
	var numStr string

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1
		numStr = sizeStr[:len(sizeStr)-2]
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024
		numStr = sizeStr[:len(sizeStr)-2]
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024
		numStr = sizeStr[:len(sizeStr)-2]
	} else if strings.HasSuffix(sizeStr, "TB") {
		multiplier = 1024 * 1024 * 1024
		numStr = sizeStr[:len(sizeStr)-2]
	} else {
		numStr = sizeStr
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	return uint64(num * float64(multiplier)), nil
}

func printQuotaInfo(quota *xfs.QuotaInfo) {
	fmt.Printf("Quota Information:\n")
	fmt.Printf("  ID: %d\n", quota.ID)
	fmt.Printf("  Type: %s\n", quota.Type)
	fmt.Printf("  Path: %s\n", quota.Path)
	fmt.Printf("  Device: %s\n", quota.Device)
	fmt.Printf("\nBlock Usage:\n")
	fmt.Printf("  Used: %s\n", xfs.FormatSize(quota.BlockUsed*1024))
	fmt.Printf("  Soft Limit: %s\n", xfs.FormatSize(quota.BlockSoft*1024))
	fmt.Printf("  Hard Limit: %s\n", xfs.FormatSize(quota.BlockHard*1024))
	if quota.BlockHard > 0 {
		fmt.Printf("  Usage: %.1f%%\n", quota.BlockUsagePercent())
	}
	fmt.Printf("\nInode Usage:\n")
	fmt.Printf("  Used: %d\n", quota.InodeUsed)
	fmt.Printf("  Soft Limit: %d\n", quota.InodeSoft)
	fmt.Printf("  Hard Limit: %d\n", quota.InodeHard)
	if quota.InodeHard > 0 {
		fmt.Printf("  Usage: %.1f%%\n", quota.InodeUsagePercent())
	}
	fmt.Printf("\nLast Updated: %s\n", quota.LastUpdated.Format("2006-01-02 15:04:05"))
}

func printQuotasTable(quotas []xfs.QuotaInfo) {
	if len(quotas) == 0 {
		fmt.Println("No quotas found.")
		return
	}

	fmt.Printf("%-8s %-12s %-12s %-12s %-10s %-10s %-10s %-8s\n",
		"ID", "Block Used", "Block Soft", "Block Hard", "Inode Used", "Inode Soft", "Inode Hard", "Status")
	fmt.Println(strings.Repeat("-", 90))

	for _, quota := range quotas {
		status := "OK"
		if quota.IsBlockExceeded() || quota.IsInodeExceeded() {
			status = "OVER"
		} else if quota.BlockUsagePercent() > 80 || quota.InodeUsagePercent() > 80 {
			status = "WARNING"
		}

		fmt.Printf("%-8d %-12s %-12s %-12s %-10d %-10d %-10d %-8s\n",
			quota.ID,
			xfs.FormatSize(quota.BlockUsed*1024),
			xfs.FormatSize(quota.BlockSoft*1024),
			xfs.FormatSize(quota.BlockHard*1024),
			quota.InodeUsed,
			quota.InodeSoft,
			quota.InodeHard,
			status)
	}
}

func printQuotasJSON(quotas []xfs.QuotaInfo) {
	// 简化的JSON输出
	fmt.Println("[")
	for i, quota := range quotas {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": %d,\n", quota.ID)
		fmt.Printf("    \"type\": \"%s\",\n", quota.Type)
		fmt.Printf("    \"block_used\": %d,\n", quota.BlockUsed)
		fmt.Printf("    \"block_soft\": %d,\n", quota.BlockSoft)
		fmt.Printf("    \"block_hard\": %d,\n", quota.BlockHard)
		fmt.Printf("    \"inode_used\": %d,\n", quota.InodeUsed)
		fmt.Printf("    \"inode_soft\": %d,\n", quota.InodeSoft)
		fmt.Printf("    \"inode_hard\": %d\n", quota.InodeHard)
		if i < len(quotas)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")
}
