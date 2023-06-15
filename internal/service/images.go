package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"suse-cli-challenge/internal/repository"
)

type ContainerImagesInterface interface {
	ListContainerImages(ctx context.Context)
}

type ContainerImagesService struct {
	Storage         repository.ChartStorageRepository
	ContainerImages repository.ContainerImagesRepository
	// Implements
	IndexInterface
}

func NewContainerImagesService(storage repository.ChartStorageRepository, containerImages repository.ContainerImagesRepository) ContainerImagesInterface {
	return &ContainerImagesService{
		Storage:         storage,
		ContainerImages: containerImages,
	}
}

func (ci *ContainerImagesService) ListContainerImages(ctx context.Context) {
	charts, err := ci.Storage.List(ctx)
	ci.handleErrorAndExit(err)

	for _, chart := range charts {
		fmt.Printf("Chart: %s\n", chart.Name)
		for _, chartVersion := range chart.VersionInfos {
			images, err := ci.ContainerImages.GetReferencedContainerImagesFromChartDir(chartVersion.PackagedLocalPath)
			fmt.Printf("- Version: %s\n", chartVersion.Version)
			fmt.Printf("  Images: \n")
			if err != nil {
				fmt.Printf("  > There was an error simulating `helm template` for analyzing the chart container images. Probably the chart cannot be rendered without additional information or it has some precondition that fails without overriding some parameter.")
			}
			for _, image := range images {
				if image != "" {
					fmt.Printf("  - %s\n", image)
				}
			}
			fmt.Println()
		}
	}
}

func (ci *ContainerImagesService) handleErrorAndExit(err error) {
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
