package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"suse-cli-challenge/internal/repository"

	"gopkg.in/yaml.v2"
)

func (ls *LocalStorageRepository) getStorageFile() string {
	// Get the user's home directory
	homeDir, _ := os.UserHomeDir()

	// Create the ~/.config directory if it doesn't exist
	return filepath.Join(homeDir, ".config/suse-cli-challenge/charts.yaml")
}

func (ls *LocalStorageRepository) buildChartStorageKey(name string) string {
	return name
}

func (ls *LocalStorageRepository) readStoredCharts() (*ChartsStorage, string, error) {
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
