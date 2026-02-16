//go:build linux && cgo

package main

func init() {
	// Register style command (tree-sitter cgo dependency, Linux-only)
	rootCmd.AddCommand(styleCmd)
}
