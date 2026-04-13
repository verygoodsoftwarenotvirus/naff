package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func buildVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("naff %s\n", version)
		},
	}
}
