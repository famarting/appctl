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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/famartinrh/appctl/pkg/catalog"
	"github.com/famartinrh/appctl/pkg/types/app"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/manifoldco/promptui"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a directory as an appctl application",
	Long: `This will create a app.yaml file in your current directory, 
this file allow you to configure how appctl will build your app for you`,
	RunE: runInitCmd,
}

func init() {
	rootCmd.AddCommand(initCmd)

}

func runInitCmd(cmd *cobra.Command, args []string) error {

	fmt.Println("This utility will guide you through creating an app.yaml file.")
	fmt.Println("app.yaml files allow to manage the development process of your apps using appctl.")
	fmt.Println()
	fmt.Println("Press ^C at any time to quit.")
	fmt.Println()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	defaultAppName := filepath.Base(wd)

	prompt := promptui.Prompt{
		Label: "application name (" + defaultAppName + ")",
		Validate: func(s string) error {
			if strings.Contains(s, " ") {
				return errors.New("Invalid name, please use do not use whitespaces")
			}
			return nil
		},
	}
	appName, err := prompt.Run()
	if err != nil {
		return err
	}
	if appName == "" {
		appName = defaultAppName
	}

	prompt = promptui.Prompt{
		Label: "description",
	}
	appDescription, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "organization",
	}
	organization, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "author",
	}
	author, err := prompt.Run()
	if err != nil {
		return err
	}

	templates, err := catalog.ListAvailableTemplates()
	if err != nil {
		return err
	}
	templateNames := []string{}
	for _, t := range templates {
		templateNames = append(templateNames, t.Template)
	}
	templateSelectPrompt := promptui.Select{
		Label: "application template",
		Items: templateNames,
	}
	_, appTemplate, err := templateSelectPrompt.Run()
	if err != nil {
		return err
	}

	return writeAppConfigFile(appName, appDescription, author, organization, appTemplate)
}

func writeAppConfigFile(appName string, description string, author string, organization string, appTemplate string) error {

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
				"description":  description,
				"author":       author,
				"organization": organization,
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

	appConfigFilePath := filepath.Join(currentDir, "app.yaml")

	fmt.Println("About to write to " + appConfigFilePath + ":")
	fmt.Println()
	fmt.Println(string(bytes))
	fmt.Println()

	prompt := promptui.Prompt{
		Label:     "Is this OK",
		IsConfirm: true,
	}
	_, err = prompt.Run()
	if err != nil {
		fmt.Println("Cancelling...")
		return err
	}

	err = ioutil.WriteFile(appConfigFilePath, bytes, 0664)
	if err != nil {
		return err
	}
	fmt.Println("app.yaml successfully created.")
	fmt.Println()
	fmt.Println("To check the tasks you can now execute run: `appctl status`")
	return nil
}
