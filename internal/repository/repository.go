package repository

import (
	"context"
	"suse-cli-challenge/internal/models"
)

type HelmRepository interface {
	ParseChartMetadata(ctx context.Context, data []byte) (*models.ChartMetadata, error)
	Package(ctx context.Context, path string) (string, error)
}

type ChartStorageRepository interface {
	Init() error
	List(ctx context.Context) ([]models.ChartVersions, error)
	Get(ctx context.Context, chartName string) (map[string]models.ChartMetadata, error)
	Append(ctx context.Context, chart models.ChartMetadata, upsert bool) error
}

type FileReaderRepository interface {
	RetrieveFileContent(ctx context.Context, path string) ([]byte, error)
	RetrieveDirectoryContent(ctx context.Context, path string) (string, error)
}

type FileWriterRepository interface {
	Write(ctx context.Context, data interface{}, path string) error
}

type KubernetesRepository interface {
	Install(ctx context.Context, chartPath, valuesPath, releaseName, namespace string) error
	CheckIfKubernetesPod() bool
}

type ContainerImagesRepository interface {
	GetReferencedContainerImagesFromChartDir(path string) ([]string, error)
}
