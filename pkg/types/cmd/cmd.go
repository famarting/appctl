package cmd

import (
	"github.com/famartinrh/appctl/pkg/types/app"
	"github.com/famartinrh/appctl/pkg/types/template"
)

// var CfgFile string
var AppFile string

var Verbosity int

var AppConfig *app.AppConfig
var ProjectDir string
var ForceDowload bool

type AppctlConfig struct {
	Verbosity  int    `yaml:"verbosity"`
	CatalogURL string `yaml:"catalogURL"`
	Force      bool   `yaml:"force"`
}

type AvailableRecipe struct {
	RecipeName     string
	TemplateName   string
	TemplateRecipe string
	Recipe         template.TemplateRecipe
}
