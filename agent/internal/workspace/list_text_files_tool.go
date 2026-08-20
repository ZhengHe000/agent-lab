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
		return nil, fmt.Errorf("create list_text_files tool: workspace is nil")
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
		Description: "Recursively list workspace-relative paths for supported regular text files " +
			"in the controlled workspace. Use this tool to discover available files when the exact " +
			"path is unknown or before reading, searching, or editing files. The result indicates " +
			"whether the listing was truncated by safety limits.",
		Parameters: listTextFilesParameters,
	}
}

type listTextFilesResponse struct {
	Paths     []string `json:"paths"`
	Truncated bool     `json:"truncated"`
}

func (l *ListTextFilesTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before tool execution: %w", err)
	}

	var args struct{}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("parse list_text_files arguments: %w", err)
	}

	fileList, err := l.workspace.ListTextFiles(ctx)
	if err != nil {
		return "", fmt.Errorf("execute ListTextFiles(ctx): %w", err)
	}

	response := listTextFilesResponse{
		Paths:     fileList.Paths,
		Truncated: fileList.Truncated,
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("marshal ListTextFiles(ctx).got: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before result returned: %w", err)
	}

	return string(encoded), nil
}

var _ tool.Tool = (*ListTextFilesTool)(nil)
