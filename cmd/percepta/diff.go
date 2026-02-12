package main

import (
	"fmt"
	"os"

	"github.com/perceptumx/percepta/internal/diff"
	"github.com/perceptumx/percepta/internal/storage"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <device> --from <firmware1> --to <firmware2>",
	Short: "Compare hardware behavior across firmware versions",
	Long: `Compares observations between two firmware versions and shows exact differences.

Shows added, removed, and modified signals between firmware versions. Useful
for detecting regressions, validating features, and understanding behavior changes.

Examples:
  # Compare two firmware versions
  percepta diff my-esp32 --from v1.0 --to v1.1

  # Check for regressions from baseline
  percepta diff my-board --from baseline --to feature-branch

  # Validate behavior change
  percepta diff test-device --from before --to after

Exit codes:
  0 - No differences detected (behavior identical)
  1 - Differences detected
  2 - Error (device not found, firmware tag missing, etc.)`,
	Args: cobra.ExactArgs(1),
	RunE: runDiff,
}

var (
	diffFromFlag string
	diffToFlag   string
)

func init() {
	diffCmd.Flags().StringVar(&diffFromFlag, "from", "", "Source firmware version tag (required)")
	diffCmd.Flags().StringVar(&diffToFlag, "to", "", "Target firmware version tag (required)")
	diffCmd.MarkFlagRequired("from")
	diffCmd.MarkFlagRequired("to")
}

func runDiff(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	// Validate flags
	if diffFromFlag == "" || diffToFlag == "" {
		return fmt.Errorf("both --from and --to flags are required")
	}

	// Initialize storage
	sqliteStorage, err := storage.NewSQLiteStorage()
	if err != nil {
		return fmt.Errorf("storage init failed: %w", err)
	}
	defer sqliteStorage.Close()

	// Get latest observation for 'from' firmware
	fromObs, err := sqliteStorage.GetLatestForFirmware(deviceID, diffFromFlag)
	if err != nil {
		return fmt.Errorf("failed to get observation for firmware '%s': %w", diffFromFlag, err)
	}

	// Get latest observation for 'to' firmware
	toObs, err := sqliteStorage.GetLatestForFirmware(deviceID, diffToFlag)
	if err != nil {
		return fmt.Errorf("failed to get observation for firmware '%s': %w", diffToFlag, err)
	}

	// Compare observations
	result := diff.Compare(fromObs, toObs)

	// Print results
	printDiffResult(result)

	// Exit with appropriate code
	if result.HasChanges() {
		os.Exit(1)
	}

	return nil
}

func printDiffResult(result *diff.DiffResult) {
	// Header
	fmt.Println("Comparing firmware versions:")
	fmt.Printf("FROM: %s (%s)\n", result.FromFirmware, result.FromTimestamp)
	fmt.Printf("TO:   %s (%s)\n", result.ToFirmware, result.ToTimestamp)
	fmt.Println()
	fmt.Printf("Device: %s\n", result.DeviceID)
	fmt.Println()

	// Check if there are any changes
	if !result.HasChanges() {
		fmt.Println("No changes detected - firmware behavior is identical.")
		return
	}

	fmt.Println("Changes detected:")
	fmt.Println()

	// Print each change with appropriate indicator
	for _, change := range result.Changes {
		switch change.Type {
		case diff.ChangeAdded:
			fmt.Printf("+ %s: %s (ADDED)\n", change.Name, change.ToState)

		case diff.ChangeRemoved:
			fmt.Printf("- %s: %s (REMOVED)\n", change.Name, change.FromState)

		case diff.ChangeModified:
			if change.Details != "" {
				fmt.Printf("~ %s: %s (%s)\n", change.Name, change.Details, "MODIFIED")
			} else {
				fmt.Printf("~ %s: %s → %s (MODIFIED)\n", change.Name, change.FromState, change.ToState)
			}
		}
	}

	// Summary
	fmt.Println()
	added, removed, modified := result.CountByType()
	fmt.Printf("Summary: %d added, %d removed, %d modified\n", added, removed, modified)
}
