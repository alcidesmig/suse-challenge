package cmd

import (
	"path/filepath"
	"strings"
	"suse-cli-challenge/internal/repository"
	"suse-cli-challenge/internal/repository/container_images"
	"suse-cli-challenge/internal/repository/file_reader/filesystem"
	"suse-cli-challenge/internal/repository/file_reader/github"
	"suse-cli-challenge/internal/repository/helm"
	"suse-cli-challenge/internal/repository/storage/local"

	"suse-cli-challenge/internal/service"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a Helm chart to the CLI's internal list and chosen storage.",
	Long: `This command allows you to add a Helm chart to the CLI's internal list for further processing.
You can provide the chart location as a GitHub repository URL or a local folder path.
The CLI will retrieve the chart from the specified location and store its information for future use.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags
		upsert, _ := cmd.Flags().GetBool("upsert")
		chartLocation := args[0]

		// Instantiate the implementations

		// // Dependency injection based in file path type
		var file repository.FileReaderRepository
		if strings.HasPrefix(args[0], "https://github.com") {
			file = github.NewGithubFileReaderRepository()
		} else {
			file = filesystem.NewFilesystemFileReaderRepository()
			chartLocation, _ = filepath.Abs(args[0])
		}
		helm := helm.NewImplHelmRepository()
		storage := local.NewLocalStorageRepository()
		cimgs := container_images.NewContainerImagesRepository()

		// Instantiate the service using the dependencies
		svc := service.NewAddService(file, helm, storage, cimgs)

		// Call action
		svc.Add(cmd.Context(), chartLocation, upsert)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().Bool("upsert", false, "Update the chart local data if it already exists")
}
