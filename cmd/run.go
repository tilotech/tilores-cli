package cmd

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"syscall"

	"github.com/spf13/cobra"
)

//nolint:unused
var serverPID uint64

var (
	port int
)

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

	runCmd.Flags().IntVarP(&port, "port", "p", 8080,
		"The port where the server runs on localhost, can also be set on environment variable TILORES_PORT")
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
	var err error
	if port != 0 {
		err = os.Setenv("TILORES_PORT", strconv.Itoa(port))
		if err != nil {
			return fmt.Errorf("faild to set environment variable TILORES_PORT with port value %s: %v", strconv.Itoa(port), err)
		}
	}
	startServerCommand := createGoCommand("run", "server.go")
	startServerCommand.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // set group process ID for the unit test to be able to kill the process
	err = startServerCommand.Start()
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
