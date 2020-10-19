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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/famartinrh/appctl/pkg/catalog"
	"github.com/famartinrh/appctl/pkg/types/app"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a directory as an appctl application",
	Long: `This will create a app.yaml file in your current directory, 
this file allow you to configure how appctl will build your app for you`,
	RunE: runInitCmd,
}

var appName string
var appTemplate string

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&appName, "name", "n", "", "Application name, the shorter the better")
	initCmd.MarkFlagRequired("name")

	initCmd.Flags().StringVarP(&appTemplate, "template", "t", "", "Appctl template to use, you can find the list of available templates at https://github.com/famartinrh/appctl/tree/master/docs/catalog")
	// initCmd.MarkFlagRequired("template")

}

func runInitCmd(cmd *cobra.Command, args []string) error {

	//TODO use template descriptor to print input values for template
	if appTemplate != "" {
		_, err := catalog.GetLocalTemplate(appTemplate)
		if err != nil {
			//TODO improve error handling and error messages
			return err
		}
	}

	appConfig := app.AppConfig{
		APIVersion: "appctl.io/v1",
		Kind:       "App",
		Metadata: app.AppMetadata{
			Name: appName,
			Annotations: map[string]string{
				// description: Simple app using Quarkus Java framework
				"description":  appName + " description",
				"author":       appName + " author",
				"organization": appName + "_org",
			},
		},
		Spec: app.AppConfigSpec{
			Templates: []string{appTemplate},
		},
	}
	bytes, err := yaml.Marshal(appConfig)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(currentDir, "app.yaml"), bytes, 0664)
	if err != nil {
		return err
	}
	fmt.Println("app.yaml successfully created")
	return nil
}
