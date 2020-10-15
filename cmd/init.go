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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runInitCmd,
}

var appName string
var appTemplate string

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd.Flags().StringVarP(&appName, "name", "n", "", "Application name, better if short and simple")
	initCmd.MarkFlagRequired("name")

	initCmd.Flags().StringVarP(&appTemplate, "template", "t", "", "Appctl template to use, you can find the list of available templates _____.com")
	// initCmd.MarkFlagRequired("template")

}

func runInitCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("init called")

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
			Template: appTemplate,
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
