package template

type Template struct {
	Template string `json:"template,omitempty"`

	Recipes map[string]TemplateRecipe `json:"recipes,omitempty"`
}

type TemplateRecipe struct {
	Makefile  string   `json:"makefile,omitempty"`
	InputVars []string `json:"input,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}
