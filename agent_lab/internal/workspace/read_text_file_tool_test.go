package workspace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadTextFileTool_Definition(t *testing.T) {
	testDir := t.TempDir()

	testTool := newReadTextFileToolInDir(testDir)

	testDfn := testTool.Definition()

	if testDfn.Name != "read_text_file" {
		t.Fatalf("want: %s, got: %s", "read_text_file", testDfn.Name)
	}

	if !json.Valid(testDfn.Parameters) {
		t.Fatalf("JSON格式错误: %v", testDfn.Parameters)
	}
}

func TestReadTextFileTool_Execute_Success(t *testing.T) {
	testDir := t.TempDir()
	testPath := filepath.Join(testDir, "target.txt")

	if err := os.WriteFile(testPath, []byte("zero"), 0o644); err != nil {
		t.Fatalf("准备测试文件失败: %v", err)
	}

	testArgs := json.RawMessage(`{"filename":"target.txt"}`)

	testTool := newReadTextFileToolInDir(testDir)

	resp, err := testTool.Execute(context.Background(), testArgs)
	if err != nil {
		t.Fatalf("执行Tool函数失败: %v", err)
	}

	if resp != "zero" {
		t.Fatalf("获取的文件内容与期望的 %s 不符, 实际: %s", "zero", resp)
	}
}

func TestReadTextFileTool_Execute_InvalidJSON(t *testing.T) {
	testDir := t.TempDir()
	
	testArgs := json.RawMessage(`{"filename":"//\\ww.dd"`)

	testTool := newReadTextFileToolInDir(testDir)

	resp, err := testTool.Execute(context.Background(), testArgs)
	if err == nil {
		t.Fatalf("执行Tool函数失败时得到的err为nil")
	}

	if resp != "" {
		t.Fatalf("获取文件错误时应获得空, 实际: %s", resp)
	}
}
