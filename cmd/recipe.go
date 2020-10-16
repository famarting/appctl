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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/famartinrh/appctl/pkg/catalog"
	"github.com/famartinrh/appctl/pkg/make"
	"github.com/famartinrh/appctl/pkg/types/app"
	"github.com/famartinrh/appctl/pkg/types/cmd"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
)

func execRecipe(cmd *cobra.Command, args []string, name string) error {
	err := loadAppConfig(cmd, args)
	if err == nil {
		ars, err := availableRecipes()
		if err != nil {
			return err
		}
		for _, recipe := range ars {
			if recipe.RecipeName == name {
				err := exec(recipe)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

func exec(recipe *cmd.AvailableRecipe) error {
	if appctl.Verbosity > 5 {
		fmt.Println("Executing recipe " + recipe.RecipeName)
	}

	makefile, err := catalog.GetMakefile(recipe.TemplateName, recipe.TemplateRecipe)
	if err != nil {
		return err
	}

	//global level vars
	vars := appctl.AppConfig.Spec.Vars
	//recipe level vars
	vars = append(vars, appctl.AppConfig.Spec.Recipes[recipe.RecipeName].Vars...)
	// APP_NAME
	vars = append(vars, app.InputVar{Name: "APP_NAME", Value: appctl.AppConfig.Metadata.Name})
	// APP_annotation
	for k, v := range appctl.AppConfig.Metadata.Annotations {
		vars = append(vars, app.InputVar{Name: "APP_" + strings.ToUpper(k), Value: v})
	}

	return make.BuildProject(makefile, appctl.ProjectDir, vars)
}

func loadAppConfig(cmd *cobra.Command, args []string) error {
	appctl.Verbosity = viper.GetInt("verbosity")
	appctl.ForceDowload = viper.GetBool("force")

	if appctl.Verbosity >= 10 {
		fmt.Println("Args: " + strings.Join(args, " "))
		fmt.Println("Force download " + strconv.FormatBool(appctl.ForceDowload))
	}

	//first identify app project dir
	appctl.ProjectDir = ""
	if len(args) == 2 {
		//in specific projectDir
		if args[1] != "." {
			appctl.ProjectDir = args[1]
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
	if appctl.Verbosity >= 5 {
		fmt.Println("Using dir " + appctl.ProjectDir + " as projectDir")
	}
	if appctl.AppFile == "" {
		//then look for default yaml
		appctl.AppFile = filepath.Join(appctl.ProjectDir, "app.yaml")
	}
	appctl.AppConfig = &app.AppConfig{}
	filebytes, err := ioutil.ReadFile(appctl.AppFile)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("There is no app.yaml file in this directory, try running \"appctl init\" first :)")
		}
		return err
	}
	err = yaml.Unmarshal(filebytes, appctl.AppConfig)
	if err != nil {
		return err
	}

	return nil
}
