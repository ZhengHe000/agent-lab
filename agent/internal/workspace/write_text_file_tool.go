package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/model"
	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/tool"
)

// WriteTextFileTool 经用户确认后，在受控工作区中创建或覆盖文本文件。
type WriteTextFileTool struct {
	workspace *Workspace
	confirmer tool.Confirmer
}

// NewWriteTextFileTool 创建使用正式工作区的文本写入工具。
func NewWriteTextFileTool(workspace *Workspace, confirmer tool.Confirmer) (*WriteTextFileTool, error) {
	if workspace == nil || workspace.root == nil {
		return nil, fmt.Errorf("create write_text_file tool: workspace is nil")
	}

	if confirmer == nil {
		return nil, fmt.Errorf(
			"create write_text_file tool: confirmer is nil",
		)
	}

	return &WriteTextFileTool{
		workspace: workspace,
		confirmer: confirmer,
	}, nil
}

type writeTextFileArguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

var writeTextFileParameters = json.RawMessage(`{
"type": "object",
"properties": {
"path":{
"type": "string",
"description": "Workspace-relative text file path using '/' separators, for example internal/config/config.go"
},
"content": {
"type": "string",
"description": "Complete text content to write to the file"
	}
},
"required":["path", "content"],
"additionalProperties": false
}`)

func (t *WriteTextFileTool) Definition() model.ToolDefinition {
	return model.ToolDefinition{
		Name: "write_text_file",
		Description: "Create or replace a supported regular text file in the controlled workspace " +
			"using a workspace-relative path. Use this tool only when the user has requested a file " +
			"change. The proposed path and complete content are shown for confirmation before any " +
			"file modification is performed.",
		Parameters: writeTextFileParameters,
	}
}

func (t *WriteTextFileTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before tool execution: %w", err)
	}

	args, err := tool.DecodeObjectArguments[writeTextFileArguments](arguments)
	if err != nil {
		return "", fmt.Errorf(
			"parse write_text_file arguments: %w",
			err,
		)
	}

	toolPath, err := t.workspace.validateTextFileWrite(args.Path, args.Content)
	if err != nil {
		return "", fmt.Errorf(
			"validate write_text_file arguments: %w",
			err,
		)
	}

	confirmed, err := t.confirmer.Confirm(
		ctx,
		tool.ConfirmationRequest{
			Action:  "write_text_file",
			Summary: fmt.Sprintf("创建或覆盖文件 %s", toolPath),
			Details: fmt.Sprintf(
				"内容长度: %d 字符, %d 字节\n\n内容:\n%s",
				utf8.RuneCountInString(args.Content),
				len(args.Content),
				args.Content,
			),
		},
	)

	if err != nil {
		return "", fmt.Errorf("Authorization write text file operation failed: %w", err)
	}

	if !confirmed {
		return fmt.Sprintf(
			"The write file %s operation was rejected, No modifications were made",
			toolPath,
		), nil
	}

	if err = t.workspace.WriteTextFile(toolPath, args.Content); err != nil {
		return "", fmt.Errorf(
			"write text file %q: %w",
			toolPath,
			err,
		)
	}

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before result returned: %w", err)
	}

	return fmt.Sprintf("File %q written successfully", toolPath), nil
}

var _ tool.Tool = (*WriteTextFileTool)(nil)
