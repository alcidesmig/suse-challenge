package cmd

import (
	yaml_writer "suse-cli-challenge/internal/repository/file_writer/yaml"
	"suse-cli-challenge/internal/repository/storage/local"
	"suse-cli-challenge/internal/service"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generates a Helm repository index file.",
	Long: `With this command, you can generate a Helm repository index file.
The CLI will scan the internal list of Helm charts and create an index file that contains metadata about each chart.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags
		filename, _ := cmd.Flags().GetString("file")

		// Instantiate the implementations
		storage := local.NewLocalStorageRepository()
		file := yaml_writer.NewYamlFileWriterRepository()

		// Instantiate the service using the dependencies
		svc := service.NewIndexService(storage, file)

		// Call action
		svc.Write(cmd.Context(), filename)
	},
}

// printCmd represents the "index print" command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print the local helm index",
	Long:  `Print the local helm index to shell. This commands do not generate any files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate the implementations

		storage := local.NewLocalStorageRepository()
		file := yaml_writer.NewYamlFileWriterRepository()
		// Instantiate the service using the dependencies

		svc := service.NewIndexService(storage, file)

		// Call action
		svc.Print(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.AddCommand(printCmd)

	indexCmd.Flags().String("file", "charts_index.yaml", "Name for the file to write the charts index.")
}
