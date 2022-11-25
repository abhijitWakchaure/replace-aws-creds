// Package cmd ...
/*
Copyright © 2022 Abhijit Wakchaure<abhijitwakchaure.2014@gmail.com>

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
	"path/filepath"

	"github.com/abhijitWakchaure/replace-aws-creds/app"
	"github.com/abhijitWakchaure/replace-aws-creds/logger"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "replace-aws-creds",
	Short: "Update the aws creadentials file in your home dir",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var log logger.Logger
		var logLevel logger.LogLevel
		daemon, _ := cmd.Flags().GetBool("daemon")
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			logLevel = logger.LogLevelDebug
		}
		log = logger.GetLogger(logLevel)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			log.Debugf("using config file: %s", viper.ConfigFileUsed())
		}
		app := &app.App{
			Logger:  log,
			Daemon:  daemon,
			Verbose: verbose,
		}
		app.ExtractAWSCreds()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable debug logs")
	rootCmd.Flags().BoolP("daemon", "d", false, "Listen to aws creds continuously")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		home = filepath.Join(home, ".aws")
		// credentialsFile := filepath.Join(home, "credentials")
		viper.AddConfigPath(home)
		viper.SetConfigName("credentials")
		viper.SetConfigType("ini")
	}
	viper.AutomaticEnv() // read in environment variables that match
}
