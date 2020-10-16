/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/famartinrh/appctl/pkg/catalog"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
	"github.com/famartinrh/appctl/pkg/types/template"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"recipes"},
	Short:   "List available recipes and display application information ",
	// Long:    `Display application information, such as the available recipes that can be executed for the application.`,
	Args:    cobra.MaximumNArgs(1),
	PreRunE: loadAppConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("In the application " + appctl.AppConfig.Metadata.Name)
		fmt.Println()
		fmt.Println("Available recipes:")
		fmt.Println()
		var globalTemplate *template.Template = nil
		if appctl.AppConfig.Spec.Template != "" {
			temp, err := catalog.GetLocalTemplate(appctl.AppConfig.Spec.Template)
			if err != nil {
				return err
			}
			globalTemplate = temp
		}
		var overridenRecipes map[string]string = map[string]string{}

		var orderedRecipes []string = []string{}
		for customRecipe := range appctl.AppConfig.Spec.Recipes {
			orderedRecipes = append(orderedRecipes, customRecipe)
		}
		sort.Strings(orderedRecipes)

		for _, customRecipe := range orderedRecipes {
			v := appctl.AppConfig.Spec.Recipes[customRecipe]

			if v.Template == "" && globalTemplate == nil {
				fmt.Println("    * " + customRecipe + " [ERROR] Recipe does not come from any template (nor specific nor global)")
			} else if v.Template != "" && v.Recipe == "" {
				fmt.Println("    * " + customRecipe + " [ERROR] Recipe does not specify a recipe to be used from template \"" + v.Template + "\"")
			} else if v.Template != "" && v.Recipe != "" {
				customRecipeTemplate, err := catalog.GetLocalTemplate(v.Template)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Println("    * " + customRecipe + " [ERROR] template \"" + v.Template + "\" does not exist")
					} else {
						return err
					}
				} else {
					if _, ok := customRecipeTemplate.Recipes[v.Recipe]; ok {
						fmt.Println("    * " + customRecipe + " (from recipe \"" + v.Recipe + "\" in template \"" + v.Template + "\")")
					} else {
						fmt.Println("    * " + customRecipe + " [ERROR] Recipe \"" + v.Recipe + "\" does not exist in template \"" + v.Template + "\"")
					}
				}

				if globalTemplate != nil {
					if _, ok := globalTemplate.Recipes[customRecipe]; ok {
						overridenRecipes[customRecipe] = customRecipe
					}
				}

			} else if globalTemplate != nil && v.Template == "" && v.Recipe == "" {
				fmt.Println("    * " + customRecipe + " [ERROR] Recipe does not specify a recipe to be used from template \"" + globalTemplate.Template + "\"")
			} else if globalTemplate != nil && v.Template == "" && v.Recipe != "" {
				if _, ok := globalTemplate.Recipes[v.Recipe]; ok {
					fmt.Println("    * " + customRecipe + " (from recipe \"" + v.Recipe + "\" in template \"" + globalTemplate.Template + "\")")
				} else {
					fmt.Println("    * " + customRecipe + " [ERROR] Recipe \"" + v.Recipe + "\" does not exist in template \"" + globalTemplate.Template + "\"")
				}
			}

		}

		if globalTemplate != nil && len(overridenRecipes) < len(globalTemplate.Recipes) {

			var orderedGlobalRecipes []string = []string{}
			for globalRecipe := range globalTemplate.Recipes {
				orderedGlobalRecipes = append(orderedGlobalRecipes, globalRecipe)
			}
			sort.Strings(orderedGlobalRecipes)
			fmt.Println()
			fmt.Println("    Inherited recipes:")
			fmt.Println()
			for _, gr := range orderedGlobalRecipes {
				if _, ok := overridenRecipes[gr]; !ok {
					fmt.Println("    * " + gr + " (from template \"" + globalTemplate.Template + "\")")
				}
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringVarP(&appctl.AppFile, "file", "f", "", "app.yaml config file (default is ./app.yaml)")

}

//TODO reimplement status command to use this function
func availableRecipes() ([]*appctl.AvailableRecipe, error) {

	var availableRecipes []*appctl.AvailableRecipe = []*appctl.AvailableRecipe{}

	var globalTemplate *template.Template = nil
	if appctl.AppConfig.Spec.Template != "" {
		temp, err := catalog.GetLocalTemplate(appctl.AppConfig.Spec.Template)
		if err != nil {
			return nil, err
		}
		globalTemplate = temp
	}
	var overridenRecipes map[string]string = map[string]string{}

	var orderedRecipes []string = []string{}
	for customRecipe := range appctl.AppConfig.Spec.Recipes {
		orderedRecipes = append(orderedRecipes, customRecipe)
	}
	sort.Strings(orderedRecipes)

	for _, customRecipe := range orderedRecipes {
		v := appctl.AppConfig.Spec.Recipes[customRecipe]

		if v.Template != "" && v.Recipe != "" {
			customRecipeTemplate, err := catalog.GetLocalTemplate(v.Template)
			if err != nil {
				if !os.IsNotExist(err) {
					return nil, err
				}
			} else {
				if templateRecipe, ok := customRecipeTemplate.Recipes[v.Recipe]; ok {
					// fmt.Println("    * " + customRecipe + " (from recipe \"" + v.Recipe + "\" in template \"" + v.Template + "\")")
					availableRecipes = append(availableRecipes, &appctl.AvailableRecipe{
						RecipeName:     customRecipe,
						TemplateName:   v.Template,
						TemplateRecipe: v.Recipe,
						Recipe:         templateRecipe,
					})
				}
			}

			if globalTemplate != nil {
				if _, ok := globalTemplate.Recipes[customRecipe]; ok {
					overridenRecipes[customRecipe] = customRecipe
				}
			}

		} else if globalTemplate != nil && v.Template == "" && v.Recipe != "" {
			if templateRecipe, ok := globalTemplate.Recipes[v.Recipe]; ok {
				// fmt.Println("    * " + customRecipe + " (from recipe \"" + v.Recipe + "\" in template \"" + globalTemplate.Template + "\")")
				availableRecipes = append(availableRecipes, &appctl.AvailableRecipe{
					RecipeName:     customRecipe,
					TemplateName:   globalTemplate.Template,
					TemplateRecipe: v.Recipe,
					Recipe:         templateRecipe,
				})
			}
		}

	}

	if globalTemplate != nil && len(overridenRecipes) < len(globalTemplate.Recipes) {

		var orderedGlobalRecipes []string = []string{}
		for globalRecipe := range globalTemplate.Recipes {
			orderedGlobalRecipes = append(orderedGlobalRecipes, globalRecipe)
		}
		sort.Strings(orderedGlobalRecipes)

		// fmt.Println("Inherited recipes:")
		for _, gr := range orderedGlobalRecipes {
			if _, ok := overridenRecipes[gr]; !ok {
				// fmt.Println("    * " + gr + " (from template \"" + globalTemplate.Template + "\")")
				availableRecipes = append(availableRecipes, &appctl.AvailableRecipe{
					RecipeName:     gr,
					TemplateName:   globalTemplate.Template,
					TemplateRecipe: gr,
					Recipe:         globalTemplate.Recipes[gr],
				})
			}
		}
	}
	return availableRecipes, nil
}
