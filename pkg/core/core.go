package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/famartinrh/appctl/pkg/catalog"
	"github.com/famartinrh/appctl/pkg/make"
	"github.com/famartinrh/appctl/pkg/types/app"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
	"github.com/famartinrh/appctl/pkg/types/template"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func ExecRecipe(args []string, name string, projectDir string, appFile string) error {
	appConfig, err := LoadAppConfig(args, projectDir, appFile)
	if err != nil {
		return err
	}
	parsedRecipes, err := AvailableRecipes(appConfig)
	if err != nil {
		return err
	}
	executed := false
	for _, parsedRecipe := range parsedRecipes {
		if executed && parsedRecipe.RecipeName == name {
			return errors.New("Internal error!! duplicated recipe " + parsedRecipe.RecipeName)
		}
		if parsedRecipe.RecipeName == name {
			if parsedRecipe.Err != nil {
				return err
			}
			fmt.Println()
			fmt.Println("-> Executing recipe \"" + parsedRecipe.RecipeName + "\"")
			err := ExecuteTasks(parsedRecipe, appConfig)
			executed = true
			if err != nil {
				return err
			}
		}
	}
	if !executed {
		return errors.New("Recipe \"" + name + "\" not found")
	}
	return nil
}

func LoadAppConfig(args []string, projectDir string, appFile string) (*app.AppConfig, error) {
	appctl.Verbosity = viper.GetInt("verbosity")
	appctl.ForceDowload = viper.GetBool("force")

	if appctl.Verbosity >= 10 {
		fmt.Println("Args: " + strings.Join(args, " "))
		fmt.Println("Force download " + strconv.FormatBool(appctl.ForceDowload))
	}

	//first identify app project dir
	if projectDir == "" && len(args) == 2 {
		//in specific projectDir
		if args[1] != "." {
			projectDir = args[1]
		}
	}
	//by default in current directory
	if projectDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		projectDir = currentDir
	}
	if appctl.Verbosity >= 5 {
		fmt.Println("Using dir " + projectDir + " as projectDir")
	}
	if appFile == "" {
		//then look for default yaml
		appFile = filepath.Join(projectDir, "app.yaml")
	}
	appConfig := &app.AppConfig{}
	filebytes, err := ioutil.ReadFile(appFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("There is no app.yaml file in this directory, try running \"appctl init\" first :)")
		}
		return nil, err
	}
	err = yaml.Unmarshal(filebytes, appConfig)
	if err != nil {
		return nil, errors.New("Error parsing app.yaml " + err.Error())
	}
	appConfig.ProjectDir = projectDir
	return appConfig, nil
}

func AvailableRecipes(appConfig *app.AppConfig) ([]*appctl.ParsedRecipe, error) {

	if appConfig.Spec.Template != "" && appConfig.Spec.Templates != nil {
		return nil, errors.New("Invalid app declaration, only one of \"template\" or \"templates\" allowed")
	}

	var availableRecipes []*appctl.ParsedRecipe = []*appctl.ParsedRecipe{}

	var globalTemplates []*template.Template = []*template.Template{}
	if appConfig.Spec.Template != "" {
		temp, err := catalog.GetLocalTemplate(appConfig.Spec.Template)
		if err != nil {
			return nil, err
		}
		globalTemplates = append(globalTemplates, temp)
	}
	if appConfig.Spec.Templates != nil {
		for _, t := range appConfig.Spec.Templates {
			temp, err := catalog.GetLocalTemplate(t)
			if err != nil {
				return nil, err
			}
			globalTemplates = append(globalTemplates, temp)
		}
	}
	var globalRecipes map[string]*template.TemplateRecipeRef = map[string]*template.TemplateRecipeRef{}
	for _, globalTemplate := range globalTemplates {
		for globalRecipeName, globalRecipe := range globalTemplate.Recipes {

			if existingGlobalRecipeRef, ok := globalRecipes[globalRecipeName]; ok {
				//rename both with templatename/recipename
				delete(globalRecipes, globalRecipeName)

				globalRecipes[existingGlobalRecipeRef.Template+"/"+existingGlobalRecipeRef.Name] = existingGlobalRecipeRef

				globalRecipes[globalTemplate.Template+"/"+globalRecipeName] = &template.TemplateRecipeRef{
					Name:     globalRecipeName,
					Template: globalTemplate.Template,
					Recipe:   globalRecipe,
				}

			} else {
				//add it
				globalRecipes[globalRecipeName] = &template.TemplateRecipeRef{
					Name:     globalRecipeName,
					Template: globalTemplate.Template,
					Recipe:   globalRecipe,
				}
			}

		}
	}

	var overridenRecipes map[string]string = map[string]string{}

	var orderedRecipes []string = []string{}
	for customRecipe := range appConfig.Spec.Recipes {
		orderedRecipes = append(orderedRecipes, customRecipe)
	}
	sort.Strings(orderedRecipes)

	for _, customRecipe := range orderedRecipes {
		tasks := appConfig.Spec.Recipes[customRecipe]

		if _, ok := globalRecipes[customRecipe]; ok {
			overridenRecipes[customRecipe] = customRecipe
		}

		if len(tasks) == 0 {
			availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
				RecipeName: customRecipe,
				Err:        errors.New("Recipe does not specify any task"),
			})
		} else {
			errmsgs := []string{}
			for _, task := range tasks {
				if err := validateRecipeTask(&task); err != nil {
					errmsgs = append(errmsgs, err.Error())
				}
			}
			if len(errmsgs) != 0 {
				availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
					RecipeName: customRecipe,
					Err:        errors.New(strings.Join(errmsgs, ". ")),
				})
			} else if len(tasks) == 1 && tasks[0].Recipes == nil {
				availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
					RecipeName:     customRecipe,
					TemplateName:   tasks[0].Template,
					TemplateRecipe: tasks[0].Recipe,
				})
			} else {
				availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
					RecipeName: customRecipe,
					Multistep:  true,
				})
			}
		}

	}

	if len(overridenRecipes) < len(globalRecipes) {

		var orderedGlobalRecipes []string = []string{}
		for globalRecipe := range globalRecipes {
			orderedGlobalRecipes = append(orderedGlobalRecipes, globalRecipe)
		}
		sort.Strings(orderedGlobalRecipes)

		for _, gr := range orderedGlobalRecipes {
			//only return not overriden global recipes
			if _, ok := overridenRecipes[gr]; !ok {
				availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
					RecipeName:     gr,
					TemplateName:   globalRecipes[gr].Template,
					TemplateRecipe: globalRecipes[gr].Name,
				})
			}
		}
	}
	return availableRecipes, nil
}

func validateRecipeTask(task *app.AppRecipeTask) error {
	if task.Recipe != "" && task.Recipes != nil {
		return errors.New("Invalid task declaration, only one of \"recipe\" or \"recipes\" allowed")
	}
	if task.Template != "" && (task.Recipe != "" || task.Recipes != nil) {
		if task.Template == "appctl" {
			if task.Apps == nil {
				return errors.New("Missing apps to run tasks on")
			}
			//TODO verify apps exists
			//TODO verify recipe or recipes exists in apps
		} else {
			customRecipeTemplate, err := catalog.GetLocalTemplate(task.Template)
			if err != nil {
				if os.IsNotExist(err) {
					return errors.New("Template \"" + task.Template + "\" not found")
				}
				return err
			}
			if task.Recipe != "" {
				if _, ok := customRecipeTemplate.Recipes[task.Recipe]; !ok {
					return errors.New("Recipe \"" + task.Recipe + "\" not found")
				}
			} else {
				for _, recipe := range task.Recipes {
					if _, ok := customRecipeTemplate.Recipes[recipe]; !ok {
						return errors.New("Recipe \"" + recipe + "\" not found")
					}
				}
			}
		}

	} else {
		return errors.New("Task does not specify template and recipe")
	}
	return nil
}

func ExecuteTasks(parsedRecipe *appctl.ParsedRecipe, appConfig *app.AppConfig) error {
	// fmt.Println("-------------------------------------------")
	// defer fmt.Println("-------------------------------------------")
	if tasks, ok := appConfig.Spec.Recipes[parsedRecipe.RecipeName]; ok {
		// tasks from app.yaml
		for _, task := range tasks {
			fmt.Println()
			if task.Name != "" {
				fmt.Println("  -> Executing task \"" + task.Name + "\"")
			}
			if task.Template == "appctl" {
				err := executeAppctlTask(&task)
				if err != nil {
					return err
				}
			} else {
				//makefile based
				err := executeMakefileTask(appConfig, &task)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// tasks from global template
		err := executeMakefileTask(appConfig, &app.AppRecipeTask{
			Template: parsedRecipe.TemplateName,
			Recipe:   parsedRecipe.TemplateRecipe,
			Vars:     []app.InputVar{},
		})
		if err != nil {
			return err
		}
	}
	fmt.Println()
	fmt.Println("    [SUCCESS]")
	fmt.Println()
	return nil
}

func executeAppctlTask(task *app.AppRecipeTask) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, targetapp := range task.Apps {

		var appConfigFiles []string = []string{}
		filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
			if info.Name() == "app.yaml" {
				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				var appConfig *app.AppConfig = &app.AppConfig{}
				err = yaml.Unmarshal(bytes, appConfig)
				if err != nil {
					return errors.New("Error parsing app.yaml " + err.Error())
				}
				if appConfig.Metadata.Name == targetapp {
					appConfigFiles = append(appConfigFiles, path)
				}
			}
			return nil
		})

		if len(appConfigFiles) == 0 {
			return errors.New("There is no app called " + targetapp)
		} else if len(appConfigFiles) > 1 {
			return errors.New("There are more than one app called " + targetapp)
		}

		targetAppConfigFile := appConfigFiles[0]

		var appctlRecipesToRun []string = []string{}
		if task.Recipe != "" {
			appctlRecipesToRun = append(appctlRecipesToRun, task.Recipe)
		}
		if len(task.Recipes) != 0 {
			appctlRecipesToRun = append(appctlRecipesToRun, task.Recipes...)
		}
		for _, appctlRecipe := range appctlRecipesToRun {
			if task.Name != "" {
				fmt.Println()
				fmt.Println("    -> Task \"" + task.Name + "\" , executing recipe \"" + appctlRecipe + "\"")
				fmt.Println()
			}
			err = ExecRecipe([]string{}, appctlRecipe, filepath.Dir(targetAppConfigFile), targetAppConfigFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func executeMakefileTask(appConfig *app.AppConfig, task *app.AppRecipeTask) error {

	var templateRecipesToRun []string = []string{}

	if task.Recipe != "" {
		templateRecipesToRun = append(templateRecipesToRun, task.Recipe)
	}
	if len(task.Recipes) != 0 {
		templateRecipesToRun = append(templateRecipesToRun, task.Recipes...)
	}

	for _, templateRecipe := range templateRecipesToRun {
		makefile, err := catalog.GetMakefile(task.Template, templateRecipe)
		if err != nil {
			return err
		}

		//global level vars
		vars := appConfig.Spec.Vars
		//recipe level vars
		vars = append(vars, task.Vars...)
		// APP_NAME
		vars = append(vars, app.InputVar{Name: "APP_NAME", Value: appConfig.Metadata.Name})
		// APP_annotation
		for k, v := range appConfig.Metadata.Annotations {
			vars = append(vars, app.InputVar{Name: "APP_" + strings.ToUpper(k), Value: v})
		}

		if task.Name != "" {
			fmt.Println()
			fmt.Println("    -> Task \"" + task.Name + "\" , executing recipe \"" + templateRecipe + "\"")
			fmt.Println()
		}
		err = make.BuildProject(makefile, appConfig.ProjectDir, vars)
		if err != nil {
			return errors.New("Task failed")
		}
	}

	return nil
}
