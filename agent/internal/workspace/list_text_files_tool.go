package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZhengHeOwo/agent-AuXuan/agent/internal/model"
	"github.com/ZhengHeOwo/agent-AuXuan/agent/internal/tool"
)

// ListTextFilesTool 列出受控工作区中的可读取文本文件。
type ListTextFilesTool struct {
	dir string
}

// NewListTextFilesTool 创建使用正式工作区的文件列表工具。
func NewListTextFilesTool() *ListTextFilesTool {
	return &ListTextFilesTool{
		dir: workspaceDir,
	}
}

func newListTextFilesToolInDir(dir string) *ListTextFilesTool {
	return &ListTextFilesTool{
		dir: dir,
	}
}

var listTextFilesParameters = json.RawMessage(`{
	"type": "object",
	"properties": {},
	 "additionalProperties": false
}`)

func (l *ListTextFilesTool) Definition() model.ToolDefinition {
	return model.ToolDefinition{
		Name: "list_text_files",
		Description: "列出受控工作区中所有可读取的文本文件名." +
			"当用户没有给出准确文件名，或需要了解工作区有哪些文件时使用.",
		Parameters: listTextFilesParameters,
	}
}

func (l *ListTextFilesTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("列出文件前上下文已结束: %w", err)
	}

	var args struct{}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("解析 list_text_files 参数失败: %w", err)
	}

	files, err := listTextFilesInDir(l.dir)
	if err != nil {
		return "", fmt.Errorf("执行 list_text_files 失败: %w", err)
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", fmt.Errorf("编码 list_text_files 执行结果失败: %w", err)
	}

	return string(result), nil
}

var _ tool.Tool = (*ListTextFilesTool)(nil)
