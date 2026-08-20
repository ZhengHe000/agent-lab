package workspace

import (
	"fmt"
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
