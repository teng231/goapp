package main

import (
	"fmt"
	"log"
	"os"
	"teng231/goapp/cmd"

	"github.com/spf13/cobra"
)

func cliInitilized() *cobra.Command {
	startCLICMD := &cobra.Command{Use: "start", Short: "Start the application", Run: cmd.StartCLI}
	// syncDBCLICMD := &cobra.Command{Use: "migrate", Short: "Run database migrations", Run: syncDBCLI}
	appCommands := []*cobra.Command{
		startCLICMD,
		// syncDBCLICMD,
	}
	// // add more flag
	// startCLICMD.Flags().String("yamlPath", "", "set config mode is read file yaml")
	// syncDBCLICMD.Flags().String("yamlPath", "", "set config mode is read file yaml")
	rootCommand := &cobra.Command{}
	rootCommand.AddCommand(appCommands...)
	return rootCommand
}

func main() {
	log.Print("🔥🔥🔥🔥$ app cli init ")
	// Execute the root command
	if err := cliInitilized().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
