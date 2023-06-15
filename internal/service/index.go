package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"
)

type IndexInterface interface {
	Print(ctx context.Context)
	Write(ctx context.Context, filename string)
}

type IndexService struct {
	Storage    repository.ChartStorageRepository
	FileWriter repository.FileWriterRepository
	// Implements
	IndexInterface
}

func NewIndexService(storage repository.ChartStorageRepository, fileWriter repository.FileWriterRepository) IndexInterface {
	return &IndexService{
		Storage:    storage,
		FileWriter: fileWriter,
	}
}

func (is *IndexService) generate(ctx context.Context) []models.ChartVersions {
	charts, err := is.Storage.List(ctx)
	is.handleErrorAndExit(err)

	return charts
}

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
