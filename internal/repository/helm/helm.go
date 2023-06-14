package helm

import (
	"suse-cli-challenge/internal/repository"

	"context"

	"suse-cli-challenge/internal/models"

	"github.com/google/go-github/v37/github"
	"gopkg.in/yaml.v2"
)

type HelmRepository struct {
	GithubClient *github.Client
	repository.HelmRepository
}

func (h *HelmRepository) ParseChartMetadata(ctx context.Context, data []byte) (*models.ChartMetadata, error) {

	chartDataParsed := models.ChartMetadata{}
	err := yaml.Unmarshal([]byte(data), &chartDataParsed)

	return &chartDataParsed, err
}

func NewHelmRepository() *HelmRepository {
	return &HelmRepository{}
}
