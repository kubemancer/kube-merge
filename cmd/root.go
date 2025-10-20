package cmd

import (
	"os"
	"path/filepath"

	"kube-merge/pkg"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kube-merge",
	Short: "Merge and delete kubeconfig credentials into and from ~/.kube/config",
	Long:  `A CLI tool to integrate and delete cluster credentials from a kubeconfig file into your main ~/.kube/config`,
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge kubeconfig credentials into ~/.kube/config",
	Long:  `A command to merge cluster credentials from a kubeconfig file into your main ~/.kube/config`,
	RunE:  pkg.RunMerge,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete kubeconfig credentials into ~/.kube/config",
	Long:  `A command to delete cluster credentials from a kubeconfig file into your main ~/.kube/config`,
	RunE:  pkg.RunDelete,
}

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultTarget := filepath.Join(homeDir, ".kube", "config")

	rootCmd.Flags().StringVarP(&pkg.TargetFile, "target", "t", defaultTarget, "Target kubeconfig file")
	rootCmd.Flags().BoolVarP(&pkg.SetContext, "set-context", "c", false, "Set merged context as current")
	rootCmd.Flags().BoolVar(&pkg.DryRun, "dry-run", false, "Show what would be merged")
	rootCmd.MarkFlagRequired("source")
	rootCmd.AddCommand(mergeCmd, deleteCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
