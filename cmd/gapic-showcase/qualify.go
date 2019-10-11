package main

import (
	"fmt"

	"github.com/googleapis/gapic-showcase/cmd/gapic-showcase/qualifier"

	"github.com/spf13/cobra"
)

func init() {
	qualifyCmd := &cobra.Command{
		Use:   "qualify",
		Short: "Tests a provided gapic generator against an acceptance suite",
		Run:   qualifier.RunProcess,
	}
	rootCmd.AddCommand(qualifyCmd)
}

func qualifyStub(cmd *cobra.Command, args []string) {
	fmt.Printf("in qualify!\n")
}
