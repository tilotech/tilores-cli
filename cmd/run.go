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
