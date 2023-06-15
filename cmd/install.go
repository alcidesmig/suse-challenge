package cmd

import (
	"suse-cli-challenge/internal/repository/kubernetes"
	"suse-cli-challenge/internal/repository/storage/local"
	"suse-cli-challenge/internal/service"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install", // @TODO assumption removed the "chart"
	Short: "Installs a specific resource in the current Kubernetes cluster.",
	Long:  `Use this command to install a specific resource in the current Kubernetes cluster.`,
	// @TODO check long doc about k8s pod
	Run: func(cmd *cobra.Command, args []string) {
		chartCmd.Help()
	},
}

// chartCmd represents the chart command
var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Installs a specific Helm chart in the current Kubernetes cluster.",
	Long: `Use this command to install a specific Helm chart in the current Kubernetes cluster.
Specify the chart name to be installed, and the CLI will initiate the installation process inside a Kubernetes pod.

This command ensures that the chart's resources are deployed to the cluster according to the chart's specifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		values, _ := cmd.Flags().GetString("values")
		releaseName, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")
		version, _ := cmd.Flags().GetString("version")

		if releaseName == "" {
			releaseName = args[0]
		}
		storage := local.NewLocalStorageRepository()
		kubernetes := kubernetes.NewKubernetesRepository()

		svc := service.NewInstallService(storage, kubernetes)
		svc.Install(cmd.Context(), args[0], values, releaseName, namespace, version)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(chartCmd)

	chartCmd.Flags().String("values", "", "Path for values.yaml file")
	chartCmd.Flags().String("name", "", "Release name for installation")
	chartCmd.Flags().String("namespace", "default", "Namespace for installing the chart")
	chartCmd.Flags().String("version", "", "Version of the chart for installing")
}

// Install(ctx context.Context, chartName, valuesPath, releaseName, namespace, version string) error
