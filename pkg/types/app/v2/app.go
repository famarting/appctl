package v2

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
	Recipes map[string]AppRecipe `yaml:"recipes,omitempty"`

	//only one of template or templates can be used, appctl may throw an error if both are set
	Template  string   `yaml:"template,omitempty"`
	Templates []string `yaml:"templates,omitempty"`

	Vars []InputVar `yaml:"vars,omitempty"`
}

type AppRecipe struct {
	Description string `yaml:"description,omitempty"`

	Vars []InputVar `yaml:"vars,omitempty"`

	Steps []AppRecipeStep `yaml:"steps,omitempty"`
}

type AppRecipeStep struct {
	Name string `yaml:"name,omitempty"`

	Vars []InputVar `yaml:"vars,omitempty"`

	//recipe call mode
	Template string `yaml:"template,omitempty"`
	//only one of recipe or recipes can be used, appctl may throw an error if both are set
	Recipe  string   `yaml:"recipe,omitempty"`
	Recipes []string `yaml:"recipes,omitempty"`
	//special field, only used when template==appctl
	Apps []string `yaml:"apps,omitempty"`

	//command mode
	//special field, for running custom commands, only this or a recipe call per task
	RunCmd string `yaml:"run,omitempty"`
}

type InputVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}
