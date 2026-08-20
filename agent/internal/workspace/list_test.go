package workspace

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestListTextFiles_nested_supported_files_are_returned_in_order(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"README.md":                 "readme",
		"internal/config/config.go": "package config",
		"internal/config/data.json": "{}",
		"assets/logo.png":           "image",
	}

	for path, content := range files {
		localPath := filepath.Join(
			dir,
			filepath.FromSlash(path),
		)

		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			t.Fatalf("mkdirAll %v: %v", localPath, err)
		}

		if err := os.WriteFile(localPath, []byte(content), 0o644); err != nil {
			t.Fatalf("writefile %v: %v", localPath, err)
		}
	}

	workspace, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf("open workspace %v: %v", dir, err)
	}

	t.Cleanup(func() {
		_ = workspace.Close()
	})

	got, err := workspace.ListTextFiles(context.Background())
	if err != nil {
		t.Fatalf("ListTextFiles() want nil, got: %v", err)
	}

	want := []string{
		"README.md",
		"internal/config/config.go",
		"internal/config/data.json",
	}

	if !slices.Equal(got.Paths, want) {
		t.Fatalf("ListTextFiles().Paths want %+v, got: %+v", want, got.Paths)
	}

	if got.Truncated {
		t.Fatalf("ListTextFiles().Truncated want false, got true")
	}
}
