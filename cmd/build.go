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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/famartinrh/appctl/pkg/catalog"
	"github.com/famartinrh/appctl/pkg/make"
	"github.com/famartinrh/appctl/pkg/types/app"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PreRunE: loadAppConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		if appctl.Verbosity > 5 {
			fmt.Println("Building app")
		}

		var makefile string = ""
		if appctl.AppConfig.Spec.Recipes.Build != "" {
			fmt.Println("Using recipe " + appctl.AppConfig.Spec.Recipes.Build)
			makefile = appctl.AppConfig.Spec.Recipes.Build
		} else if appctl.AppConfig.Spec.Template != "" {
			fmt.Println("Using template " + appctl.AppConfig.Spec.Template)
			templateMakefile, err := catalog.GetMakefile(appctl.AppConfig.Spec.Template, "build")
			if err != nil {
				return err
			}
			makefile = templateMakefile
		} else {
			return errors.New("Nothing to execute, missing build recipe or template in app.yaml")
		}

		vars := appctl.AppConfig.Spec.Vars
		// APP_NAME
		vars = append(vars, app.InputVar{Name: "APP_NAME", Value: appctl.AppConfig.Metadata.Name})
		// APP_annotation
		for k, v := range appctl.AppConfig.Metadata.Annotations {
			vars = append(vars, app.InputVar{Name: "APP_" + strings.ToUpper(k), Value: v})
		}

		return make.BuildProject(makefile, appctl.ProjectDir, vars)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	buildCmd.PersistentFlags().StringVarP(&appctl.AppFile, "file", "f", "", "app.yaml config file (default is ./app.yaml)")

	buildCmd.Flags().Bool("force", false, "Force dowload of template files and recipes")
	viper.BindPFlag("force", buildCmd.Flags().Lookup("force"))
}

func loadAppConfig(cmd *cobra.Command, args []string) error {
	appctl.Verbosity = viper.GetInt("verbosity")
	appctl.ForceDowload = viper.GetBool("force")

	if appctl.Verbosity >= 10 {
		fmt.Println("Args: " + strings.Join(args, ", "))
	}

	//first identify app project dir
	appctl.ProjectDir = ""
	if len(args) != 0 {
		//in specific projectDir
		if args[0] != "." {
			appctl.ProjectDir = args[0]
		}
	}
	//by default in current directory
	if appctl.ProjectDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		appctl.ProjectDir = currentDir
	}
	if appctl.Verbosity > 5 {
		fmt.Println("Using dir " + appctl.ProjectDir + " as projectDir")
	}
	if appctl.AppFile == "" {
		//then look for default yaml
		appctl.AppFile = filepath.Join(appctl.ProjectDir, "app.yaml")
	}
	appctl.AppConfig = &app.AppConfig{}
	filebytes, err := ioutil.ReadFile(appctl.AppFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(filebytes, appctl.AppConfig)
	if err != nil {
		return err
	}

	return nil
}
