package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"

	"gopkg.in/yaml.v2"
)

// ChartsStorage struct represents the structure of the YAML storage file containing the stored charts.
// It has a Charts field which maps chart names to their versions and metadata.
type ChartsStorage struct {
	Charts map[string]map[string]models.ChartMetadata `yaml:"charts"`
}

// getStorageFile returns the path to the storage file.
func (ls *LocalChartStorageRepository) getStorageFile() string {
	// Get the user's home directory
	homeDir, _ := os.UserHomeDir()

	// Create the ~/.config directory if it doesn't exist
	return filepath.Join(homeDir, ".config/suse-cli-challenge/charts.yaml")
}

// buildChartStorageKey builds the chart storage key from the given chart name.
func (ls *LocalChartStorageRepository) buildChartStorageKey(name string) string {
	return name
}

// readStoredCharts reads the stored charts from the storage file.
func (ls *LocalChartStorageRepository) readStoredCharts() (*ChartsStorage, string, error) {
	storageFile := ls.getStorageFile()

	// Read existing YAML data
	yamlData, err := ioutil.ReadFile(storageFile)
	if err != nil {
		return nil, storageFile, fmt.Errorf("%w: %s", repository.ErrReadingCharts, err.Error())
	}

	// Unmarshal existing YAML data into ChartsStorage
	data := ChartsStorage{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, storageFile, fmt.Errorf("%w: %s", repository.ErrReadingCharts, err.Error())
	}
	return &data, storageFile, nil
}
