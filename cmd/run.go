package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"

	"github.com/spf13/cobra"
)

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
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			shutdownWebserver()
			os.Exit(1)
		}()
		err := runGraphQLServer()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&port, "port", "p", 8080,
		"The port where the server runs on localhost")
}

func runGraphQLServer() error {
	steps := []step{
		generateCode,
		startWebserver,
	}

	return executeSteps(steps)
}

func startWebserver() error {
	startServerCommand := createGoCommand("run", "server.go", strconv.Itoa(port))
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

func shutdownWebserver() {
	pid := atomic.LoadUint64(&serverPID)
	if pid != 0 {
		_ = syscall.Kill(-int(pid), syscall.SIGTERM)
	}
}
