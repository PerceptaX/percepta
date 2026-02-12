package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "percepta",
	Short: "AI firmware development with hardware validation",
	Long: `Percepta uses computer vision to observe, validate, and compare real-world hardware behavior.

Vision-based hardware testing for embedded systems:
- Observe LED states, displays, and boot behavior via camera
- Assert expected behavior with assertions
- Compare firmware versions to detect regressions
- Generate BARR-C compliant code with AI validation

Quick Start:
  1. Add device:    percepta device add my-board
  2. Observe:       percepta observe my-board
  3. Assert:        percepta assert my-board "led power is ON"
  4. Generate code: percepta generate "Blink LED at 1Hz" --board esp32

Learn more: https://github.com/Perceptax/percepta`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(knowledgeCmd)
	rootCmd.AddCommand(generateCmd)
}
