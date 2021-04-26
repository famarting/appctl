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
	"github.com/famartinrh/appctl/pkg/cmd"
	"github.com/famartinrh/appctl/pkg/make"
	app "github.com/famartinrh/appctl/pkg/types/app/v2"
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
			err := ExecuteRecipeSteps(parsedRecipe, appConfig)
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
		recipeObj := appConfig.Spec.Recipes[customRecipe]

		steps := recipeObj.Steps

		if _, ok := globalRecipes[customRecipe]; ok {
			overridenRecipes[customRecipe] = customRecipe
		}

		if len(steps) == 0 {
			availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
				RecipeName: customRecipe,
				Err:        errors.New("Recipe does not specify any step"),
			})
		} else {
			errmsgs := []string{}
			for _, step := range steps {
				if err := validateRecipeStep(&step); err != nil {
					errmsgs = append(errmsgs, err.Error())
				}
			}
			if len(errmsgs) != 0 {
				availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
					RecipeName: customRecipe,
					Err:        errors.New(strings.Join(errmsgs, ". ")),
				})
			} else if len(steps) == 1 && steps[0].Recipes == nil {
				if steps[0].RunCmd != "" {
					availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
						RecipeName:        customRecipe,
						RecipeDescription: recipeObj.Description,
						CommandMode:       true,
					})
				} else {
					availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
						RecipeName:        customRecipe,
						RecipeDescription: recipeObj.Description,
						TemplateName:      steps[0].Template,
						TemplateRecipe:    steps[0].Recipe,
					})
				}
			} else {
				availableRecipes = append(availableRecipes, &appctl.ParsedRecipe{
					RecipeName:        customRecipe,
					RecipeDescription: recipeObj.Description,
					Multistep:         true,
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

func validateRecipeStep(step *app.AppRecipeStep) error {
	if step.Recipe != "" && step.Recipes != nil {
		return errors.New("Invalid step declaration, only one of \"recipe\" or \"recipes\" allowed")
	} else if step.RunCmd != "" && step.Template != "" {
		return errors.New("Invalid step declaration, only one of \"run\" or \"template\" allowed")
	}
	if step.Template != "" && (step.Recipe != "" || step.Recipes != nil) {
		if step.Template == "appctl" {
			if step.Apps == nil {
				return errors.New("Missing apps to run steps on")
			}
			//TODO verify apps exists
			//TODO verify recipe or recipes exists in apps
		} else {
			customRecipeTemplate, err := catalog.GetLocalTemplate(step.Template)
			if err != nil {
				if os.IsNotExist(err) {
					return errors.New("Template \"" + step.Template + "\" not found")
				}
				return err
			}
			if step.Recipe != "" {
				if _, ok := customRecipeTemplate.Recipes[step.Recipe]; !ok {
					return errors.New("Recipe \"" + step.Recipe + "\" not found")
				}
			} else {
				for _, recipe := range step.Recipes {
					if _, ok := customRecipeTemplate.Recipes[recipe]; !ok {
						return errors.New("Recipe \"" + recipe + "\" not found")
					}
				}
			}
		}

	} else if step.RunCmd != "" {
		//TODO do some validation of the command
	} else {
		return errors.New("Step does not specify template/recipe or cmd")
	}
	return nil
}

func ExecuteRecipeSteps(parsedRecipe *appctl.ParsedRecipe, appConfig *app.AppConfig) error {
	// fmt.Println("-------------------------------------------")
	// defer fmt.Println("-------------------------------------------")
	if recipe, ok := appConfig.Spec.Recipes[parsedRecipe.RecipeName]; ok {
		// steps from app.yaml
		for _, step := range recipe.Steps {
			fmt.Println()
			if step.Name != "" {
				fmt.Println("  -> Executing step \"" + step.Name + "\"")
			}
			if parsedRecipe.CommandMode {
				err := executeCustomCommandStep(appConfig, parsedRecipe.RecipeName, &recipe, &step)
				if err != nil {
					return err
				}
			} else if step.Template == "appctl" {
				err := executeAppctlStep(&step)
				if err != nil {
					return err
				}
			} else {
				//makefile based
				err := executeMakefileStep(appConfig, &recipe, &step)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// steps from global template
		err := executeMakefileStep(appConfig, &app.AppRecipe{Vars: []app.InputVar{}}, &app.AppRecipeStep{
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

func executeCustomCommandStep(appConfig *app.AppConfig, recipeName string, recipe *app.AppRecipe, step *app.AppRecipeStep) error {

	vars := loadVars(appConfig, recipe, step)

	if step.Name != "" {
		fmt.Println()
		fmt.Println("    -> Step \"" + step.Name + "\" , from recipe \"" + recipeName + "\"")
		fmt.Println()
	}
	err := cmd.RunCustomCommand(step.RunCmd, appConfig.ProjectDir, vars)
	if err != nil {
		return err
	}

	return nil
}

func executeAppctlStep(step *app.AppRecipeStep) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, targetapp := range step.Apps {

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
		if step.Recipe != "" {
			appctlRecipesToRun = append(appctlRecipesToRun, step.Recipe)
		}
		if len(step.Recipes) != 0 {
			appctlRecipesToRun = append(appctlRecipesToRun, step.Recipes...)
		}
		for _, appctlRecipe := range appctlRecipesToRun {
			if step.Name != "" {
				fmt.Println()
				fmt.Println("    -> Step \"" + step.Name + "\" , from recipe \"" + appctlRecipe + "\"")
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

func executeMakefileStep(appConfig *app.AppConfig, recipe *app.AppRecipe, step *app.AppRecipeStep) error {

	var templateRecipesToRun []string = []string{}

	if step.Recipe != "" {
		templateRecipesToRun = append(templateRecipesToRun, step.Recipe)
	}
	if len(step.Recipes) != 0 {
		templateRecipesToRun = append(templateRecipesToRun, step.Recipes...)
	}

	for _, templateRecipe := range templateRecipesToRun {
		makefile, err := catalog.GetMakefile(step.Template, templateRecipe)
		if err != nil {
			return err
		}

		vars := loadVars(appConfig, recipe, step)

		if step.Name != "" {
			fmt.Println()
			fmt.Println("    -> Step \"" + step.Name + "\" , from recipe \"" + templateRecipe + "\"")
			fmt.Println()
		}
		err = make.BuildProject(makefile, appConfig.ProjectDir, vars)
		if err != nil {
			return errors.New("Step failed")
		}
	}

	return nil
}
