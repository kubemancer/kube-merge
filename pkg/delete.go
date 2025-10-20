package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func RunDelete(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("context name required")
	}
	contextName := args[0]
	var targetConfig *clientcmdapi.Config
	targetConfig, err := clientcmd.LoadFromFile(TargetFile)
	if err != nil {
		return fmt.Errorf("loading target kubeconfig: %w", err)
	}
	context, exists := targetConfig.Contexts[contextName]
	if !exists {
		return fmt.Errorf("context '%s' not found", contextName)
	}

	fmt.Printf("Into: %s\n\n", TargetFile)

	if DryRun {
		fmt.Printf("Would delete context: %s\n", contextName)
		fmt.Println("Dry run - no changes made")
		return nil
	}
	// Clear current-context if it's the one being deleted
	if targetConfig.CurrentContext == contextName {
		var firstContext string
		for ctx := range targetConfig.Contexts {
			if ctx != contextName {
				firstContext = ctx
				break
			}
		}
		if firstContext != "" {
			targetConfig.CurrentContext = firstContext
			fmt.Printf("⚠ Switched current-context to: %s\n", firstContext)
		} else {
			targetConfig.CurrentContext = ""
			fmt.Printf("⚠ Cleared current-context (no contexts remaining)\n")
		}
	}
	clusterName := context.Cluster
	authInfo := context.AuthInfo
	delete(targetConfig.Contexts, contextName)
	delete(targetConfig.Clusters, clusterName)
	delete(targetConfig.AuthInfos, authInfo)

	err = clientcmd.WriteToFile(*targetConfig, TargetFile)
	if err != nil {
		return fmt.Errorf("deleting from kubeconfig: %w", err)
	}

	fmt.Println("\n✓ Successfully deleted context, cluster and user")
	return nil
}
