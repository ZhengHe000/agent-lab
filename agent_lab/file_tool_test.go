package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		filename string // 输入的文件名
		wantErr  bool   // 结果是否错误
	}{
		{filename: "note.txt", wantErr: false},
		{filename: "study-note_01.txt", wantErr: false},
		{filename: "../secret.txt", wantErr: true},
		{filename: "a/b.txt", wantErr: true},
		{filename: "my file.txt", wantErr: true},
		{filename: "report.md", wantErr: true},
		{filename: "a..b.txt", wantErr: true},
	}

	for i, tt := range tests {
		err := validateFilename(tt.filename)

		if (err != nil) != tt.wantErr {
			t.Errorf("case: %d, Filename: %s, Want: %t, Got: %t, err: %v", i, tt.filename, tt.wantErr, (err != nil), err)
		}
	}
}

func TestWriteTextFileInDir(t *testing.T) {
	testDir := t.TempDir()
	filePath, err := writeTextFileInDir(testDir, "note.txt", "hello")
	if err != nil {
		t.Fatalf("<第一次> 调用writeTextFileInDir函数时异常, 错误: %v", err)
	}

	wantPath := filepath.Join(testDir, "note.txt")
	if filePath != wantPath {
		t.Errorf("<记录> writeTextFileInDir返回的路径与期望不匹配, Want: %q, Got: %q", wantPath, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("调用os.ReadFile函数时异常, 错误: %v", err)
	}

	if got := string(data); got != "hello" {
		t.Errorf("<记录> writeTextFileInDir对测试路径下的文件写入内容与期望不匹配, Want: %s, Got: %s", "hello", got)
	}

	_, err = writeTextFileInDir(testDir, "note.txt", "updated")
	if err != nil {
		t.Fatalf("<第二次> 调用writeTextFileInDir函数时异常, 错误: %v", err)
	}

	data, err = os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("调用os.ReadFile函数时异常, 错误: %v", err)
	}

	if got := string(data); got != "updated" {
		t.Errorf("<记录> writeTextFileInDir对测试路径下的文件写入内容与期望不匹配, Want: %s, Got: %s", "updated", got)
	}
}
