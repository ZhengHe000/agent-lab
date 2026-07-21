package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfirmAndWriteTextFileInDir(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		filename   string
		content    string
		wantErr    error
		wantFile   bool
		wantOutput string
	}{
		{
			name:       "确认后写入文件",
			input:      "y\n",
			filename:   "note.txt",
			content:    "hello",
			wantErr:    nil,
			wantFile:   true,
			wantOutput: "确认写入",
		},
		{
			name:       "拒绝后不写入文件",
			input:      "n\n",
			filename:   "note.txt",
			content:    "hello",
			wantErr:    ErrWriteCancelled,
			wantFile:   false,
			wantOutput: "确认写入",
		},
		{
			name:       "无效确认输入不写入文件",
			input:      "maybe\n",
			filename:   "note.txt",
			content:    "hello",
			wantErr:    ErrInvalidInput,
			wantFile:   false,
			wantOutput: "确认写入",
		},
		{
			name:       "非法文件名不显示确认提示也不写入",
			input:      "y\n",
			filename:   "../secret.txt",
			content:    "hello",
			wantErr:    ErrInvalidFilename,
			wantFile:   false,
			wantOutput: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()

			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			filePath, gotErr := confirmAndWriteTextFileInDir(testDir, input, &output, tc.filename, tc.content)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("实际得到的错误链中不包含预期的原始错误, want: %v, got: %v", tc.wantErr, gotErr)
			}
			if tc.wantOutput == "" {
				if output.Len() != 0 {
					t.Errorf("校验失败时不应显示确认提示，实际输出: %q", output.String())
				}
			} else if !strings.Contains(output.String(), tc.wantOutput) {
				t.Errorf("确认提示未包含 %q, 实际输出: %q", tc.wantOutput, output.String())
			}

			entries, readErr := os.ReadDir(testDir)
			if readErr != nil {
				t.Fatalf("读取临时目录失败: %v", readErr)
			}

			if !tc.wantFile {
				if filePath != "" {
					t.Errorf("未写入时应返回空路径，实际得到 %q", filePath)
				}
				if len(entries) != 0 {
					t.Errorf("未确认或输入无效时不应创建文件，实际目录: %v", entries)
				}
				return
			}

			wantPath := filepath.Join(testDir, tc.filename)
			if wantPath != filePath {
				t.Errorf("filePath = %q, want %q", filePath, wantPath)
			}
			data, readErr := os.ReadFile(filePath)
			if readErr != nil {
				t.Fatalf("读取已写入文件失败: %v", readErr)
			}

			if got := string(data); got != tc.content {
				t.Errorf("文件内容 = %q, want %q", got, tc.content)
			}

		})
	}
}
