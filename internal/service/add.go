package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"
)

type AddInterface interface {
	Add(ctx context.Context, chartLocation string, upsert bool)
}

type AddService struct {
	Helm           repository.HelmRepository
	Storage        repository.ChartStorageRepository
	FileReader     repository.FileReaderRepository
	ContainerImage repository.ContainerImagesRepository
	// Implements
	AddInterface
}

func NewAddService(
	fileReader repository.FileReaderRepository,
	helm repository.HelmRepository,
	storage repository.ChartStorageRepository,
	cimages repository.ContainerImagesRepository,
) AddInterface {
	return &AddService{
		Helm:           helm,
		Storage:        storage,
		FileReader:     fileReader,
		ContainerImage: cimages,
	}
}

func (as *AddService) getMetadata(ctx context.Context, chartLocation string) (chartContent *models.ChartMetadata, err error) {
	chartRawContent, err := as.FileReader.RetrieveFileContent(ctx, fmt.Sprintf("%s/%s", chartLocation, "Chart.yaml"))
	if err != nil {
		return nil, err
	}
	chartContent, err = as.Helm.ParseChartMetadata(ctx, chartRawContent)
	return chartContent, err
}

func (as *AddService) Add(ctx context.Context, chartLocation string, upsert bool) {
	path, err := as.FileReader.RetrieveDirectoryContent(ctx, chartLocation)
	as.handleErrorAndExit(err)

	packagedChartLocation, err := as.Helm.Package(ctx, path)
	as.handleErrorAndExit(err)

	chartContent, err := as.getMetadata(ctx, chartLocation)
	as.handleErrorAndExit(err)

	chartContent.URL = chartLocation
	chartContent.PackagedLocalPath = packagedChartLocation

	err = as.Storage.Append(ctx, *chartContent, upsert)
	as.handleErrorAndExit(err)

	fmt.Println("Chart saved successfully.")
}

func (as *AddService) handleErrorAndExit(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, repository.ErrReadingCharts) {
		fmt.Printf("Error: Failed to read charts.")
	} else if errors.Is(err, repository.ErrVersionAlreadyExists) {
		fmt.Printf("Error: Version already exists. Please use --upsert if you want to override it.")
	} else if errors.Is(err, repository.ErrSavingChart) {
		fmt.Printf("Error: Failed to save chart.")
	} else if errors.Is(err, repository.ErrConfiguringCLI) {
		fmt.Printf("Error: Failed to configure CLI.")
	} else if errors.Is(err, repository.ErrCloningRepo) {
		fmt.Printf("Error: Failed to clone repository.")
	} else if errors.Is(err, repository.ErrInvalidRepoURL) {
		fmt.Printf("Error: Invalid repository URL.")
	} else if errors.Is(err, repository.ErrGettingHomeDir) {
		fmt.Printf("Error: Failed to get user home directory.")
	} else if errors.Is(err, repository.ErrInvalidChartDir) {
		fmt.Printf("Error: Invalid chart directory.")
	} else if errors.Is(err, repository.ErrLoadingChart) {
		fmt.Printf("Error: Failed while loading chart.")
	} else {
		fmt.Printf("Unknown error occurred.")
	}

	fmt.Printf("\nException: %s\n\n", err.Error())
	os.Exit(1)
}
