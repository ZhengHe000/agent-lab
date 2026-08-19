package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/model"
	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/tool"
)

// ReadTextFileTool 读取受控工作区中的指定文本文件。
type ReadTextFileTool struct {
	workspace *Workspace
}

// NewReadTextFileTool 创建使用正式工作区的文本读取工具。
func NewReadTextFileTool(workspace *Workspace) (*ReadTextFileTool, error) {
	if workspace == nil || workspace.root == nil {
		return nil, fmt.Errorf("created Read_Text_File_Tool failed: workspace is nil")
	}

	return &ReadTextFileTool{
		workspace: workspace,
	}, nil
}

type readTextFileArguments struct {
	Path string `json:"path"`
}

var readTextFileParameters = json.RawMessage(`{
            "type": "object",
            "properties": {
            "path": {
                "type": "string",
                "description": "Workspace-relative text file path using '/' separators, for example internal/config/config.go"
            }
        },
        "required": ["path"],
        "additionalProperties": false
    }`)

func (r *ReadTextFileTool) Definition() model.ToolDefinition {
	return model.ToolDefinition{
		Name: "read_text_file",
		Description: "Read an allowed text file from the controlled workspace." +
			"using a workspace-relative path.",
		Parameters: readTextFileParameters,
	}
}

func (r *ReadTextFileTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before tool execution: %w", err)
	}

	var args readTextFileArguments
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("parsed read_text_file arguments failed: %w", err)
	}

	content, err := r.workspace.ReadTextFile(args.Path)
	if err != nil {
		return "", fmt.Errorf("executed read_text_file failed: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before result returned: %w", err)
	}

	return content, nil
}

var _ tool.Tool = (*ReadTextFileTool)(nil)
