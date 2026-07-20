package main

import (
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
		t.Fatalf("期望错误链中包含: %v, 实际得到: %v", ErrReadFile, err)
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
		t.Fatalf("期望错误链中包含: %v, 实际得到: %v", ErrReadFile, err)
	}
}

func createTestDirOrFile(t *testing.T, dir string, testFileNames []string, testDirNames []string) {
	t.Helper()

	target := "target.txt" 
	targetPath := filepath.Join(dir, target)

	if err := os.WriteFile(targetPath, []byte("target"), 0o644); err != nil {
		t.Fatalf("创建符号链接目标文件失败: %v", err)
	}

	link := "link.txt"
	linkPath := filepath.Join(dir, link)

	if err := os.Symlink(targetPath, linkPath); err != nil {
		t.Skipf("当前环境不支持创建符号链接，跳过测试: %v", err)
	}

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

	testFileName := []string{ // 混合创建顺序
		".n-1.txt", //名字非法的 .txt 文件
		"y-2.txt",  // 合法 .txt 文件
		"n-1.md",   // .md 文件
		"y-3.txt",  // 合法 .txt 文件
		"y-1.txt",  // 合法 .txt 文件
	}

	testDirName := []string{
		"folder.txt", // 名为 folder.txt 的目录
	}

	createTestDirOrFile(t, testDir, testFileName, testDirName)

	entries, err := listTextFilesInDir(testDir)
	if err != nil {
		t.Fatalf("获取文件列表失败: %v", err)
	}

	want := []string{ // 期望entries这个切片中出现的内容
		"y-1.txt",
		"y-2.txt",
		"y-3.txt",
	}

	if !slices.Equal(entries, want) {
		t.Fatal("期待与实际得到文件不符")
	}
}
