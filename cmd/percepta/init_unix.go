//go:build !windows

package main

func init() {
	// Register camera-based observe command (requires V4L2 on Linux/macOS)
	rootCmd.AddCommand(observeCmd)
}
