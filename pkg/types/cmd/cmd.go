package cmd

var AppFile string

var Verbosity int
var ForceDowload bool

type AppctlConfig struct {
	Verbosity  int    `yaml:"verbosity"`
	CatalogURL string `yaml:"catalogURL"`
	Force      bool   `yaml:"force"`
}

type ParsedRecipe struct {
	RecipeName string
	Multistep  bool
	//only if no multi step
	TemplateName   string
	TemplateRecipe string

	Err error

	// tasks []ParsedRecipeTask
}

// type ParsedRecipeTask struct {
// 	TaskName     string
// 	TemplateName string

// 	TemplateRecipe string
// 	Recipe         template.TemplateRecipe

// 	TemplateRecipes []string
// }
