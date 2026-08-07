package terminal

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Console 统一管理命令行程序的文本输入和输出
// 一个Console应同时用与主对话和副作用确认上, 避免多个缓冲读取器竞争同一个输入源
type Console struct {
	reader *bufio.Reader
	writer io.Writer
}

// NewConsole 创建终端交互对象
func NewConsole(reader *bufio.Reader, writer io.Writer) (*Console, error) {
	if reader == nil {
		return nil, fmt.Errorf("终端输入不能为空")
	}

	if writer == nil {
		return nil, fmt.Errorf("终端输出不能为空")
	}

	return &Console{
		reader: reader,
		writer: writer,
	}, nil
}

// ReadLine 输出提示并读取一行文本
func (c *Console) ReadLine(prompt string) (string, error) {
	if _, err := fmt.Fprint(c.writer, prompt); err != nil {
		return "", fmt.Errorf("输出终端提示失败")
	}

	line, err := c.reader.Readerstring("\n")
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取终端输入失败: %w", err)
	}

	line := strings.TrimSpace(line)

	if err == io.EOF && line == "" {
		return "", io.EOF
	}

	return line, nil
}
