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

	"github.com/famartinrh/appctl/pkg/types/cmd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// var appFile string

// var Verbosity int

// var appConfig *app.AppConfig
// var projectDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "appctl",
	Short: "Unified development experience using make",
	Long: `Unified development experience using make. 
With appctl you can build/test/package your applications running the same commands, 
no matter the languages or frameworks used`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appctl/appctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().IntP("verbosity", "v", 0, "number for the log level verbosity")
	viper.BindPFlag("v", rootCmd.Flags().Lookup("verbosity"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// use default config file from .appctl folder in home directory

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		appctlConfigFilePath := filepath.Join(home, ".appctl", "appctl.yaml")

		_, err = os.Stat(appctlConfigFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				appctlCfg := &cmd.AppctlConfig{
					Verbosity:  3,
					CatalogURL: "https://famartinrh.github.io/appctl",
					Force:      false,
				}
				bytes, err := yaml.Marshal(appctlCfg)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				err = ioutil.WriteFile(appctlConfigFilePath, bytes, 0664)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		viper.SetConfigFile(appctlConfigFilePath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetInt("verbosity") > 5 {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
