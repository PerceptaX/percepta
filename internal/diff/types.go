package diff

import "github.com/perceptumx/percepta/internal/core"

// ChangeType represents the type of change detected
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// SignalChange represents a change in a signal between two observations
type SignalChange struct {
	Type      ChangeType
	Name      string
	FromState string // Human-readable representation of old state
	ToState   string // Human-readable representation of new state
	Details   string // Additional details about the change
}

// DiffResult contains all changes detected between two observations
type DiffResult struct {
	DeviceID      string
	FromFirmware  string
	ToFirmware    string
	FromTimestamp string
	ToTimestamp   string
	Changes       []SignalChange
}

// HasChanges returns true if any changes were detected
func (d *DiffResult) HasChanges() bool {
	return len(d.Changes) > 0
}

// CountByType returns counts of added, removed, and modified signals
func (d *DiffResult) CountByType() (added, removed, modified int) {
	for _, change := range d.Changes {
		switch change.Type {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeModified:
			modified++
		}
	}
	return
}

// NormalizedSignal represents a signal with normalized values for comparison
type NormalizedSignal struct {
	Name    string
	Type    string
	Signal  core.Signal
	OnState bool
	Blink   bool
	BlinkHz float64
	Color   core.RGB
	Text    string
	History []core.DisplayTextEntry
	Changed bool
}
