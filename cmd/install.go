package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install", // @TODO assumption removed the "chart"
	Short: "Installs a specific Helm chart in the current Kubernetes cluster.",
	Long: `Use this command to install a specific Helm chart in the current Kubernetes cluster.
Specify the chart name to be installed, and the CLI will initiate the installation process inside a Kubernetes pod.

This command ensures that the chart's resources are deployed to the cluster according to the chart's specifications.`,
// @TODO check long doc about k8s pod
Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
