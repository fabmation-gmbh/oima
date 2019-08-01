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
  "github.com/spf13/cobra"
  "os"

  "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"

  "github.com/fabmation-gmbh/oima/internal"
)


var cfgFile string

var rootCmd = &cobra.Command{
 Use:   "oima <command> [flags]",
 Short: "oima Manages OCI/ Docker Image Signatures in you 'sigstore'",
 Long: `oima Manages OCI/ Docker Image Signatures in you 'sigstore'.
Why? Because its impossible to keep track of all Signatures.

For Example, you have to remove the Signature for the
Docker Image 'docker.io/library/hello_world:vulnerable',
now you have to find out the Digest of the Image and
manually delete the Directory/ Signature.

This Tool automates this Process and helps to keep
track of all signed Images.`,
  Version: internal.GetVersion(),
  Run: run,
}


func run(cmd *cobra.Command, args []string) {

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


  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func initConfig() {
  if cfgFile != "" {
    viper.SetConfigFile(cfgFile)
  } else {
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    viper.AddConfigPath(home)
    viper.SetConfigName(".oima")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}