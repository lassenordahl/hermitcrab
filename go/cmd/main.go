package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hermitcrab",
		Short: "Hermitcrab is a CLI for managing versions",
		Long:  `A Fast and Flexible CLI for managing and maintaining your versions powered by Cobra.`,
	}

	var versionsCmd = &cobra.Command{
		Use:   "versions",
		Short: "Manage versions",
		Long:  `Manage versions in the system`,
	}

	versionsCmd.AddCommand(cmdAdd, cmdValidate)
	rootCmd.AddCommand(versionsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
