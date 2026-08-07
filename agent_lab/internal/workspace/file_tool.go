package workspace

import (
	"bufio"
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

func writeTextFile(filename string, content string) (string, error) {
	return writeTextFileInDir(workspaceDir, filename, content)
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

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("%w: %w", ErrWriteFile, err)
	}

	return filePath, nil
}

func confirmWrite(
	reader io.Reader,
	writer io.Writer,
	filename string,
	content string,
) (bool, error) {

	const promptTemplate = `即将覆盖写入文件。

文件名：%s
内容长度：%d 字符，%d 字节
内容预览：
%s

确认写入？请输入 y 或 n:`

	characterCount := utf8.RuneCountInString(content)
	byteCount := len(content)

	prompt := fmt.Sprintf(
		promptTemplate,
		filename,
		characterCount,
		byteCount,
		content,
	)

	_, err := writer.Write([]byte(prompt))
	if err != nil {
		return false, fmt.Errorf("输出确认提示失败: %w", err)
	}

	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("读取确认输入失败: %w", err)
		}
		return false, fmt.Errorf("未读取到确认输入")
	}

	decision := strings.TrimSpace(scanner.Text())

	switch decision {
	case "y", "Y":
		return true, nil
	case "n", "N":
		return false, nil
	default:
		return false, fmt.Errorf("收到无效输入: %q, %w", decision, ErrInvalidInput)
	}

}

func confirmAndWriteTextFile(
	reader io.Reader,
	writer io.Writer,
	filename string,
	content string,
) (string, error) {
	return confirmAndWriteTextFileInDir(
		workspaceDir,
		reader,
		writer,
		filename,
		content,
	)
}

func confirmAndWriteTextFileInDir(
	dir string,
	reader io.Reader,
	writer io.Writer,
	filename string,
	content string,
) (string, error) {
	if err := validateFilename(filename); err != nil {
		return "", err
	}

	if err := validateContent(content); err != nil {
		return "", err
	}

	confirmed, err := confirmWrite(reader, writer, filename, content)
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", ErrWriteCancelled
	}

	return writeTextFileInDir(dir, filename, content)
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
