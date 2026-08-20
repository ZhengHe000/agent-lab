package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/model"
	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/tool"
)

// ListTextFilesTool 列出受控工作区中的可读取文本文件。
type ListTextFilesTool struct {
	workspace *Workspace
}

// NewListTextFilesTool 创建使用正式工作区的文件列表工具。
func NewListTextFilesTool(workspace *Workspace) (*ListTextFilesTool, error) {
	if workspace == nil || workspace.root == nil {
		return nil, fmt.Errorf("create list_text_file tool: workspace is nil")
	}

	return &ListTextFilesTool{
		workspace: workspace,
	}, nil
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
		return "", fmt.Errorf("context canceled before tool execution: %w", err)
	}

	var args struct{}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("parse list_text_files arguments: %w", err)
	}

	got, err := l.workspace.ListTextFiles(ctx)
	if err != nil {
		return "", fmt.Errorf("execute ListTextFiles(ctx): %w", err)
	}

	result, err := json.Marshal(got)
	if err != nil {
		return "", fmt.Errorf("marshal ListTextFiles(ctx).got: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before result returned: %w", err)
	}

	return string(result), nil
}

var _ tool.Tool = (*ListTextFilesTool)(nil)
