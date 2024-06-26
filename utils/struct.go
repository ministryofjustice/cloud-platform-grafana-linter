package utils

type YAML struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Labels    struct {
			GrafanaDashboard string `yaml:"grafana_dashboard"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Data map[string]interface{} `yaml:"data"`
}

type JSON struct {
	UID string `json:"uid"`
}

type ConfigMaps struct {
	Title     string `json:"title"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}
