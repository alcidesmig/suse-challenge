package cmd

import (
	"suse-cli-challenge/internal/repository/container_images"
	"suse-cli-challenge/internal/repository/storage/local"
	"suse-cli-challenge/internal/service"

	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Provides a list of container images used in all added charts.",
	Long: `This command retrieves a list of container images used in all the Helm charts that have been added to the CLI's internal list.
By executing this command, you can obtain an overview of the container images referenced by the charts.

It can be helpful for managing and tracking the container images associated with your charts.

For running this command, it is necessary to have "helm" installed.
This command works only for charts that can be rendered with "helm template" without additional arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate the implementations
		storage := local.NewLocalStorageRepository()
		containerImages := container_images.NewContainerImagesRepository()

		// Instantiate the service using the dependencies
		svc := service.NewContainerImagesService(storage, containerImages)

		// Call action
		svc.ListContainerImages(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
}
