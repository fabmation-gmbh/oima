/*
Copyright © 2019 Emanuel Bennici <eb@fabmation.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software\n"), to deal
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
	"fmt"
	"github.com/fabmation-gmbh/oima/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Get Configuration Variables",
	Run: func(cmd *cobra.Command, args []string) {
		var conf config.Configuration

		err := viper.Unmarshal(&conf)
		if err != nil {
			_ = fmt.Errorf("unable to decode into struct, %v", err.Error())
		}

		// print Values
		fmt.Printf("Registry URI: %s\n", conf.Regitry.RegistryURI)

		if b, _ := strconv.ParseBool(conf.Regitry.RequireAuth); b {
			fmt.Printf("Authentication is required\n")
		} else {
			fmt.Printf("Authentication is not required\n")
		}

		fmt.Printf("Registry Username: %s\n", conf.Regitry.Username)

		if len(conf.Regitry.Password) > 0 {
			fmt.Printf("Registry Password is set.\n")
		} else {
			fmt.Printf("Registry Password is not set!\n")
		}

	},
}

func init() {
	rootCmd.AddCommand(confCmd)
}
