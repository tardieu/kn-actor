/*
Copyright IBM Corporation 2022

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
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and publish an actor project",
	Long:  `Build and publish an actor project.`,
	RunE:  build,
}

var (
	buildPushOption        bool
	buildImageOption       string
	buildKindClusterOption string
)

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&buildImageOption, "image", "", "Container image name (required)")
	buildCmd.Flags().BoolVar(&buildPushOption, "push", false, "Attempt to push image")
	buildCmd.Flags().StringVar(&buildKindClusterOption, "kind", "knative", `Kind cluster name for "kind.local" images`)

	buildCmd.MarkFlagRequired("image")
}

func build(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	dockerCmd := exec.Command("docker", "build", "-t", buildImageOption, ".")
	fmt.Println(strings.Join(dockerCmd.Args, " "))
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	err := dockerCmd.Run()
	if err != nil {
		return err
	}

	if buildPushOption {
		if strings.HasPrefix(buildImageOption, "kind.local/") {
			kindCmd := exec.Command("kind", "load", "docker-image", buildImageOption, "--name", buildKindClusterOption)
			fmt.Println(strings.Join(kindCmd.Args, " "))
			kindCmd.Stdout = os.Stdout
			kindCmd.Stderr = os.Stderr
			err = kindCmd.Run()
			if err != nil {
				return err
			}
		} else {
			pushCmd := exec.Command("docker", "push", buildImageOption)
			fmt.Println(strings.Join(pushCmd.Args, " "))
			pushCmd.Stdout = os.Stdout
			pushCmd.Stderr = os.Stderr
			err = pushCmd.Run()
			if err != nil {
				return err
			}
		}
	}

	if buildPushOption {
		fmt.Println("Actor image created and published.")
	} else {
		fmt.Println("Actor image created.")
	}
	return nil
}
