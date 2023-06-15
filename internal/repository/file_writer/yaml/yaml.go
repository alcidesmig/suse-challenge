package yaml_writer

import (
	"context"
	"fmt"
	"io/ioutil"
	"suse-cli-challenge/internal/repository"

	"gopkg.in/yaml.v2"
)

type YamlFileWriterRepository struct {
	repository.FileWriterRepository
}

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

func NewYamlFileWriterRepository() repository.FileWriterRepository {
	return &YamlFileWriterRepository{}
}
