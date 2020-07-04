/*
Copyright © 2020 Brandon Caton brandon.r.caton@gmail.com

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
)

// TODO: fail if run flags were given to root command
var rootCmd = &cobra.Command{
	Use:   "knitt",
	Short: "Kubernetes Native Integration Test Tool",
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
