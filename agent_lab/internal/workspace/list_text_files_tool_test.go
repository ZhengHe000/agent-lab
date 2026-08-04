package workspace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestListTextFilesTool_Definition(t *testing.T) {
	testDir := t.TempDir()

	testTool := newListTextFilesToolInDir(testDir)

	testDfn := testTool.Definition()

	if testDfn.Name != "list_text_files" {
		t.Fatalf("want: %s, got: %s", "list_text_files", testDfn.Name)
	}

	if !json.Valid(testDfn.Parameters) {
		t.Fatalf("JSON格式错误: %v", testDfn.Parameters)
	}
}
func TestListTextFilesTool_Execute_Success(t *testing.T) {
	testDir := t.TempDir()

	testPaths := []string{
		filepath.Join(testDir, "note.txt"),
		filepath.Join(testDir, "plan.txt"),
		filepath.Join(testDir, "ignore.md"),
	}

	wantFiles := `["note.txt","plan.txt"]`

	for _, testPath := range testPaths {
		if err := os.WriteFile(testPath, []byte("test"), 0o644); err != nil {
			t.Fatalf("创建测试文件 %s 失败: %v", testPath, err)
		}
	}

	testTool := newListTextFilesToolInDir(testDir)

	args := json.RawMessage(`{}`)
	result, err := testTool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("执行Tool函数失败: %v", err)
	}

	if result != wantFiles {
		t.Fatalf("want: %v, got: %v", wantFiles, result)
	}
}
