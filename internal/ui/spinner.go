package ui

import (
	"fmt"
	"os"
	"time"
)

// Spinner displays a progress indicator for long-running operations
type Spinner struct {
	message string
	done    chan bool
	stopped bool
}

// NewSpinner creates and starts a new spinner with the given message
func NewSpinner(message string) *Spinner {
	s := &Spinner{
		message: message,
		done:    make(chan bool),
		stopped: false,
	}
	go s.spin()
	return s
}

func (s *Spinner) spin() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-s.done:
			return
		default:
			fmt.Fprintf(os.Stderr, "\r%s %s", frames[i], s.message)
			i = (i + 1) % len(frames)
			time.Sleep(80 * time.Millisecond)
		}
	}
}

// Stop stops the spinner and shows a completion indicator
func (s *Spinner) Stop(success bool) {
	if s.stopped {
		return
	}
	s.stopped = true
	s.done <- true
	close(s.done)

	if success {
		fmt.Fprintf(os.Stderr, "\r✓ %s\n", s.message)
	} else {
		fmt.Fprintf(os.Stderr, "\r✗ %s\n", s.message)
	}
}

// StopWithMessage stops the spinner with a custom message
func (s *Spinner) StopWithMessage(success bool, message string) {
	if s.stopped {
		return
	}
	s.stopped = true
	s.done <- true
	close(s.done)

	if success {
		fmt.Fprintf(os.Stderr, "\r✓ %s\n", message)
	} else {
		fmt.Fprintf(os.Stderr, "\r✗ %s\n", message)
	}
}
