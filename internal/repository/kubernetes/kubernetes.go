package kubernetes

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"suse-cli-challenge/internal/repository"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

type KubernetesRepository struct {
	repository.KubernetesRepository
}

func (kr *KubernetesRepository) init() (*action.Configuration, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

func (kr *KubernetesRepository) readValues(valuesPath string) (output map[string]interface{}, err error) {
	if valuesPath == "" {
		return nil, nil
	}

	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Values file \"%s\" does not exists.", valuesPath)
	}
	content, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, output)
	if err != nil {
		return nil, err
	}
	return output, err
}

func (kr *KubernetesRepository) CheckIfKubernetesPod() bool {
	// https://kubernetes.io/docs/tutorials/services/connect-applications-service/#environment-variables
	_, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	return exists
}

func (kr *KubernetesRepository) Install(ctx context.Context, chartPath, valuesPath, releaseName, namespace string) error {
	values, err := kr.readValues(valuesPath)
	if err != nil {
		return fmt.Errorf("%w: %v", repository.ErrReadingValues, err.Error())
	}

	actionConfig, err := kr.init()
	if err != nil {
		return fmt.Errorf("%w: %v", repository.ErrInstallingChart, err.Error())
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("%w: %v", repository.ErrLoadingChart, err.Error())
	}
	install := createInstall(actionConfig, namespace, releaseName)

	_, err = install.Run(chart, values)
	if err != nil {
		return fmt.Errorf("%w: %v", repository.ErrInstallingChart, err.Error())
	}
	return nil
}

func NewKubernetesRepository() repository.KubernetesRepository {
	return &KubernetesRepository{}
}
