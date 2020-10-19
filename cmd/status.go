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

	"github.com/famartinrh/appctl/pkg/core"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"recipes"},
	Short:   "List available recipes and display application information ",
	// Long:    `Display application information, such as the available recipes that can be executed for the application.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appConfig, err := core.LoadAppConfig(args, "", appctl.AppFile)
		if err != nil {
			return err
		}
		fmt.Println("In the application " + appConfig.Metadata.Name)
		fmt.Println()
		fmt.Println("Available recipes:")
		fmt.Println("  (use \"appctl <recipe>\" to execute a recipe and perform it's actions in your app")
		fmt.Println()

		parsedRecipes, err := core.AvailableRecipes(appConfig)
		if err != nil {
			return err
		}
		for _, parsedRecipe := range parsedRecipes {
			if parsedRecipe.Err != nil {
				fmt.Println("    * " + parsedRecipe.RecipeName + " [ERROR] " + parsedRecipe.Err.Error())
			} else if parsedRecipe.Multistep {
				fmt.Println("    * " + parsedRecipe.RecipeName + " (from multistep recipe)")
			} else {
				fmt.Println("    * " + parsedRecipe.RecipeName + " (from recipe \"" + parsedRecipe.TemplateRecipe + "\" in template \"" + parsedRecipe.TemplateName + "\")")
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
