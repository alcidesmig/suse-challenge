package cmd

import (
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
		helm := helm.NewHelmRepository()
		storage := local.NewLocalStorageRepository()
		var file repository.FileReaderRepository
		if strings.HasPrefix(args[0], "https://github.com") {
			file = github.NewGithubFileReaderRepository()
		} else {
			file = filesystem.NewFilesystemFileReaderRepository()
		}
		cimgs := container_images.NewContainerImagesRepository()
		svc := service.NewAddService(file, helm, storage, cimgs)

		upsert, _ := cmd.Flags().GetBool("upsert")
		svc.Add(cmd.Context(), args[0], upsert)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().Bool("upsert", false, "Update the chart local data if it already exists")
}
