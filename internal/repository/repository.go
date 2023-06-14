package repository

import (
	"context"
	"suse-cli-challenge/internal/models"
)

type HelmRepository interface {
	ParseChartMetadata(ctx context.Context, data []byte) (*models.ChartMetadata, error)
}

type GithubRepository interface {
	RetrieveFileContent(ctx context.Context, path string) ([]byte, error)
	ParseURL(repoURL string) (*models.GithubMetadata, error)
}
