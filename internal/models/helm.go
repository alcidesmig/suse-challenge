package models

// ChartMetadata represents some metadata of a chart.
type ChartMetadata struct {
	// Helm
	APIVersion  string `yaml:"apiVersion"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`

	// Internal
	URL               string `yaml:"provided_url"`
	PackagedLocalPath string `yaml:"packaged_local_path"`
}

// ChartDependency represents a dependency of a chart.
type ChartDependency struct {
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	Repository   string   `yaml:"repository,omitempty"`
	ImportValues []string `yaml:"import-values,omitempty"`
	Alias        string   `yaml:"alias,omitempty"`
}

// ChartMaintainer represents the maintainer of a chart.
type ChartMaintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email,omitempty"`
	URL   string `yaml:"url,omitempty"`
}

// ChartVersions represents the versions of a chart.
type ChartVersions struct {
	Name         string             `yaml:"name"`
	VersionInfos []ChartVersionInfo `yaml:"data"`
}

// ChartVersionInfo represents the information of a specific chart version.
type ChartVersionInfo struct {
	Description       string `yaml:"description"`
	Version           string `yaml:"version"`
	URL               string `yaml:"url"`
	PackagedLocalPath string `yaml:"packaged_local_path"`
}
