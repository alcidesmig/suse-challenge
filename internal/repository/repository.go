package repository

import (
	"context"
	"suse-cli-challenge/internal/models"
)

// HelmRepository represents a repository for interacting with Helm charts.
type HelmRepository interface {
	ParseChartMetadata(ctx context.Context, data []byte) (*models.ChartMetadata, error)
	Package(ctx context.Context, path string) (string, error)
}

// ChartStorageRepository represents a repository for storing and retrieving chart metadata.
type ChartStorageRepository interface {
	Init() error
	List(ctx context.Context) ([]models.ChartVersions, error)
	Get(ctx context.Context, chartName string) (map[string]models.ChartMetadata, error)
	Append(ctx context.Context, chart models.ChartMetadata, upsert bool) error
}

// FileReaderRepository represents a repository for reading files and directories.
type FileReaderRepository interface {
	RetrieveFileContent(ctx context.Context, path string) ([]byte, error)
	RetrieveDirectoryContent(ctx context.Context, path string) (string, error)
}

// FileWriterRepository represents a repository for writing data to files.
type FileWriterRepository interface {
	Write(ctx context.Context, data interface{}, path string) error
}

// HelmKubernetesRepository represents a repository for interacting with Kubernetes using helm charts.
type HelmKubernetesRepository interface {
	Install(ctx context.Context, chartPath, valuesPath, releaseName, namespace string) error
	CheckIfKubernetesPod() bool
}

// ContainerImagesRepository represents a repository for retrieving container images from a chart directory.
type ContainerImagesRepository interface {
	GetReferencedContainerImagesFromChartDir(path string) ([]string, error)
}
