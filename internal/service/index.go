package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"
)

// IndexInterface represents an interface for managing chart indexes.
type IndexInterface interface {
	Print(ctx context.Context)
	Write(ctx context.Context, filename string)
}

// IndexService is a service that implements IndexInterface.
type IndexService struct {
	Storage    repository.ChartStorageRepository
	FileWriter repository.FileWriterRepository
	// Implements
	IndexInterface
}

// NewIndexService creates a new instance of IndexInterface implementation.
func NewIndexService(storage repository.ChartStorageRepository, fileWriter repository.FileWriterRepository) IndexInterface {
	return &IndexService{
		Storage:    storage,
		FileWriter: fileWriter,
	}
}

// generate retrieves the stored charts and generates the chart index.
func (is *IndexService) generate(ctx context.Context) []models.ChartVersions {
	charts, err := is.Storage.List(ctx)
	is.handleErrorAndExit(err)

	return charts
}

// Print prints the chart index in stdout.
func (is *IndexService) Print(ctx context.Context) {
	charts := is.generate(ctx)
	for _, chart := range charts {
		fmt.Printf("Chart Name: %s\nVersions:\n", chart.Name)
		for _, info := range chart.VersionInfos {
			fmt.Printf("- Version: %s\n", info.Version)
			fmt.Printf("  Description: %s\n", info.Description)
			fmt.Printf("  URL: %s\n", info.URL)
		}
		fmt.Println()
	}
}

// Write writes the chart index to a file.
func (is *IndexService) Write(ctx context.Context, filename string) {
	if filename == "" {
		fmt.Println("Error: Please provide one name for the index file.")
		os.Exit(1)
	}
	charts := is.generate(ctx)

	err := is.FileWriter.Write(ctx, charts, filename)
	is.handleErrorAndExit(err)

	fmt.Printf("Index successfully generated in \"%s\"\n", filename)
}

// handleErrorAndExit handles errors and exits the application if necessary.
func (is *IndexService) handleErrorAndExit(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, repository.ErrReadingCharts) {
		fmt.Printf("Error: Failed to read charts.")
	} else if errors.Is(err, repository.ErrWritingFile) {
		fmt.Printf("Error: Failed while writing charts index.")
	} else {
		fmt.Printf("Unknown error occurred.")
	}

	fmt.Printf("\nException: %s\n\n", err.Error())
	os.Exit(1)
}
