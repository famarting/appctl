package app

type AppConfig struct {
	Kind       string      `yaml:"kind,omitempty"`
	APIVersion string      `yaml:"apiVersion,omitempty"`
	Metadata   AppMetadata `yaml:"metadata,omitempty"`

	Spec AppConfigSpec `yaml:"spec,omitempty"`

	//ProjectDir for internal usage
	ProjectDir string `yaml:"-"`
}

type AppMetadata struct {
	Name        string            `yaml:"name,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type AppConfigSpec struct {
	Recipes map[string][]AppRecipeTask `yaml:"recipes,omitempty"`

	//only one of template or templates can be used, appctl may throw an error if both are set
	Template  string   `yaml:"template,omitempty"`
	Templates []string `yaml:"templates,omitempty"`

	Vars []InputVar `yaml:"vars,omitempty"`
}

type AppRecipeTask struct {
	Name     string `yaml:"name,omitempty"`
	Template string `yaml:"template,omitempty"`

	//only one of recipe or recipes can be used, appctl may throw an error if both are set
	Recipe  string   `yaml:"recipe,omitempty"`
	Recipes []string `yaml:"recipes,omitempty"`

	Vars []InputVar `yaml:"vars,omitempty"`
	//special field, only used when template==appctl
	Apps []string `yaml:"apps,omitempty"`
}

type InputVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}
