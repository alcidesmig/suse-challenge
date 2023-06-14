package models

type ChartMetadata struct {
	APIVersion   string            `json:"apiVersion"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	KubeVersion  string            `json:"kubeVersion,omitempty"`
	Description  string            `json:"description,omitempty"`
	Type         string            `json:"type,omitempty"`
	Keywords     []string          `json:"keywords,omitempty"`
	Home         string            `json:"home,omitempty"`
	Sources      []string          `json:"sources,omitempty"`
	Dependencies []ChartDependency `json:"dependencies,omitempty"`
	Maintainers  []ChartMaintainer `json:"maintainers,omitempty"`
	Icon         string            `json:"icon,omitempty"`
	AppVersion   string            `json:"appVersion,omitempty"`
	Deprecated   bool              `json:"deprecated,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
}

type ChartDependency struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Repository   string   `json:"repository,omitempty"`
	Condition    string   `json:"condition,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	ImportValues []string `json:"import-values,omitempty"`
	Alias        string   `json:"alias,omitempty"`
}

type ChartMaintainer struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}
