/*
Copyright © 2019 FABMation GmbH <eb@fabmation.de>

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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fabmation-gmbh/oima/internal"
	"github.com/fabmation-gmbh/oima/pkg/config"
)


var (
	cfgFile string		// Application Config File
	debug	bool		// Print Debug Messages
)

var Config config.Configuration
var applicationName = os.Args[0]

var rootCmd = &cobra.Command{
	Use:   "oima <command> [flags]",
	Short: "OCI/ Docker Image Signature Management Tool",
	Long: `oima Manages OCI/ Docker Image Signatures in you 'sigstore'.

Its impossible to keep track of all Signatures.

For Example, you have to remove the Signature for the
Docker Image 'docker.io/library/hello_world:vulnerable',
now you have to find out the Digest of the Image and
manually delete the Directory/ Signature.

This Tool automates this Process and helps to keep
track of all signed Images.`,
	Version: internal.GetVersion(),
	//RunE:     run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oima.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Print Debug Messages (defaults to false)")

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetConfigName(".oima")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// unmarshal Config Struct
	err := viper.Unmarshal(&Config)
	if err != nil {
		_ = fmt.Errorf("unable to decode into struct, %v", err.Error())
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("No Config File found! Maybe run '%s configure' first\n", applicationName)
		os.Exit(1)
	}
}
