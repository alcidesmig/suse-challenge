package local

import "suse-cli-challenge/internal/models"

// chart-name:
//   - v1
//   - v2
type ChartsStorage struct {
	Charts map[string]map[string]models.ChartMetadata `yaml:"charts"`
}
