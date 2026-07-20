package main

import (
	"errors"
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

	if !errors.Is(err, nil) {
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

	if !errors.Is(err, nil) {
		t.Fatalf("期望错误链中包含: %v, 实际得到: %v", ErrReadFile, err)
	}
}

func createTestDirOrFile(t *testing.T, dir string, testFileName []string, testDirName []string) {
	t.Helper()

	for _, fileName := range testFileName {
		filePath := filepath.Join(dir, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			t.Errorf("文件名 %s 创建失败 文件路径: %s 错误: %v, 为保证其他内容正常创建用于测试,程序跳过本轮继续\n", fileName, filePath, err)
			continue
		}
		
		file.Close()
	}

	for _, dirName := range testDirName {
		dirPath := filepath.Join(dir, dirName)
		if err := os.Mkdir(dirPath, 0755); err != nil {
			t.Errorf("目录名 %s 创建失败 目录路径: %s 错误: %v, 为保证其他内容正常创建用于测试,程序跳过本轮继续\n", dirName, dirPath, err)
			continue
		}
	}
}

func TestListTextFilesInDirFiltersAndSorts(t *testing.T) {
	testDir := t.TempDir()

	testFileName := []string{
		"y-1.txt",  // 合法 .txt 文件
		"y-2.txt",  // 合法 .txt 文件
		"y-3.txt",  // 合法 .txt 文件
		"n-1.md",   // .md 文件
		".n-1.txt", //名字非法的 .txt 文件
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

	slices.Sort(entries)

	if !slices.Equal(entries, want) {
		t.Fatal("期待与实际得到文件不符")
	}
}
