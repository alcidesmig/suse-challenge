package yaml_writer

import (
	"context"
	"fmt"
	"io/ioutil"
	"suse-cli-challenge/internal/repository"

	"gopkg.in/yaml.v2"
)

// YamlFileWriterRepository implements repository.FileWriterRepository for YAML type.
type YamlFileWriterRepository struct {
	repository.FileWriterRepository
}

// Write writes the provided data as YAML to the specified file path.
func (yfr *YamlFileWriterRepository) Write(ctx context.Context, data interface{}, path string) error {
	// Convert interface to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrWritingFile, err.Error())
	}

	// Write YAML to file
	err = ioutil.WriteFile(path, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrWritingFile, err.Error())
	}
	return nil
}

// NewYamlFileWriterRepository creates a new instance of FileWriterRepository for YamlFileWriterRepository implementation.
func NewYamlFileWriterRepository() repository.FileWriterRepository {
	return &YamlFileWriterRepository{}
}
