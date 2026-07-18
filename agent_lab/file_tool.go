package main

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

type ByteTooLongError struct { // ByteTooLongError 字节过长错误结构体原型
	limit  int
	actual int
}

func (t *ByteTooLongError) Error() string {
	return fmt.Sprintf("内容过长, 限制 %d 字节, 当前 %d 字节", t.limit, t.actual)
}

type RuneTooLongError struct { // RuneTooLongError 字符过长错误结构体原型
	limit  int
	actual int
}

func (r *RuneTooLongError) Error() string {
	return fmt.Sprintf("内容过长, 限制 %d 字符, 当前 %d 字符", r.limit, r.actual)
}

func validateContent(content string) error { // validateContent 对Content内容校验

	if strings.TrimSpace(content) == "" { // 判断是否为空
		return ErrContentEmpty
	}

	if characterCount := utf8.RuneCountInString(content); characterCount > 10000 { // 判断字符
		return &RuneTooLongError{
			limit:  10000,
			actual: characterCount,
		}
	}

	if size := len(content); size > 40000 { // 判断字符
		return &ByteTooLongError{
			limit:  40000,
			actual: size,
		}
	}

	return nil
}

var filenameRule = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,95}\.txt$`) // filenameRule 编译可复用的规则

func validateFilename(filename string) error { // validateFilename 文件名校验

	if strings.Contains(filename, "..") { // 排除文件名存在..(两个点)的错误项
		return ErrInvalidFilename
	}

	if !filenameRule.MatchString(filename) { //用正则表达式规则判断 传入值是否合规
		return ErrInvalidFilename
	}

	return nil
}

const workspaceDir = `./AIWorkspace` // workspaceDir 是正式文件工具使用的受控工作区
func writeTextFile(filename string, content string) (string, error) { // writeTextFile 在受控工作目录中覆盖写入文本文件
	return writeTextFileInDir(workspaceDir, filename, content)
}

func writeTextFileInDir(dir string, filename string, content string) (string, error) { //  writeTextFileInDir 在指定的受控工作区中覆盖写入文本文件
	if err := validateFilename(filename); err != nil { // 检验文件名
		return "", err
	}

	if err := validateContent(content); err != nil { // 检验写入内容
		return "", err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil { // 使用os.MkdirAll在指定盘创建完整目录
		return "", fmt.Errorf("%w: %w", ErrMkdirAll, err)
	}

	filePath := filepath.Join(dir, filename) // 拼出完整文件路径

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil { // 在指定路径写入指定内容
		return "", fmt.Errorf("%w: %w", ErrWriteFile, err)
	}

	return filePath, nil // 返回文件路径和空
}

func confirmWrite(
	reader io.Reader,
	writer io.Writer,
	filename string,
	content string,
) (bool, error) {

	//设置模板
	const promptTemplate = `即将覆盖写入文件。

文件名：%s
内容长度：%d 字符，%d 字节
内容预览：
%s

确认写入？请输入 y 或 n:`

	characterCount := utf8.RuneCountInString(content) // 计算字符数
	byteCount := len(content)                         // 计算字节数

	prompt := fmt.Sprintf( // 装填模板参数
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

	scanner := bufio.NewScanner(reader) // 创建扫描器

	if !scanner.Scan() { // 阻塞程序直到得到输入 + 错误处理
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("读取确认输入失败: %w", err)
		}
		return false, fmt.Errorf("未读取到确认输入")
	}

	decision := strings.TrimSpace(scanner.Text()) // 接收输入

	switch decision { // 筛选输入
	case "y", "Y":
		return true, nil
	case "n", "N":
		return false, nil
	default:
		return false, fmt.Errorf("无效确认输入 %q,请输入 y 或 n", decision)
	}

}

func confirmAndWriteTextFile( //组装file_tool使用流程
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
	if err := validateFilename(filename); err != nil { // 使用文件名校验
		return "", err
	}

	if err := validateContent(content); err != nil { // 使用文件内容校验
		return "", err
	}

	confirmed, err := confirmWrite(reader, writer, filename, content) // 进行确认
	if err != nil {
		return "", err
	}

	if !confirmed { // 拒绝或异常时中断函数
		return "", ErrWriteCancelled
	}

	return writeTextFileInDir(dir, filename, content) // 调用写入函数, 返回string 和 error
}
