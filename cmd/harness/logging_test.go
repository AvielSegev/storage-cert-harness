package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLoggerToMirrorsHarnessMessages(t *testing.T) {
	var terminal, runLog bytes.Buffer
	logger := newLoggerTo(false, &terminal, &runLog)
	logger.Info("message", "component", "test")

	if !strings.Contains(terminal.String(), "message") {
		t.Fatalf("terminal output does not contain log message: %q", terminal.String())
	}
	if !strings.Contains(runLog.String(), "message") {
		t.Fatalf("run log does not contain log message: %q", runLog.String())
	}
}
