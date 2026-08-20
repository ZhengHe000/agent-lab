package workspace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestListTextFilesTool_returns_structured_file_list(t *testing.T) {
	dir := t.TempDir()

	localPath := filepath.Join(
		dir,
		"internal",
		"config",
		"config.go",
	)

	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		t.Fatalf(
			"os.MkdirAll(%q) error = %v, want nil",
			filepath.Dir(localPath),
			err,
		)
	}

	if err := os.WriteFile(
		localPath,
		[]byte("package config"),
		0o644,
	); err != nil {
		t.Fatalf(
			"os.WriteFile(%q) error = %v, want nil",
			localPath,
			err,
		)
	}

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

	listTool, err := NewListTextFilesTool(projectWorkspace)
	if err != nil {
		t.Fatalf(
			"NewListTextFilesTool() error = %v, want nil",
			err,
		)
	}

	encoded, err := listTool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf(
			"Execute() error = %v, want nil",
			err,
		)
	}

	var response listTextFilesResponse
	if err := json.Unmarshal([]byte(encoded), &response); err != nil {
		t.Fatalf(
			"json.Unmarshal() error = %v, want nil",
			err,
		)
	}

	wantPaths := []string{
		"internal/config/config.go",
	}

	if !slices.Equal(response.Paths, wantPaths) {
		t.Fatalf(
			"response.Paths = %v, want %v",
			response.Paths,
			wantPaths,
		)
	}

	if response.Truncated {
		t.Fatal("response.Truncated = true, want false")
	}
}
