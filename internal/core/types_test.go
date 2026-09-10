package core

import (
	"bytes"
	"testing"
)

func TestToolOutputMirrorsSubprocessStreams(t *testing.T) {
	var terminal, runLog bytes.Buffer
	rc := &RunCtx{LogOutput: &runLog}
	if _, err := rc.ToolOutput(&terminal).Write([]byte("tool output\n")); err != nil {
		t.Fatal(err)
	}

	if terminal.String() != runLog.String() {
		t.Fatalf("terminal and run log differ: terminal=%q run-log=%q", terminal.String(), runLog.String())
	}
}
