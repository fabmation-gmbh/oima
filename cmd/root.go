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
	"github.com/fabmation-gmbh/oima/pkg/s3"
	"os"

	"github.com/apsdehal/go-logger"
	"github.com/awnumar/memguard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fabmation-gmbh/oima/internal"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/config"
	"github.com/fabmation-gmbh/oima/pkg/credential"
	"github.com/fabmation-gmbh/oima/pkg/ui"
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
	Run: func(cmd *cobra.Command, args []string) { ui.StartUI()	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// initialize CredStore struct
		internal.Cred = new(credential.CredStore)

		// set Log Level
		if debug {
			Log.SetLogLevel(logger.DebugLevel)
		}

		Config = internal.GetConfig()

		// map Credentials into CredStore
		if len(Config.Registry.Password) > 0 {
			// map User Password into CredStore
			err := internal.Cred.AddCredential("password", []byte(Config.Registry.Password))
			if err != nil {
				Log.Fatal(err.Error())
			}
		}

		// map S3 Credentials into CredStore
		if Config.S3.Enabled {
			// copy S3 access Key ID into CredStore
			err := internal.Cred.AddCredential("s3_accessKeyID", []byte(Config.S3.AccessKeyID))
			if err != nil { Log.Fatal(err.Error()); memguard.SafeExit(1) }

			// copy S3 Secret access Key into CredStore
			err = internal.Cred.AddCredential("s3_secretAccessKeyID", []byte(Config.S3.SecretAccessKey))
			if err != nil { Log.Fatal(err.Error()); memguard.SafeExit(1) }

			// initialize S3
			S3Data := &s3.S3Minio{}
			err = S3Data.InitS3()
			if err != nil {
				Log.Fatalf("Error while initializing MinIO S3 Client: %s", err.Error())
				memguard.SafeExit(1)
			}
		} else {
			Log.Debugf("The S3 Component of this CLI was disabled by User Configuration.")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Log.Error(err.Error())
		memguard.SafeExit(1)
	}

	memguard.SafeExit(0)
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

	// unmarshal Config struct
	err := viper.Unmarshal(&Config)
	if err != nil {
		_ = fmt.Errorf("unable to decode into struct, %v", err.Error())
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("No Config File found! Maybe run '%s configure' first\n", applicationName)
		memguard.SafeExit(1)
	}
}
