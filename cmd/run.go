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
	"sync/atomic"
	"syscall"

	"github.com/spf13/cobra"
)

//nolint:unused
var serverPID uint64

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the " + applicationName + "GraphQL API webserver for testing",
	Long:  "Starts the " + applicationName + "GraphQL API webserver for testing.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runGraphQLServer()
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

func runGraphQLServer() error {
	err := generateGraphQLCode()
	if err != nil {
		return err
	}
	err = startWebserver()
	if err != nil {
		return err
	}
	return nil
}

func startWebserver() error {
	startServerCommand := createGoCommand("run", "server.go")
	startServerCommand.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // set group process ID for the unit test to be able to kill the process
	err := startServerCommand.Start()
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	atomic.StoreUint64(&serverPID, uint64(startServerCommand.Process.Pid))
	err = startServerCommand.Wait()
	if err != nil {
		return fmt.Errorf("an error occurred while waiting on server process: %v", err)
	}
	return nil
}

func generateGraphQLCode() error {
	err := createGoCommand("generate", "./...").Run()
	if err != nil {
		return fmt.Errorf("failed to run go generate: %v", err)
	}
	return nil
}
