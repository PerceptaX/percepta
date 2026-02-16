package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewSpinner(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Testing...", &buf)
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
	var buf bytes.Buffer
	spinner := newSpinner("Operation", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.Stop(true)
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator ✓")
	}

	if !strings.Contains(output, "Operation") {
		t.Error("Expected message in output")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_Stop_Failure(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Operation", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.Stop(false)
	output := buf.String()

	if !strings.Contains(output, "✗") {
		t.Error("Expected failure indicator ✗")
	}

	if !strings.Contains(output, "Operation") {
		t.Error("Expected message in output")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_StopWithMessage_Success(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Loading", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.StopWithMessage(true, "Completed successfully")
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator ✓")
	}

	if !strings.Contains(output, "Completed successfully") {
		t.Error("Expected custom message in output")
	}
}

func TestSpinner_StopWithMessage_Failure(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Processing", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.StopWithMessage(false, "Failed to process")
	output := buf.String()

	if !strings.Contains(output, "✗") {
		t.Error("Expected failure indicator ✗")
	}

	if !strings.Contains(output, "Failed to process") {
		t.Error("Expected custom message in output")
	}
}

func TestSpinner_DoubleStop(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Test", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.Stop(true)

	if !strings.Contains(buf.String(), "✓") {
		t.Error("Expected success indicator on first stop")
	}

	// Record buffer length after first stop
	lenAfterFirst := buf.Len()

	// Second stop should be a no-op
	spinner.Stop(false)

	if buf.Len() != lenAfterFirst {
		t.Error("Expected no additional output on second stop")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to remain stopped")
	}
}

func TestSpinner_DoubleStopWithMessage(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Test", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.StopWithMessage(true, "Done")

	if !strings.Contains(buf.String(), "Done") {
		t.Error("Expected custom message on first stop")
	}

	lenAfterFirst := buf.Len()

	// Second stop should be a no-op
	spinner.StopWithMessage(false, "Failed")

	if buf.Len() != lenAfterFirst {
		t.Error("Expected no additional output on second stop with message")
	}
}

func TestSpinner_ImmediateStop(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Quick test", &buf)

	spinner.Stop(true)
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator even for immediate stop")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_Lifecycle(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Lifecycle test", &buf)

	if spinner.stopped {
		t.Error("Expected spinner to not be stopped initially")
	}

	time.Sleep(200 * time.Millisecond)

	spinner.Stop(true)
	output := buf.String()

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator in output")
	}
}

func TestSpinner_MultipleSpinners(t *testing.T) {
	var buf1, buf2, buf3 bytes.Buffer
	spinner1 := newSpinner("Task 1", &buf1)
	spinner2 := newSpinner("Task 2", &buf2)
	spinner3 := newSpinner("Task 3", &buf3)

	time.Sleep(150 * time.Millisecond)

	spinner2.Stop(true)
	time.Sleep(50 * time.Millisecond)
	spinner1.Stop(false)
	time.Sleep(50 * time.Millisecond)
	spinner3.Stop(true)

	if !spinner1.stopped || !spinner2.stopped || !spinner3.stopped {
		t.Error("Expected all spinners to be stopped")
	}
}

func TestSpinner_LongRunning(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("Long operation", &buf)

	time.Sleep(500 * time.Millisecond)

	spinner.Stop(true)
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator after long run")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	spinner := newSpinner("", &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.Stop(true)
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator even with empty message")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}

func TestSpinner_LongMessage(t *testing.T) {
	var buf bytes.Buffer
	longMessage := strings.Repeat("Very long message ", 20)
	spinner := newSpinner(longMessage, &buf)

	time.Sleep(100 * time.Millisecond)

	spinner.Stop(true)
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Expected success indicator with long message")
	}

	if !spinner.stopped {
		t.Error("Expected spinner to be stopped")
	}
}
