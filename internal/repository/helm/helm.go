package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"suse-cli-challenge/internal/repository"

	"context"

	"suse-cli-challenge/internal/models"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

type HelmRepository struct {
	repository.HelmRepository
}

// use chartutil?
func (h *HelmRepository) ParseChartMetadata(ctx context.Context, data []byte) (*models.ChartMetadata, error) {

	chartDataParsed := models.ChartMetadata{}
	err := yaml.Unmarshal([]byte(data), &chartDataParsed)

	return &chartDataParsed, err
}

func (h *HelmRepository) Package(ctx context.Context, path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w: %s", repository.ErrGettingHomeDir, err.Error())
	}
	ok, err := chartutil.IsChartDir(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s", repository.ErrInvalidChartDir, err.Error())
	}
	if !ok {
		return "", fmt.Errorf("%w: the target directory is not a valid chart directory", repository.ErrInvalidChartDir)
	}
	chart, err := loader.Load(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s", repository.ErrLoadingChart, err.Error())
	}
	targetDirectory := filepath.Join(homeDir, ".config/suse-cli-challenge/charts")
	packagedChartLocation, err := chartutil.Save(chart, targetDirectory)
	if err != nil {
		return "", fmt.Errorf("%w: %s", repository.ErrSavingChart, err.Error())
	}

	return packagedChartLocation, nil
}

func NewHelmRepository() repository.HelmRepository {
	return &HelmRepository{}
}
