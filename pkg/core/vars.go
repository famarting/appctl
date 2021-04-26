package core

import (
	"strings"

	app "github.com/famartinrh/appctl/pkg/types/app/v2"
)

func loadVars(appConfig *app.AppConfig, recipe *app.AppRecipe, step *app.AppRecipeStep) []app.InputVar {
	//global level vars
	vars := appConfig.Spec.Vars
	//recipe level vars
	vars = append(vars, recipe.Vars...)
	//step level vars
	vars = append(vars, step.Vars...)
	// APP_NAME
	vars = append(vars, app.InputVar{Name: "APP_NAME", Value: appConfig.Metadata.Name})
	// APP_annotation
	for k, v := range appConfig.Metadata.Annotations {
		vars = append(vars, app.InputVar{Name: "APP_" + strings.ToUpper(k), Value: v})
	}
	return vars
}
