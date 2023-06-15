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
		storage := local.NewLocalStorageRepository()
		containerImages := container_images.NewContainerImagesRepository()
		svc := service.NewContainerImagesService(storage, containerImages)
		svc.ListContainerImages(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
