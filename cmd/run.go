/*
Copyright © 2021 Tilo Tech GmbH

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
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var serverPID int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the " + applicationName + "GraphQL API webserver for testing",
	Long:  "Starts the " + applicationName + "GraphQL API webserver for testing.",
	Run: func(cmd *cobra.Command, args []string) {
		err := startWebserver()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startWebserver() error {
	startServerCommand := createGoCommand("run", "server.go")
	startServerCommand.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // set group process ID for the unit test to be able to kill the process
	err := startServerCommand.Start()
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	serverPID = startServerCommand.Process.Pid
	err = startServerCommand.Wait()
	if err != nil {
		return fmt.Errorf("an error occured while waiting on server process: %v", err)
	}
	return nil
}
