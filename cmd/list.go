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
	"github.com/spf13/cobra"
	"os"

	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/registry"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Images of the Remote Registry",

	Run: func(cmd *cobra.Command, args []string) {
		// ATTENTION: This is only for testing/ debugging!
		var dockerRegistry registry.DockerRegistry
		err := dockerRegistry.Init()
		if err != nil {
			Log.Panicf("Error while Initialize DockerRegistry: %s", err.Error())
		}

		// list Repos
		repos := dockerRegistry.ListRepositories()
		Log.Debugf("Returned %d Repositories", len(repos))

		err = dockerRegistry.FetchAll()
		if err != nil {
			Log.Fatalf("Error while Fetching All Informations from Registry '%s': %s", dockerRegistry.URI, err.Error())
			os.Exit(1)
		}

		for _, v := range dockerRegistry.ListRepositories() {
			Log.Noticef(">>>>> Repository: %s (%d Images) <<<<<", v.Name, len(v.Images))
			image, _ := v.ListImages()

			for _, img := range image {
				Log.Debugf(">> Image: %s", img.Name)
				tags, err := img.ListImageTags()

				if err != nil {
					Log.Errorf("Error while getting Image Tags of Image '%s': %s", img.Name, err.Error())
					os.Exit(1)
				}

				for _, tag := range tags { Log.Infof("  -- Tag: %s || Content Digest: %s", tag.TagName, tag.ContentDigest) }
			}
		}
	},
}

func init() {
	imageCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
