package models

type ChartMetadata struct {
	APIVersion  string `yaml:"apiVersion"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`

	URL               string `yaml:"provided_url"`
	PackagedLocalPath string `yaml:"packaged_local_path"`
}

type ChartDependency struct {
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	Repository   string   `yaml:"repository,omitempty"`
	ImportValues []string `yaml:"import-values,omitempty"`
	Alias        string   `yaml:"alias,omitempty"`
}

type ChartMaintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email,omitempty"`
	URL   string `yaml:"url,omitempty"`
}

type ChartVersions struct {
	Name         string             `yaml:"name"`
	VersionInfos []ChartVersionInfo `yaml:"data"`
}

type ChartVersionInfo struct {
	Description       string `yaml:"description"`
	Version           string `yaml:"version"`
	URL               string `yaml:"url"`
	PackagedLocalPath string `yaml:"packaged_local_path"`
}
