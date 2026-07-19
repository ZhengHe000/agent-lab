package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadTextFileInDir(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		setup    func(t *testing.T, dir string, filename string, content string)
		want     string
		wantErr  error
	}{
		{
			name:     "正常读取",
			filename: "note.txt",
			content:  "hello",
			setup: func(t *testing.T, dir string, filename string, content string) {
				t.Helper()

				filePath := filepath.Join(dir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
					t.Fatalf("测试文件创建失败, Err: %v", err)
				}
			},
			want:    "hello",
			wantErr: nil,
		},
		{
			name:     "非法文件名",
			filename: "../secret.txt",
			want:     "",
			wantErr:  ErrInvalidFilename,
		},
		{
			name:     "文件不存在",
			filename: "missing.txt",
			want:     "",
			wantErr:  ErrReadFile,
		},
		{
			name:     "文件刚好达到上限",
			filename: "note.txt",
			content:  strings.Repeat("a", maxReadBytes),
			setup: func(t *testing.T, dir string, filename string, content string) {
				t.Helper()

				filePath := filepath.Join(dir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
					t.Fatalf("测试文件创建失败, Err: %v", err)
				}
			},
			want:    strings.Repeat("a", maxReadBytes),
			wantErr: nil,
		},
		{
			name:     "文件超过上限",
			filename: "note.txt",
			content:  strings.Repeat("a", maxReadBytes+1),
			setup: func(t *testing.T, dir string, filename string, content string) {
				t.Helper()

				filePath := filepath.Join(dir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
					t.Fatalf("测试文件创建失败, Err: %v", err)
				}
			},
			want:    "",
			wantErr: ErrFileTooLarge,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()

			if tc.setup != nil {
				tc.setup(t, testDir, tc.filename, tc.content)
			}

			gotStr, gotErr := readTextFileInDir(testDir, tc.filename)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("wantErr: %v, gotErr: %v", tc.wantErr, gotErr)
			}

			if gotStr != tc.want {
				t.Fatalf("want: %v, gotStr: %v", tc.want, gotStr)
			}
		})
	}
}
