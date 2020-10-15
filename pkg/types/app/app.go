package app

type AppConfig struct {
	Kind       string      `yaml:"kind,omitempty"`
	APIVersion string      `yaml:"apiVersion,omitempty"`
	Metadata   AppMetadata `yaml:"metadata,omitempty"`

	Spec AppConfigSpec `yaml:"spec,omitempty"`
}

type AppMetadata struct {
	Name        string            `yaml:"name,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type AppConfigSpec struct {
	Recipes  AppRecipes `yaml:"recipes,omitempty"`
	Template string     `yaml:"template,omitempty"`
	Vars     []InputVar `yaml:"vars,omitempty"`
}

type AppRecipes struct {
	Build string `yaml:"build,omitempty"`
}

type InputVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}
