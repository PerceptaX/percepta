package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "percepta",
	Short: "Perception kernel for physical hardware",
	Long:  "Percepta uses computer vision to observe, validate, and compare real-world hardware behavior.",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(observeCmd)
	rootCmd.AddCommand(assertCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(styleCmd)
	rootCmd.AddCommand(knowledgeCmd)
	rootCmd.AddCommand(generateCmd)
}
