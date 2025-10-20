package pkg

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	ClusterToDelete string
	TargetFile      string
	SetContext      bool
	DryRun          bool
)

func RunMerge(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("source file path required")
	}
	sourceFile := args[0]
	if len(sourceFile) == 0 {
		return fmt.Errorf("source file path cannot be empty")
	}

	sourceConfig, err := clientcmd.LoadFromFile(sourceFile)
	if err != nil {
		return fmt.Errorf("loading source kubeconfig: %w", err)
	}

	var targetConfig *clientcmdapi.Config
	if _, err := os.Stat(TargetFile); err == nil {
		targetConfig, err = clientcmd.LoadFromFile(TargetFile)
		if err != nil {
			return fmt.Errorf("loading target kubeconfig: %w", err)
		}
	} else {
		targetConfig = clientcmdapi.NewConfig()
	}

	fmt.Printf("Merging from: %s\n", sourceFile)
	fmt.Printf("Into: %s\n\n", TargetFile)
	fmt.Printf("  Clusters: %v\n", getKeys(sourceConfig.Clusters))
	fmt.Printf("  Users: %v\n", getKeys(sourceConfig.AuthInfos))
	fmt.Printf("  Contexts: %v\n", getKeys(sourceConfig.Contexts))

	if DryRun {
		fmt.Println("\nDry run - no changes made")
		return nil
	}

	for name, cluster := range sourceConfig.Clusters {
		targetConfig.Clusters[name] = cluster
	}
	for name, authInfo := range sourceConfig.AuthInfos {
		targetConfig.AuthInfos[name] = authInfo
	}
	for name, context := range sourceConfig.Contexts {
		targetConfig.Contexts[name] = context
	}

	if SetContext && sourceConfig.CurrentContext != "" {
		targetConfig.CurrentContext = sourceConfig.CurrentContext
		fmt.Printf("\nCurrent context set to: %s\n", sourceConfig.CurrentContext)
	}

	err = clientcmd.WriteToFile(*targetConfig, TargetFile)
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Println("\n✓ Successfully merged kubeconfig")
	return nil
}

func getKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
