package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xfs-quota-kit/pkg/xfs"
)

// NewReportCommand 创建报告命令
func NewReportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate quota reports",
		Long:  `Generate comprehensive quota usage reports.`,
	}

	cmd.AddCommand(
		newReportGenerateCommand(),
		newReportFilesystemCommand(),
	)

	return cmd
}

func newReportGenerateCommand() *cobra.Command {
	var format string
	var output string

	cmd := &cobra.Command{
		Use:   "generate [path]",
		Short: "Generate quota usage report",
		Long:  `Generate a comprehensive quota usage report for the specified filesystem.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manager := xfs.NewQuotaManager()

			report, err := manager.GenerateReport(path)
			if err != nil {
				return fmt.Errorf("failed to generate report: %w", err)
			}

			switch format {
			case "table":
				printReport(report)
			case "json":
				printReportJSON(report)
			default:
				printReport(report)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")

	return cmd
}

func newReportFilesystemCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filesystem [path]",
		Short: "Show filesystem information",
		Long:  `Show detailed filesystem information and quota status.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manager := xfs.NewQuotaManager()

			info, err := manager.GetFilesystemInfo(path)
			if err != nil {
				return fmt.Errorf("failed to get filesystem info: %w", err)
			}

			isXFS, _ := manager.IsXFSFilesystem(path)

			fmt.Printf("Filesystem Information:\n")
			fmt.Printf("  Path: %s\n", path)
			fmt.Printf("  Type: %s\n", info["type"])
			fmt.Printf("  XFS: %t\n", isXFS)
			fmt.Printf("  Block Size: %v\n", info["block_size"])
			fmt.Printf("  Total Size: %s\n", info["total_size"])
			fmt.Printf("  Used Size: %s\n", info["used_size"])
			fmt.Printf("  Free Size: %s\n", info["free_size"])
			fmt.Printf("  Total Inodes: %v\n", info["total_inodes"])
			fmt.Printf("  Free Inodes: %v\n", info["free_inodes"])

			return nil
		},
	}

	return cmd
}

func printReport(report *xfs.QuotaReport) {
	fmt.Printf("Quota Report for %s\n", report.Filesystem)
	fmt.Printf("Generated at: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total Quotas: %d\n", report.TotalQuotas)
	fmt.Printf("  Over Quota: %d\n", report.OverQuotas)
	fmt.Printf("  Warning: %d\n", report.WarningQuotas)

	if len(report.Quotas) > 0 {
		fmt.Printf("\nDetailed Quota Information:\n")
		printQuotasTable(report.Quotas)
	}
}

func printReportJSON(report *xfs.QuotaReport) {
	fmt.Printf("{\n")
	fmt.Printf("  \"filesystem\": \"%s\",\n", report.Filesystem)
	fmt.Printf("  \"total_quotas\": %d,\n", report.TotalQuotas)
	fmt.Printf("  \"over_quotas\": %d,\n", report.OverQuotas)
	fmt.Printf("  \"warning_quotas\": %d,\n", report.WarningQuotas)
	fmt.Printf("  \"generated_at\": \"%s\",\n", report.GeneratedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("  \"quotas\": ")
	printQuotasJSON(report.Quotas)
	fmt.Printf("}\n")
}
