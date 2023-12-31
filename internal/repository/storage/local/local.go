package local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"

	"gopkg.in/yaml.v2"
)

// LocalChartStorageRepository represents a local storage implementation for ChartStorageRepository.
type LocalChartStorageRepository struct {
	repository.ChartStorageRepository
}

func (ls *LocalChartStorageRepository) Init() error {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return repository.ErrConfiguringCLI
	}

	// Create the ~/.config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".config/suse-cli-challenge")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("%w: %s", repository.ErrConfiguringCLI, err.Error())
		}
	}

	// Define the path to the YAML file
	chatsFilePath := filepath.Join(configDir, "charts.yaml")

	// Create the file if it doesn't exist
	if _, err := os.Stat(chatsFilePath); os.IsNotExist(err) {
		_, err = os.Create(chatsFilePath)
		if err != nil {
			return fmt.Errorf("%w: %s", repository.ErrConfiguringCLI, err.Error())
		}
	} else {
		return nil
	}

	data := map[string]interface{}{"charts": map[string]interface{}{}}

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrConfiguringCLI, err.Error())
	}

	err = ioutil.WriteFile(chatsFilePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrConfiguringCLI, err.Error())
	}
	return nil
}

// Append appends a chart metadata to the local storage repository.
func (ls *LocalChartStorageRepository) Append(ctx context.Context, chart models.ChartMetadata, upsert bool) error {
	data, storageFile, err := ls.readStoredCharts()
	if err != nil {
		return err
	}

	chartKey := ls.buildChartStorageKey(chart.Name)
	chartVersions, exists := data.Charts[chartKey]

	if exists {
		if _, versionExists := chartVersions[chart.Version]; versionExists && !upsert {
			return fmt.Errorf("%w: (version %s) already exists", repository.ErrVersionAlreadyExists, chart.Version)
		}

		chartVersions[chart.Version] = chart
	} else {
		data.Charts[chartKey] = map[string]models.ChartMetadata{
			chart.Version: chart,
		}
	}

	// Marshal the updated data into YAML format
	updatedData, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrSavingChart, err.Error())
	}

	// Write the updated YAML data back to the file
	err = ioutil.WriteFile(storageFile, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrSavingChart, err.Error())
	}
	return nil
}

// Get retrieves the chart versions info for the given chart name from the local storage repository.
func (ls *LocalChartStorageRepository) Get(ctx context.Context, chartName string) (map[string]models.ChartMetadata, error) {
	data, _, err := ls.readStoredCharts()
	if err != nil {
		return nil, err
	}
	chartKey := ls.buildChartStorageKey(chartName)
	chartVersions, exists := data.Charts[chartKey]

	if exists {
		return chartVersions, nil
	}
	return nil, nil
}

// List returns a list of all the stored charts and their versions from the local storage repository.
func (ls *LocalChartStorageRepository) List(ctx context.Context) ([]models.ChartVersions, error) {
	data, _, err := ls.readStoredCharts()
	if err != nil {
		return nil, err
	}

	charts := []models.ChartVersions{}
	for name, value := range data.Charts {
		versions := []models.ChartVersionInfo{}
		for _, chartMetadata := range value {
			versions = append(versions, models.ChartVersionInfo{
				Description:       chartMetadata.Description,
				Version:           chartMetadata.Version,
				URL:               chartMetadata.URL,
				PackagedLocalPath: chartMetadata.PackagedLocalPath,
			})
		}
		charts = append(charts, models.ChartVersions{Name: name, VersionInfos: versions})
	}

	return charts, nil
}

// NewLocalStorageRepository creates a new instance of ChartStorageRepository for LocalChartStorageRepository implementation.
func NewLocalStorageRepository() repository.ChartStorageRepository {
	lr := LocalChartStorageRepository{}
	lr.Init()
	return &lr
}
