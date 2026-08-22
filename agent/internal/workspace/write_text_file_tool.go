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
		return nil, fmt.Errorf("create read_text_file tool: workspace is nil")
	}

	if confirmer == nil {
		return nil, fmt.Errorf("写入工具确认器不能为空")
	}

	return &WriteTextFileTool{
		workspace: workspace,
		confirmer: confirmer,
	}, nil
}

type writeTextFileArguments struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

var writeTextFileParameters = json.RawMessage(`{
"type": "object",
"properties": {
"filename":{
"type": "string",
"description": "要创建或覆盖的文本文件名, 例如 note.txt"
},
"content": {
"type": "string",
"description": "要写入文件的完整文本内容"
	}
},
"required":["filename", "content"],
"additionalProperties": false
}`)

func (t *WriteTextFileTool) Definition() model.ToolDefinition {
	return model.ToolDefinition{
		Name: "write_text_file",
		Description: "经过用户明确确认后, 在受控工作区中创建或覆盖文本文件." +
			"仅当用户要求保存或创建或修改文件内容时使用",
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

	confirmed, err := t.confirmer.Confirm(
		ctx,
		tool.ConfirmationRequest{
			Action:  "write_text_file",
			Summary: fmt.Sprintf("创建或覆盖文件 %s", args.Filename),
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
			args.Filename,
		), nil
	}

	if err = t.workspace.WriteTextFile(args.Filename, args.Content); err != nil {
		return "", fmt.Errorf(
			"write text file %q: %w",
			args.Filename,
			err,
		)
	}

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before result returned: %w", err)
	}

	return fmt.Sprintf("File %q written successfully", args.Filename), nil
}

var _ tool.Tool = (*WriteTextFileTool)(nil)
