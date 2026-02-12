//go:build linux || darwin

package main

func init() {
	// Register camera-based commands
	// Linux: Uses V4L2
	// macOS: Uses AVFoundation
	rootCmd.AddCommand(observeCmd)
	rootCmd.AddCommand(assertCmd)
}
