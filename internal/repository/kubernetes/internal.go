package kubernetes

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

// init initializes the Helm action configuration.
func (kr *ImplHelmKubernetesRepository) init() (*action.Configuration, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "memory", log.Printf); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

// readValues reads and parses one helm values.yaml file.
func (kr *ImplHelmKubernetesRepository) readValues(valuesPath string) (output map[string]interface{}, err error) {
	if valuesPath == "" {
		return nil, nil
	}

	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("values file \"%s\" does not exists", valuesPath)
	}
	content, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}
	output = make(map[string]interface{})
	err = yaml.Unmarshal(content, output)
	if err != nil {
		return nil, err
	}
	return output, err
}

// createInstall creates a new action.Install instance with the provided configuration, namespace, and release name.
func createInstall(config *action.Configuration, namespace, releaseName string) *action.Install {
	install := action.NewInstall(config)
	install.ReleaseName = releaseName
	install.Namespace = namespace
	install.Wait = false
	install.WaitForJobs = false
	return install
}
