package workspace

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

func TestWriteTextFileInDirRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		testName    string
		useFilename string
		useContent  string
	}{
		{testName: "禁止上级目录引用的文件名", useFilename: "../secret.txt", useContent: "hello"},
		{testName: "非法文件后缀", useFilename: "note.md", useContent: "hello"},
		{testName: "禁止内容为空的写入", useFilename: "note.txt", useContent: "    "},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			testDir := t.TempDir()

			gotFilePath, err := writeTextFileInDir(testDir, tc.useFilename, tc.useContent)

			if err == nil {
				t.Fatal("非法输入时，期望 writeTextFileInDir 返回错误，但实际没有错误")
			}

			if gotFilePath != "" {
				t.Errorf("期望返回的路径为空, 实际得到: %q", gotFilePath)
			}

			entries, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("调用 os.ReadDir 时异常, 错误: %v", err)
			}

			if len(entries) != 0 {
				t.Errorf("非法输入创建了文件或目录: %v", entries)
			}
		})
	}
}
