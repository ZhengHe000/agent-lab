package workspace

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadTextFile_nested_allowed_file_is_read(t *testing.T) {
	testDir := t.TempDir()

	nestedDir := filepath.Join(testDir, "internal", "config")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("MkdirAll %q: %v, want nil", nestedDir, err)
	}

	filePath := filepath.Join(nestedDir, "config.go")
	if err := os.WriteFile(filePath, []byte("package config"), 0o644); err != nil {
		t.Fatalf("write file %q: %v, want nil", filePath, err)
	}

	workspace, err := OpenWorkspace(testDir)
	if err != nil {
		t.Fatalf("OpenWorkspace (%q): %v, want nil", testDir, err)
	}
	t.Cleanup(func() {
		_ = workspace.Close()
	})

	got, err := workspace.ReadTextFile("internal/config/config.go")
	if err != nil {
		t.Fatalf("read file %q: %v, want nil", filePath, err)
	}

	if got != "package config" {
		t.Fatalf("got: %v, want: %v", got, "package config")
	}
}

func TestReadTextFile_symlink_parent_is_rejected(t *testing.T) {
	workspaceDir := t.TempDir()
	outsideDir := t.TempDir()

	targetPath := filepath.Join(outsideDir, "secret.go")
	if err := os.WriteFile(targetPath, []byte("secret"), 0o644); err != nil {
		t.Fatalf("write file %q: %v, want nil", targetPath, err)
	}

	linkPath := filepath.Join(workspaceDir, "config")
	if err := os.Symlink(outsideDir, linkPath); err != nil {
		t.Skipf("os.Symlink(): %v", err)
	}

	workspace, err := OpenWorkspace(workspaceDir)
	if err != nil {
		t.Fatalf("OpenWorkspace (%q): %v, want nil", workspaceDir, err)
	}
	t.Cleanup(func() {
		_ = workspace.Close()
	})

	_, err = workspace.ReadTextFile("config/secret.go")
	if !errors.Is(err, ErrSymlinkPath) {
		t.Fatalf("want errors.Is(err, ErrSymlinkPath), but got: %v", err)
	}
}

func TestReadTextFile_oversized_file_is_rejected(t *testing.T) {
	dir := t.TempDir()
	content := bytes.Repeat([]byte("a"), maxReadBytes+1)

	filePath := filepath.Join(dir, "large.txt")
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v, want nil", filePath, err)
	}

	workspace, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf("OpenWorkspace(%q) error = %v, want nil", dir, err)
	}
	t.Cleanup(func() {
		_ = workspace.Close()
	})

	_, err = workspace.ReadTextFile("large.txt")
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("want errors.Is(err, ErrFileTooLarge), but got: %v", err)
	}
}
