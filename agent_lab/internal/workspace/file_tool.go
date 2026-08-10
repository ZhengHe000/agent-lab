package workspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

type ByteTooLongError struct {
	limit  int
	actual int
}

func (t *ByteTooLongError) Error() string {
	return fmt.Sprintf("内容过长, 限制 %d 字节, 当前 %d 字节", t.limit, t.actual)
}

type RuneTooLongError struct {
	limit  int
	actual int
}

func (r *RuneTooLongError) Error() string {
	return fmt.Sprintf("内容过长, 限制 %d 字符, 当前 %d 字符", r.limit, r.actual)
}

const ( // 默认条件

	// workspaceDir 受控工作区。
	workspaceDir = `./AIWorkspace`

	// maxReadBytes 单次的最大字节数。
	maxReadBytes = 40_000
)

func validateContent(content string) error { // validateContent 判断内容合法性

	if strings.TrimSpace(content) == "" {
		return ErrContentEmpty
	}

	if characterCount := utf8.RuneCountInString(content); characterCount > 10000 {
		return &RuneTooLongError{
			limit:  10000,
			actual: characterCount,
		}
	}

	if size := len(content); size > 40000 {
		return &ByteTooLongError{
			limit:  40000,
			actual: size,
		}
	}

	return nil
}

var filenameRule = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,95}\.txt$`)

func validateFilename(filename string) error { // validateFilename 判断文件名合法性

	if strings.Contains(filename, "..") {
		return ErrInvalidFilename
	}

	if !filenameRule.MatchString(filename) {
		return ErrInvalidFilename
	}

	return nil
}

func writeTextFileInDir(dir string, filename string, content string) (string, error) {
	if err := validateFilename(filename); err != nil {
		return "", err
	}

	if err := validateContent(content); err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("%w: %w", ErrMkdirAll, err)
	}

	filePath := filepath.Join(dir, filename)

	info, err := os.Lstat(filePath)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("%w: 拒绝写入链接文件", ErrWriteFile)
		}

		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("%w: 只能覆盖普通文件", ErrWriteFile)
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("%w: 获取目标文件信息失败: %w", ErrWriteFile, err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("%w: %w", ErrWriteFile, err)
	}

	return filePath, nil
}

func readTextFile(filename string) (string, error) {
	return readTextFileInDir(workspaceDir, filename)
}

func readTextFileInDir(dir string, filename string) (string, error) {
	if err := validateFilename(filename); err != nil {
		return "", err
	}

	filePath := filepath.Join(dir, filename)

	info, err := os.Lstat(filePath)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrReadFile, err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("%w: 拒绝读取符号链接", ErrReadFile)
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%w: 只能读取普通文件", ErrReadFile)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrReadFile, err)
	}
	defer file.Close()

	limitedReader := io.LimitReader(file, int64(maxReadBytes+1))

	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrReadFile, err)
	}

	if len(data) > maxReadBytes {
		return "", ErrFileTooLarge
	}

	return string(data), nil
}

func listTextFiles() ([]string, error) {
	return listTextFilesInDir(workspaceDir)
}

func listTextFilesInDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("%w:获取目录内容失败: %w", ErrReadFile, err)
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if err := validateFilename(entry.Name()); err != nil {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("%w:获取文件元数据失败: %w", ErrReadFile, err)
		}

		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}

		if !info.Mode().IsRegular() {
			continue
		}

		files = append(files, entry.Name())
	}

	return files, nil
}
