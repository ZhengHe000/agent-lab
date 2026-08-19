package workspace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func openTestWorkspace(t *testing.T, dir string) *Workspace {
	t.Helper()

	workspace, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf("Helper_OpenWorkspace() create workspace failed, %v: %v", dir, err)
	}

	t.Cleanup(func() {
		_ = workspace.Close()
	})

	return workspace
}

func TestReadTextFileTool_definition_exposes_path_schema(t *testing.T) {
	testDir := t.TempDir()

	workspace := openTestWorkspace(t, testDir)

	readTool, err := NewReadTextFileTool(workspace)
	if err != nil {
		t.Fatalf("got ReadTool failed: %v", err)
	}

	definition := readTool.Definition()
	if definition.Name != "read_text_file" {
		t.Fatalf("want: %v, got: %v", "read_text_file", definition.Name)
	}
	var schema struct {
		Required   []string                   `json:"required"`
		Properties map[string]json.RawMessage `json:"properties"`
	}

	if err := json.Unmarshal(definition.Parameters, &schema); err != nil {
		t.Fatalf("parsed definition.Parameters failed: %v", err)
	}

	if _, exists := schema.Properties["path"]; !exists {
		t.Fatalf("want true, got %v", exists)
	}

	// 用于检查确认第一次工具升级时淘汰的字段"filename"不存在, 注:filename被替换为path.
	if _, exists := schema.Properties["filename"]; exists {
		t.Fatalf("want false, got %v", exists)
	}

	if !slices.Contains(schema.Required, "path") {
		t.Fatalf("want required comprise path arguments, but not")
	}
}

func TestReadTextFileTool_nested_path_returns_content(t *testing.T) {
	dir := t.TempDir()
	nestedDir := filepath.Join(dir, "internal", "config")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v, want nil", nestedDir, err)
	}

	filePath := filepath.Join(nestedDir, "config.go")
	if err := os.WriteFile(filePath, []byte("package config"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v, want nil", filePath, err)
	}

	projectWorkspace := openTestWorkspace(t, dir)

	readTool, err := NewReadTextFileTool(projectWorkspace)
	if err != nil {
		t.Fatalf("NewReadTextFileTool() error = %v, want nil", err)
	}

	arguments := json.RawMessage(
		`{"path":"internal/config/config.go"}`,
	)

	got, err := readTool.Execute(context.Background(), arguments)
	if err != nil {
		t.Fatalf("Execute(%s) error = %v, want nil", arguments, err)
	}

	if got != "package config" {
		t.Fatalf("Execute(%s) = %q, want %q", arguments, got, "package config")
	}
}

func TestReadTextFileTool_invalid_json_returns_error(t *testing.T) {
	projectWorkspace := openTestWorkspace(t, t.TempDir())

	readTool, err := NewReadTextFileTool(projectWorkspace)
	if err != nil {
		t.Fatalf(
			"NewReadTextFileTool() error = %v, want nil",
			err,
		)
	}

	arguments := json.RawMessage(`{"path":`)

	got, err := readTool.Execute(context.Background(), arguments)
	if err == nil {
		t.Fatalf(
			"Execute(%s) error = nil, want non-nil",
			arguments,
		)
	}

	if got != "" {
		t.Fatalf(
			"Execute(%s) = %q, want empty string",
			arguments,
			got,
		)
	}
}
