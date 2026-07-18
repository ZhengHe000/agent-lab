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

	if len(strings.TrimSpace(filename)) == 0 || len(strings.TrimSpace(content)) == 0 { // 确认后续参数非空
		return false, ErrContentEmpty
	}

	//设置模板
	const template = `
	|即将覆盖写入文件		  
	|文件名: %s
	|内容长度: %d字节. %d 字符  	
	|内容预览:
	|%v|
	|确认写入？请输入 y(确认) 或 n(拒绝):`

	charCount := utf8.RuneCountInString(content) // 计算字符数
	byteCount := len(content)                    // 计算字节数

	prompt := fmt.Sprintf(template, filename, byteCount, charCount, content) // 装填模板参数
	_, err := writer.Write([]byte(prompt))
	if err != nil {
		return false, fmt.Errorf("错误, 使用writer时异常, Err: %v", err)
	}

	scanners := bufio.NewScanner(reader) // 创建扫描器

	for scanners.Scan() { // 阻塞程序直到得到输入

		decision := strings.TrimSpace(scanners.Text()) // 接收输入

		if len(decision) == 0 { // 检验
			return false, fmt.Errorf("禁止输入空内容, 必须输入y 或 n 来做选择")
		}

		if decision == "y" || decision == "Y" { // 按照正确输入做决策
			return true, nil
		}

		if decision == "n" || decision == "N" { // 按照正确输入做决策
			return false, nil
		}

		return false, fmt.Errorf("其他异常, 接收到输入为: %v", decision) // 未知异常直接拒绝文件写入
	}

	if err := scanners.Err(); err != nil { // 检验扫描结束的正确性
		return false, fmt.Errorf("错误, Scanner异常结束, Err: %v", err)
	}

	return false, fmt.Errorf("函数: confirmWrite 出现未知问题导致执行到此, 此处为兜底默认拒绝执行策略") // 默认拒绝行为
}
