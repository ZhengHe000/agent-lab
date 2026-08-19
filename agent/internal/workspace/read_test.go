package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"bytes"
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
		{
			name:     "拒绝读取目录",
			filename: "folder.txt",
			setup: func(t *testing.T, dir string, filename string, content string) {
				t.Helper()

				dirPath := filepath.Join(dir, filename)
				if err := os.Mkdir(dirPath, 0o755); err != nil {
					t.Fatalf("创建测试目录失败: %v", err)
				}
			},
			want:    "",
			wantErr: ErrReadFile,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()

			if tc.setup != nil {
				tc.setup(t, testDir, tc.filename, tc.content)
			}

			got, gotErr := readTextFileInDir(testDir, tc.filename)

			if got != tc.want {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}

			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("wantErr: %v, gotErr: %v", tc.wantErr, gotErr)
			}
		})
	}
}

func TestReadTextFileInDirRejectsSymlink(t *testing.T) {
	testDir := t.TempDir()

	targetPath := filepath.Join(testDir, "target.txt")
	if err := os.WriteFile(targetPath, []byte("secret"), 0o644); err != nil {
		t.Fatalf("创建符号链接目标文件失败: %v", err)
	}

	linkName := "link.txt"
	linkPath := filepath.Join(testDir, linkName)

	if err := os.Symlink(targetPath, linkPath); err != nil {
		t.Skipf("当前环境不支持创建符号链接，跳过测试: %v", err)
	}

	got, err := readTextFileInDir(testDir, linkName)

	if got != "" {
		t.Errorf("拒绝读取符号链接时应返回空内容，实际得到: %q", got)
	}

	if !errors.Is(err, ErrReadFile) {
		t.Fatalf("Err: %v, 期望错误链包含: %v", err, ErrReadFile)
	}
}

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
