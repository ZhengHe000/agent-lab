package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteTextFile_creates_and_replaces_nested_file(t *testing.T) {
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

	toolPath := "internal/config/config.go"

	if err := projectWorkspace.WriteTextFile(
		toolPath,
		"package config",
	); err != nil {
		t.Fatalf(
			"WriteTextFile() create error = %v, want nil",
			err,
		)
	}

	if err := projectWorkspace.WriteTextFile(
		toolPath,
		"package updated",
	); err != nil {
		t.Fatalf(
			"WriteTextFile() replace error = %v, want nil",
			err,
		)
	}

	localPath := filepath.Join(
		dir,
		filepath.FromSlash(toolPath),
	)

	got, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatalf(
			"os.ReadFile(%q) error = %v, want nil",
			localPath,
			err,
		)
	}

	if string(got) != "package updated" {
		t.Fatalf(
			"os.ReadFile(%q) = %q, want %q",
			localPath,
			string(got),
			"package updated",
		)
	}
}
