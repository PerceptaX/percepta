//go:build linux

package main

func init() {
	// Register style command (tree-sitter cgo dependency, Linux-only)
	rootCmd.AddCommand(styleCmd)
}
