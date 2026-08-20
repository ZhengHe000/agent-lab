package workspace

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestListTextFilesInDirMissingDir(t *testing.T) {
	testDir := t.TempDir()

	filePath := filepath.Join(testDir, "this_dir_not_exist")

	entries, err := listTextFilesInDir(filePath)

	want := []string{}

	if !slices.Equal(entries, want) {
		t.Errorf("期待填入不存在的目录时返回 []string{}, 但实际返回: %v", entries)
	}

	if err != nil {
		t.Fatalf("列出不存在目录时不应返回错误，实际得到: %v", err)
	}
}

func TestListTextFilesInDirEmptyDir(t *testing.T) {
	testDir := t.TempDir()

	entries, err := listTextFilesInDir(testDir)

	want := []string{}

	if !slices.Equal(entries, want) {
		t.Errorf("期待填入空目录时返回 []string{}, 但实际返回: %v", entries)
	}

	if err != nil {
		t.Fatalf("列出空目录时不应返回错误，实际得到: %v", err)
	}
}

func createTestDirOrFile(t *testing.T, dir string, testFileNames []string, testDirNames []string) {
	t.Helper()

	for _, fileName := range testFileNames {
		filePath := filepath.Join(dir, fileName)

		if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
			t.Fatalf("创建测试文件 %q 失败: %v", fileName, err)
		}
	}

	for _, dirName := range testDirNames {
		dirPath := filepath.Join(dir, dirName)

		if err := os.Mkdir(dirPath, 0o755); err != nil {
			t.Fatalf("创建测试目录 %q 失败: %v", dirName, err)
		}
	}
}

func TestListTextFilesInDirFiltersAndSorts(t *testing.T) {
	testDir := t.TempDir()

	testFileName := []string{
		".n-1.txt",
		"y-2.txt",
		"n-1.md",
		"y-3.txt",
		"y-1.txt",
	}

	testDirName := []string{
		"folder.txt",
	}

	createTestDirOrFile(t, testDir, testFileName, testDirName)

	entries, err := listTextFilesInDir(testDir)
	if err != nil {
		t.Fatalf("获取文件列表失败: %v", err)
	}

	want := []string{
		"y-1.txt",
		"y-2.txt",
		"y-3.txt",
	}

	if !slices.Equal(entries, want) {
		t.Fatal("期待与实际得到文件不符")
	}
}

func TestListTextFilesInDirSkipsSymlink(t *testing.T) {
	testDir := t.TempDir()

	targetPath := filepath.Join(testDir, "target.txt")
	if err := os.WriteFile(targetPath, []byte("target"), 0o644); err != nil {
		t.Fatalf("创建符号链接目标文件失败: %v", err)
	}

	linkName := "link.txt"
	linkPath := filepath.Join(testDir, linkName)

	if err := os.Symlink(targetPath, linkPath); err != nil {
		t.Skipf("当前环境不支持创建符号链接，跳过测试: %v", err)
	}

	entries, err := listTextFilesInDir(testDir)
	if err != nil {
		t.Fatalf("列出文件失败: %v", err)
	}

	want := []string{"target.txt"}

	if !slices.Equal(entries, want) {
		t.Fatalf("entries = %v, want %v", entries, want)
	}
}

func TestListTextFiles_nested_supported_files_are_returned_in_order(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"README.md":                 "readme",
		"internal/config/config.go": "package config",
		"internal/config/data.json": "{}",
		"assets/logo.png":           "image",
	}

	for path, content := range files {
		localPath := filepath.Join(
			dir,
			filepath.FromSlash(path),
		)

		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			t.Fatalf("mkdirAll %v: %v", localPath, err)
		}

		if err := os.WriteFile(localPath, []byte(content), 0o644); err != nil {
			t.Fatalf("writefile %v: %v", localPath, err)
		}
	}

	workspace, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf("open workspace %v: %v", dir, err)
	}

	t.Cleanup(func() {
		_ = workspace.Close()
	})

	got, err := workspace.ListTextFiles(context.Background())
	if err != nil {
		t.Fatalf("ListTextFiles() want nil, got: %v", err)
	}

	want := []string{
		"README.md",
		"internal/config/config.go",
		"internal/config/data.json",
	}

	if !slices.Equal(got.Paths, want) {
		t.Fatalf("ListTextFiles().Paths want %+v, got: %+v", want, got.Paths)
	}

	if got.Truncated {
		t.Fatalf("ListTextFiles().Truncated want false, got true")
	}
}
