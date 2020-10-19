package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	appctl "github.com/famartinrh/appctl/pkg/types/cmd"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/famartinrh/appctl/pkg/types/template"
)

func GetMakefile(templateName string, recipe string) (string, error) {
	template, err := GetLocalTemplate(templateName)
	if err != nil {
		return "", err
	}

	makefileFile, err := getLocalMakefile(template, recipe)
	if err != nil {
		return "", err
	}
	return makefileFile, nil
}

func GetLocalTemplate(templateName string) (*template.Template, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	templateDir := filepath.Join(home, ".appctl", "templates", templateName)
	os.MkdirAll(templateDir, os.ModePerm)

	templateDescriptorFilePath := filepath.Join(templateDir, templateName+".json")

	if appctl.ForceDowload {
		catalogURL := viper.GetString("catalogURL")
		err = downloadFile(templateDescriptorFilePath, catalogURL+"/catalog/v1/"+templateName+"/index.json")
		if err != nil {
			return nil, err
		}
	}

	filebytes, err := ioutil.ReadFile(templateDescriptorFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			//download file
			catalogURL := viper.GetString("catalogURL")
			err = downloadFile(templateDescriptorFilePath, catalogURL+"/catalog/v1/"+templateName+"/index.json")
			if err != nil {
				return nil, err
			}
			filebytes, err = ioutil.ReadFile(templateDescriptorFilePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var template *template.Template = &template.Template{}
	err = json.Unmarshal(filebytes, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func getLocalMakefile(template *template.Template, recipe string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	templateDir := filepath.Join(home, ".appctl", "templates", template.Template)
	os.MkdirAll(templateDir, os.ModePerm)

	makefileFilePath := filepath.Join(templateDir, template.Recipes[recipe].Makefile)

	if appctl.ForceDowload {
		os.Remove(makefileFilePath)
	}

	_, err = os.Stat(makefileFilePath)
	// filebytes, err := ioutil.ReadFile(filepath.Join(templateDir, template.Recipes[recipe].Makefile))
	if err != nil {
		if os.IsNotExist(err) {
			//download file
			catalogURL := viper.GetString("catalogURL")
			err = downloadFile(makefileFilePath, catalogURL+"/catalog/v1/"+template.Template+"/"+template.Recipes[recipe].Makefile)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return makefileFilePath, nil
}

func ListAvailableTemplates() ([]template.Template, error) {
	catalogURL := viper.GetString("catalogURL")

	resp, err := http.Get(catalogURL + "/catalog/v1/index.json")
	if err != nil {
		fmt.Println("Error queriying templates " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if appctl.Verbosity >= 10 {
		fmt.Println(resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list []template.Template = []template.Template{}
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func downloadFile(filepath string, url string) error {

	if appctl.Verbosity >= 10 {
		fmt.Println("Dowloading file from " + url)
		fmt.Println("Dowloading file to " + filepath)
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file " + err.Error())
		return err
	}
	defer resp.Body.Close()

	if appctl.Verbosity >= 10 {
		fmt.Println(resp.Status)
	}

	if resp.StatusCode == 200 {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
	} else if resp.StatusCode == 404 {
		return os.ErrNotExist
	} else {
		return errors.New(resp.Status)
	}

	return nil
}
