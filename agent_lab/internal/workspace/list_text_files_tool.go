package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
)

type ListTextFilesTool struct { // 列出受控工作区全部文件的Tool
	dir string
}

func NewListTextFilesTool() *ListTextFilesTool { // 受控工作区
	return &ListTextFilesTool{
		dir: workspaceDir,
	}
}

func newListTextFilesToolInDir(dir string) *ListTextFilesTool { // 测试使用可注入目录
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
