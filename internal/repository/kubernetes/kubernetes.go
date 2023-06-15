package kubernetes

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"suse-cli-challenge/internal/repository"

	"helm.sh/helm/v3/pkg/chart/loader"
)

// ImplHelmKubernetesRepository represents a HelmKubernetesRepository implementation for installing helm charts.
type ImplHelmKubernetesRepository struct {
	repository.HelmKubernetesRepository
}

func (kr *ImplHelmKubernetesRepository) CheckIfKubernetesPod() bool {
	// https://kubernetes.io/docs/tutorials/services/connect-applications-service/#environment-variables
	_, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	return exists
}

// Install deploys a Helm chart to the current referenced Kubernetes.
func (kr *ImplHelmKubernetesRepository) Install(ctx context.Context, chartPath, valuesPath, releaseName, namespace string) error {
	values, err := kr.readValues(valuesPath)
	if err != nil {
		return fmt.Errorf("%w: %v", repository.ErrReadingValues, err.Error())
	}
	fmt.Printf("%s\n", reflect.TypeOf(values))

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

// NewImplKubernetesRepository creates a new instance of HelmKubernetesRepository for ImplKubernetesRepository implementation.
func NewImplKubernetesRepository() repository.HelmKubernetesRepository {
	return &ImplHelmKubernetesRepository{}
}
