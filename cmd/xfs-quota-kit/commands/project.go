package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xfs-quota-kit/pkg/xfs"
)

// NewProjectCommand 创建项目管理命令
func NewProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage XFS project quotas",
		Long:  `Manage XFS project quotas and project directories.`,
	}

	cmd.AddCommand(
		newProjectCreateCommand(),
		newProjectRemoveCommand(),
		newProjectListCommand(),
	)

	return cmd
}

func newProjectCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [path]",
		Short: "Create a new project",
		Long:  `Create a new XFS project quota for the specified directory.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			path := args[1]
			manager := xfs.NewQuotaManager()

			project, err := manager.CreateProject(name, path)
			if err != nil {
				return fmt.Errorf("failed to create project: %w", err)
			}

			fmt.Printf("Project created successfully:\n")
			fmt.Printf("  Name: %s\n", project.Name)
			fmt.Printf("  ID: %d\n", project.ID)
			fmt.Printf("  Path: %s\n", project.Path)
			return nil
		},
	}

	return cmd
}

func newProjectRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a project",
		Long:  `Remove an XFS project quota.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			manager := xfs.NewQuotaManager()

			err := manager.RemoveProject(name)
			if err != nil {
				return fmt.Errorf("failed to remove project: %w", err)
			}

			fmt.Printf("Project '%s' removed successfully\n", name)
			return nil
		},
	}

	return cmd
}

func newProjectListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Long:  `List all XFS project quotas.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := xfs.NewQuotaManager()

			projects, err := manager.GetProjects()
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}

			if len(projects) == 0 {
				fmt.Println("No projects found.")
				return nil
			}

			fmt.Printf("%-8s %-20s %s\n", "ID", "Name", "Path")
			fmt.Println("----------------------------------------")
			for _, project := range projects {
				fmt.Printf("%-8d %-20s %s\n", project.ID, project.Name, project.Path)
			}

			return nil
		},
	}

	return cmd
}
