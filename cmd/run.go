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
	"log"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Deploys test pod into a kubernetes environment to run remote tests",
	Long: `
Run the tests against a kubernetes service. A test pod is deployed into the target namespace which will run the test scripts given.
A detailed workflow:
  - The given pod template is applied to the target namespace
  - The pod is checked to be in the Ready state before continuing
  - The test project or scripts are copied from the local machine to the remote pod
  - The tests are ran with the given entrypoint and logs are streamed back to standard out
  - The exit status is checked, if set to silent a non-zero exit code from the entrypoint will not cause the command to fail
  - If a timeout value is given the command will pause for the specified amount of time
  - After the timeout the pod resource deleted
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Println("Additional arguments were given and will be ignored")
		}

		podTemplatePath, _ := cmd.Flags().GetString("pod-template")     //check err in the future, probably move it to separate function
		namespace, _ := cmd.Flags().GetString("namespace")              //check err in the future, probably move it to separate function
		localDirectory, _ := cmd.Flags().GetString("local-directory")   //check err in the future, probably move it to separate function
		remoteDirectory, _ := cmd.Flags().GetString("remote-directory") //check err in the future, probably move it to separate function
		entryPoint, _ := cmd.Flags().GetString("entry-point")           //check err in the future, probably move it to separate function
		timeout, _ := cmd.Flags().GetString("timeout")                  //check err in the future, probably move it to separate function
		failSilently, _ := cmd.Flags().GetBool("fail-silently")         //check err in the future, probably move it to separate function

		executeTest(&TestDriver{
			podTemplatePath: podTemplatePath,
			namespace:       namespace,
			localDirectory:  localDirectory,
			remoteDirectory: remoteDirectory,
			entryPoint:      entryPoint,
			failSilently:    failSilently,
			timeout:         timeout,
			k:               &Kubectl{},
		})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("local-directory", "d", ".", "Local directory containing test project or scripts")
	runCmd.Flags().StringP("remote-directory", "r", "/home", "Remote directory where test scripts will be copied to")
	runCmd.Flags().StringP("entry-point", "e", "", "Command to run to start test")
	runCmd.Flags().StringP("pod-template", "p", "", "Path to yaml file containing pod template")
	runCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace to deploy test pod to")
	runCmd.Flags().StringP("timeout", "t", "0s", "Time that the test pod should be left alive for after tests have completed. Examples: 1h, 60s, 30m, 1h10m10s")
	runCmd.Flags().BoolP("fail-silently", "s", false, "If set to true non-zero exit code of entrypoint will not cause the command to fail")
	runCmd.MarkFlagRequired("entry-point")
	runCmd.MarkFlagRequired("pod-template")
}

func executeTest(testDriver ITestDriver) {
	fmt.Println("executeTest called")
	testDriver.createTestPod()
	testDriver.copyArtifactsToTestPod()
	testDriver.runTest()
	testDriver.cleanUpResources()
}
