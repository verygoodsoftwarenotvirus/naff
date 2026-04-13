package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "naff",
		Short: "NAFF generates production-ready Go service codebases from type definitions",
	}

	rootCmd.AddCommand(
		buildGenerateCmd(),
		buildValidateCmd(),
		buildInitCmd(),
		buildVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
