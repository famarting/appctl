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

	"github.com/famartinrh/appctl/pkg/core"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "appctl [recipe] [flags] [PATH]",
	Short: "Unified development experience across all your projects",
	// 	Long: `Unified development experience using make.
	// With appctl you can build/test/package your applications running the same commands,
	// no matter the languages or frameworks used`,
	Args:         cobra.MaximumNArgs(2),
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initializeConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		} else {
			recipe := args[0]
			return core.ExecRecipe(args, recipe, "", appctl.AppFile)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.appctl/appctl.yaml)")

	rootCmd.PersistentFlags().IntP("verbosity", "v", 0, "number for the log level verbosity")
	viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	rootCmd.PersistentFlags().Bool("force", false, "force dowload of template files and recipes")
	viper.BindPFlag("force", rootCmd.PersistentFlags().Lookup("force"))

	rootCmd.Flags().StringVarP(&appctl.AppFile, "file", "f", "", "app.yaml config file used when executing recipes (default is ./app.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initializeConfig(cmd *cobra.Command) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// use default config file from .appctl folder in home directory

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		appctlConfigFilePath := filepath.Join(home, ".appctl", "appctl.yaml")

		_, err = os.Stat(appctlConfigFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				appctlCfg := &appctl.AppctlConfig{
					Verbosity:  3,
					CatalogURL: "https://famartinrh.github.io/appctl",
					// Force:      false,
				}
				bytes, err := yaml.Marshal(appctlCfg)
				if err != nil {
					return err
				}
				os.MkdirAll(filepath.Dir(appctlConfigFilePath), os.ModePerm)
				err = ioutil.WriteFile(appctlConfigFilePath, bytes, 0664)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		viper.SetConfigFile(appctlConfigFilePath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetInt("verbosity") > 5 {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	return nil
}
