package workspace

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/tool"
)

type recordingConfirmer struct {
	calls int
}

func (c *recordingConfirmer) Confirm(
	ctx context.Context,
	request tool.ConfirmationRequest,
) (bool, error) {
	c.calls++
	return true, nil
}

func TestWriteTextFileTool_invalid_path_is_rejected_before_confirmation(
	t *testing.T,
) {
	dir := t.TempDir()

	projectWorkspace, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf(
			"OpenWorkspace(%q) error = %v, want nil",
			dir,
			err,
		)
	}
	t.Cleanup(func() {
		_ = projectWorkspace.Close()
	})

	confirmer := &recordingConfirmer{}

	writeTool, err := NewWriteTextFileTool(
		projectWorkspace,
		confirmer,
	)
	if err != nil {
		t.Fatalf(
			"NewWriteTextFileTool() error = %v, want nil",
			err,
		)
	}

	arguments := json.RawMessage(
		`{"path":"../secret.txt","content":"secret"}`,
	)

	result, err := writeTool.Execute(
		context.Background(),
		arguments,
	)
	if err == nil {
		t.Fatalf(
			"Execute(%s) error = nil, want non-nil",
			arguments,
		)
	}

	if result != "" {
		t.Fatalf(
			"Execute(%s) result = %q, want empty string",
			arguments,
			result,
		)
	}

	if confirmer.calls != 0 {
		t.Fatalf(
			"Confirm() calls = %d, want 0",
			confirmer.calls,
		)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf(
			"os.ReadDir(%q) error = %v, want nil",
			dir,
			err,
		)
	}

	if len(entries) != 0 {
		t.Fatalf(
			"os.ReadDir(%q) returned %d entries, want 0",
			dir,
			len(entries),
		)
	}
}
