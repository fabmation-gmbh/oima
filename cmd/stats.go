/*
Copyright © 2019 Emanuel Bennici <eb@fabmation.de>

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
	"fmt"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/registry"

	"github.com/spf13/cobra"

	"os"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Shows Stats of the Registry",
	Long: `Shows Statistics of the Registry like
Number of Repositories, Images, Tags, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		// ATTENTION: This is only for testing/ debugging!
		var dockerRegistry registry.DockerRegistry
		err := dockerRegistry.Init()
		if err != nil {
			Log.Panicf("Error while Initialize DockerRegistry: %s", err.Error())
		}

		err = dockerRegistry.FetchAll()
		if err != nil {
			Log.Fatalf("Error while Fetching All Informations from Registry '%s': %s", dockerRegistry.URI, err.Error())
			os.Exit(1)
		}

		stats := dockerRegistry.Stats()
		fmt.Printf(
			"\n\n>>>>> Statistics for Registry %s <<<<<\n\n" +
				"Repositories: %8d\n" +
			"Images:       %8d\n" +
			"Tags:         %8d\n" +
			"S3Signatures: %8d\n",
			dockerRegistry.URI, stats.Repos, stats.Images, stats.Tags, stats.S3Signatures)
	},
}

func init() {
	registryCmd.AddCommand(statsCmd)
}
