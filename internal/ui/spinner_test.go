package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner("Testing...")
	defer spinner.Stop(true)

	if spinner == nil {
		t.Fatal("Expected non-nil spinner")
	}

	if spinner.message != "Testing..." {
		t.Errorf("Expected message 'Testing...', got '%s'", spinner.message)
	}

	if spinner.stopped {
		t.Error("Expected spinner to not be stopped initially")
	}

	if spinner.done == nil {
		t.Error("Expected done channel to be initialized")
	}
}

func TestSpinner_Stop_Success(t *testing.T) {
	spinner := NewSpinner("Operation")

	// Let it spin for a moment
	time.Sleep(100 * time.Millisecond)

	output := captureStderr(func() {
		spinner.Stop(true)
	})

	// Verify success indicator is present
	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator ✓")
	}

	if !strings.Contains(output, "Operation") {
		t.Error("Expected message in output")
	}

	// Verify spinner is stopped
	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_Stop_Failure(t *testing.T) {
	spinner := NewSpinner("Operation")

	// Let it spin for a moment
	time.Sleep(100 * time.Millisecond)

	output := captureStderr(func() {
		spinner.Stop(false)
	})

	// Verify failure indicator is present
	if !strings.Contains(output, "✗") {
		t.Error("Expected failure indicator ✗")
	}

	if !strings.Contains(output, "Operation") {
		t.Error("Expected message in output")
	}

	// Verify spinner is stopped
	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_StopWithMessage_Success(t *testing.T) {
	spinner := NewSpinner("Loading")

	// Let it spin for a moment
	time.Sleep(100 * time.Millisecond)

	output := captureStderr(func() {
		spinner.StopWithMessage(true, "Completed successfully")
	})

	// Verify success indicator is present
	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator ✓")
	}

	if !strings.Contains(output, "Completed successfully") {
		t.Error("Expected custom message in output")
	}

	// Original message should not be in final output (only custom message)
	if strings.Contains(output, "Loading") && !strings.Contains(output, "Completed") {
		t.Error("Expected custom message to replace original")
	}
}

func TestSpinner_StopWithMessage_Failure(t *testing.T) {
	spinner := NewSpinner("Processing")

	// Let it spin for a moment
	time.Sleep(100 * time.Millisecond)

	output := captureStderr(func() {
		spinner.StopWithMessage(false, "Failed to process")
	})

	// Verify failure indicator is present
	if !strings.Contains(output, "✗") {
		t.Error("Expected failure indicator ✗")
	}

	if !strings.Contains(output, "Failed to process") {
		t.Error("Expected custom message in output")
	}
}

func TestSpinner_DoubleStop(t *testing.T) {
	spinner := NewSpinner("Test")

	// Let it spin for a moment
	time.Sleep(100 * time.Millisecond)

	// First stop
	output1 := captureStderr(func() {
		spinner.Stop(true)
	})

	if !strings.Contains(output1, "✓") {
		t.Error("Expected success indicator on first stop")
	}

	// Second stop should be a no-op
	output2 := captureStderr(func() {
		spinner.Stop(false)
	})

	// Second stop should produce no output
	if strings.Contains(output2, "✓") || strings.Contains(output2, "✗") {
		t.Error("Expected no output on second stop")
	}

	// Spinner should still be marked as stopped
	if !spinner.stopped {
		t.Error("Expected spinner to remain stopped")
	}
}

func TestSpinner_DoubleStopWithMessage(t *testing.T) {
	spinner := NewSpinner("Test")

	// Let it spin for a moment
	time.Sleep(100 * time.Millisecond)

	// First stop
	output1 := captureStderr(func() {
		spinner.StopWithMessage(true, "Done")
	})

	if !strings.Contains(output1, "Done") {
		t.Error("Expected custom message on first stop")
	}

	// Second stop should be a no-op
	output2 := captureStderr(func() {
		spinner.StopWithMessage(false, "Failed")
	})

	// Second stop should produce no output
	if strings.Contains(output2, "Failed") {
		t.Error("Expected no output on second stop with message")
	}
}

func TestSpinner_ImmediateStop(t *testing.T) {
	// Test stopping immediately after creation
	spinner := NewSpinner("Quick test")

	output := captureStderr(func() {
		spinner.Stop(true)
	})

	// Should still work correctly
	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator even for immediate stop")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_Lifecycle(t *testing.T) {
	// Test complete lifecycle: create -> spin -> stop
	spinner := NewSpinner("Lifecycle test")

	// Verify initial state
	if spinner.stopped {
		t.Error("Expected spinner to not be stopped initially")
	}

	// Let it spin
	time.Sleep(200 * time.Millisecond)

	// Stop successfully
	output := captureStderr(func() {
		spinner.Stop(true)
	})

	// Verify final state
	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator in output")
	}
}

func TestSpinner_MultipleSpinners(t *testing.T) {
	// Test that multiple spinners can coexist
	spinner1 := NewSpinner("Task 1")
	spinner2 := NewSpinner("Task 2")
	spinner3 := NewSpinner("Task 3")

	// Let them spin
	time.Sleep(150 * time.Millisecond)

	// Stop them in different orders
	captureStderr(func() {
		spinner2.Stop(true)
	})

	time.Sleep(50 * time.Millisecond)

	captureStderr(func() {
		spinner1.Stop(false)
	})

	time.Sleep(50 * time.Millisecond)

	captureStderr(func() {
		spinner3.Stop(true)
	})

	// Verify all are stopped
	if !spinner1.stopped || !spinner2.stopped || !spinner3.stopped {
		t.Error("Expected all spinners to be stopped")
	}
}

func TestSpinner_LongRunning(t *testing.T) {
	// Test a spinner that runs for a longer duration
	spinner := NewSpinner("Long operation")

	// Let it spin through multiple frames
	time.Sleep(500 * time.Millisecond)

	output := captureStderr(func() {
		spinner.Stop(true)
	})

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator after long run")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_EmptyMessage(t *testing.T) {
	// Test spinner with empty message
	spinner := NewSpinner("")

	time.Sleep(100 * time.Millisecond)

	output := captureStderr(func() {
		spinner.Stop(true)
	})

	// Should still show indicator
	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator even with empty message")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_LongMessage(t *testing.T) {
	// Test spinner with a very long message
	longMessage := strings.Repeat("Very long message ", 20)
	spinner := NewSpinner(longMessage)

	time.Sleep(100 * time.Millisecond)

	output := captureStderr(func() {
		spinner.Stop(true)
	})

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator with long message")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}
