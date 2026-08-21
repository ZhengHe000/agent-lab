package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/model"
	"github.com/ZhengHeOwo/agent-AnXuan/agent/internal/tool"
)

// WriteTextFileTool 经用户确认后，在受控工作区中创建或覆盖文本文件。
type WriteTextFileTool struct {
	dir       string
	confirmer tool.Confirmer
}

// NewWriteTextFileTool 创建使用正式工作区的文本写入工具。
func NewWriteTextFileTool(confirmer tool.Confirmer) (*WriteTextFileTool, error) {
	return newWriteTextFileTool(DefaultDir, confirmer)
}

func newWriteTextFileTool(dir string, confirmer tool.Confirmer) (*WriteTextFileTool, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, fmt.Errorf("写入工具工作区目录不能为空")
	}

	if confirmer == nil {
		return nil, fmt.Errorf("写入工具确认器不能为空")
	}

	return &WriteTextFileTool{
		dir:       dir,
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
		return "", fmt.Errorf("写入文件前上下文已经结束: %w", err)
	}

	args, err := tool.DecodeObjectArguments[writeTextFileArguments](arguments)
	if err != nil {
		return "", fmt.Errorf(
			"parse write_text_file arguments: %w",
			err,
		)
	}

	if err := validateFilename(args.Filename); err != nil {
		return "", fmt.Errorf("write_text_file 文件名无效: %w", err)
	}

	if err := validateContent(args.Content); err != nil {
		return "", fmt.Errorf("write_text_file 内容无效: %w", err)
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
		return "", fmt.Errorf("确认 write_text_file 操作失败: %w", err)
	}

	if !confirmed {
		return fmt.Sprintf(
			"用户拒绝写入文件 %s, 未执行任何修改",
			args.Filename,
		), nil
	}

	_, err = writeTextFileInDir(t.dir, args.Filename, args.Content)
	if err != nil {
		return "", fmt.Errorf("执行 write_text_file 失败: %w", err)
	}

	return fmt.Sprintf("已成功写入文件 %q", args.Filename), nil
}

var _ tool.Tool = (*WriteTextFileTool)(nil)
