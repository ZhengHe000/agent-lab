package workspace

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestLifeCycle(t *testing.T) {
	dir := t.TempDir()

	workspace, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	if workspace == nil || workspace.root == nil {
		t.Fatal("OpenWorkspace returned an incomplete workspace")
	}
	if err = workspace.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

}

func TestInvalidDirIsRejected(t *testing.T) {
	testDir := t.TempDir()

	paths := []struct {
		name string
		path string
	}{
		{
			name: "empty",
			path: "",
		},
		{
			name: "missing",
			path: filepath.Join(testDir, "missing"),
		},
	}

	for _, tt := range paths {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OpenWorkspace(tt.path)
			if err == nil {
				t.Fatalf("want Err, got nil")
			}

			if !errors.Is(err, ErrOpenWorkspace) {
				t.Fatalf("want: %v, got: %v", ErrOpenWorkspace, err)
			}
		})
	}

}
