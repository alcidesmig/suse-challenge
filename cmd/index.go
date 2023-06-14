package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generates a Helm repository index file.",
	Long: `With this command, you can generate a Helm repository index file.
The CLI will scan the internal list of Helm charts and create an index file that contains metadata about each chart.
This index file is useful for clients that want to discover and install charts from your repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("index called")
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
