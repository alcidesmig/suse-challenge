package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"
)

type InstallInterface interface {
	Install(ctx context.Context, chartName, valuesPath, releaseName, namespace, version string) error
}

type InstallService struct {
	Storage    repository.ChartStorageRepository
	Kubernetes repository.KubernetesRepository
	// Implements
	InstallInterface
}

func NewInstallService(storage repository.ChartStorageRepository, kubernetes repository.KubernetesRepository) InstallInterface {
	return &InstallService{
		Storage:    storage,
		Kubernetes: kubernetes,
	}
}

func (is *InstallService) handleMultipleVersions(chartName string, chartVersionsMetadata map[string]models.ChartMetadata) *models.ChartMetadata {
	if len(chartVersionsMetadata) > 1 {
		fmt.Printf("More than one version was found for chart \"%s\".\nAvailable versions:\n", chartName)
		for key := range chartVersionsMetadata {
			fmt.Printf("- %s\n", key)
		}
		fmt.Printf("Please use --version to specify the desired version.\n")
		os.Exit(1)
	}

	for _, value := range chartVersionsMetadata {
		return &value
	}
	return nil
}

func (is *InstallService) handleSpecificVersion(chartName string, chartVersionsMetadata map[string]models.ChartMetadata, version string) *models.ChartMetadata {
	chartMetadata, exists := chartVersionsMetadata[version]
	if !exists {
		fmt.Printf("Version \"%s\" was not found for chart \"%s\".\n", version, chartName)
		os.Exit(1)
	}
	return &chartMetadata
}

func (is *InstallService) Install(ctx context.Context, chartName, valuesPath, releaseName, namespace, version string) error {
	isRunningInsideKubernetesPod := is.Kubernetes.CheckIfKubernetesPod()
	if !isRunningInsideKubernetesPod {
		fmt.Printf("To continue, please connect to a Kubernetes pod and execute the commands inside it. The installation process must be performed within the pod environment.\n")
		os.Exit(1)
	}

	chartVersionsMetadata, err := is.Storage.Get(ctx, chartName)
	is.handleErrorAndExit(err)

	if len(chartVersionsMetadata) == 0 {
		fmt.Printf("No charts were found with the name \"%s\". Please use the \"add\" command to add it.\n", chartName)
		os.Exit(1)
	}

	var chartMetadata *models.ChartMetadata
	if version == "" {
		chartMetadata = is.handleMultipleVersions(chartName, chartVersionsMetadata)
	} else {
		chartMetadata = is.handleSpecificVersion(chartName, chartVersionsMetadata, version)
	}

	err = is.Kubernetes.Install(ctx, chartMetadata.PackagedLocalPath, valuesPath, releaseName, namespace)
	is.handleErrorAndExit(err)

	return nil
}

func (is *InstallService) handleErrorAndExit(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, repository.ErrReadingCharts) {
		fmt.Printf("Error: Failed to read charts.")
	} else if errors.Is(err, repository.ErrConfiguringCLI) {
		fmt.Printf("Error: Failed to configure CLI.")
	} else if errors.Is(err, repository.ErrReadingValues) {
		fmt.Printf("Error: Failed reading values.yaml file.")
	} else if errors.Is(err, repository.ErrInstallingChart) {
		fmt.Printf("Error: Failed while installing chart.")
	} else if errors.Is(err, repository.ErrLoadingChart) {
		fmt.Printf("Error: Failed while loading chart.")
	} else {
		fmt.Printf("Unknown error occurred.")
	}

	fmt.Printf("\nException: %s\n\n", err.Error())
	os.Exit(1)
}
