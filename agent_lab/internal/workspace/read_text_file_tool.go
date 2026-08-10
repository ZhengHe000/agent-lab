package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
)

// ReadTextFileTool 读取受控工作区中的指定文本文件。
type ReadTextFileTool struct {
	dir string
}

// NewReadTextFileTool 创建使用正式工作区的文本读取工具。
func NewReadTextFileTool() *ReadTextFileTool {
	return &ReadTextFileTool{
		dir: workspaceDir,
	}
}

func newReadTextFileToolInDir(dir string) *ReadTextFileTool {
	return &ReadTextFileTool{
		dir: dir,
	}
}

type readTextFileArguments struct {
	Filename string `json:"filename"`
}

var readTextFileParameters = json.RawMessage(`{
            "type": "object",
            "properties": {
            "filename": {
                "type": "string",
                "description": "要读取的文本文件名, 例如note.txt"
            }
        },
        "required": ["filename"],
        "additionalProperties": false
    }`)

func (r *ReadTextFileTool) Definition() model.ToolDefinition {
	return model.ToolDefinition{
		Name: "read_text_file",
		Description: "读取受控工作区中指定文本文件的完整内容." +
			"当用户需要查看某个工作区文本文件时使用.",
		Parameters: readTextFileParameters,
	}
}

func (r *ReadTextFileTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("读取文件前上下文已结束: %w", err)
	}

	var args readTextFileArguments
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("解析 read_text_file 参数失败: %w", err)
	}

	content, err := readTextFileInDir(r.dir, args.Filename)
	if err != nil {
		return "", fmt.Errorf("执行 read_text_file 失败: %w", err)
	}

	return content, nil
}

var _ tool.Tool = (*ReadTextFileTool)(nil)
