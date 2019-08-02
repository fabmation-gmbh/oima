/*
Copyright © 2019 Emanuel Bennici <benniciemanuel78@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Hidden: true,
	Deprecated: "DO NOT USE THIS COMMAND!!",
	Use:   "configure",
	Short: "Create and Setup Configuration File",
	Long: `configure creates a Config File under
 HOME/.oima.yaml and configures it with your Input.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		err := configure()

		if err != nil {
			_ = fmt.Errorf("error while creating Configuration: %s", err.Error())
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// TODO: Add Flags to make non-interactive configuration Possible
}

func configure() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(">>>>> Docker Registry Configuration <<<<<")
	fmt.Print("Enter Docker Registry URI: ")
	registryURI, _ := reader.ReadString('\n')
	Config.Regitry.RegistryURI = registryURI

	fmt.Print("Is Authentication Required? [Y|n]: ")
	authReq, _ := reader.ReadString('\n')

	if len(authReq) == 0 || strings.ToLower(authReq) == "y" {
		Config.Regitry.RequireAuth = true

		fmt.Print("Enter Username: ")
		dockerUsername, _ := reader.ReadString('\n')
		dockerPasswd, _ := terminal.ReadPassword(int(syscall.Stdin))

		if len(dockerPasswd) == 0 || len(dockerUsername) == 0 {
			fmt.Println("Authentication is required. Please enter a Valid Username or Password!!")
			os.Exit(1)
		}

		Config.Regitry.Username = dockerUsername
		Config.Regitry.Password = string(dockerPasswd)
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		_ = fmt.Errorf("unable to decode into struct, %v", err.Error())
		return err
	}

	// save Configuration
	err = viper.WriteConfig()
	if err != nil {
		_ = fmt.Errorf("Error while writing Configuration File: %s", err.Error())
		return err
	} else {
		fmt.Println("[INFO] Config File successfully written down :)")
	}

	return nil
}