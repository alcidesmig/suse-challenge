package kubernetes

import "helm.sh/helm/v3/pkg/action"

func createInstall(config *action.Configuration, namespace, releaseName string) *action.Install {
	install := action.NewInstall(config)
	install.ReleaseName = releaseName
	install.Namespace = namespace
	install.Wait = false
	install.WaitForJobs = false
	return install
}
