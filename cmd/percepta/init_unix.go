//go:build linux || darwin

package main

func init() {
	// Register camera-based observe command
	// Linux: Uses V4L2
	// macOS: Uses AVFoundation
	rootCmd.AddCommand(observeCmd)
}
