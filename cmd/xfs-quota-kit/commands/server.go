package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewServerCommand 创建服务器命令
func NewServerCommand() *cobra.Command {
	var port int
	var host string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start REST API server",
		Long:  `Start the REST API server for remote quota management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Starting XFS Quota Kit server on %s:%d\n", host, port)
			fmt.Println("REST API server is not implemented yet")

			// 这里实现实际的服务器逻辑
			// 使用gin框架创建REST API

			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "server port")
	cmd.Flags().StringVar(&host, "host", "0.0.0.0", "server host")

	return cmd
}

// NewCompletionCommand 创建自动补全命令
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(xfs-quota-kit completion bash)

Zsh:
  $ source <(xfs-quota-kit completion zsh)

Fish:
  $ xfs-quota-kit completion fish | source

PowerShell:
  PS> xfs-quota-kit completion powershell | Out-String | Invoke-Expression
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
			}
		},
	}

	return cmd
}
