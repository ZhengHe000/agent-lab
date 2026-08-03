package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
)

type ReadTextFileTool struct {
	dir string
}

func NewReadTextFileTool() *ReadTextFileTool { // 固定使用受控工作区
	return &ReadTextFileTool{
		dir: workspaceDir,
	}
}

func newReadTextFileToolInDir(dir string) *ReadTextFileTool { // 用于在workspace的测试中注入临时测试目录 不对外开放
	return &ReadTextFileTool{
		dir: dir,
	}
}

type readTextFileArguments struct { // 工具内部协议
	FileName string `json:"filename"`
}

var readTextFileParameters = json.RawMessage(`{
            "type": "object",
            "properties": {
            "filename": {
                "type": "string",
                "description": "要读取的文本文件名, 例如note.txt",
            }
        },
        "required": ["filename"],
        "additionalProperties": false
    }`)

func (r *ReadTextFileTool) Definition() *model.ToolDefinition {
	return &model.ToolDefinition{
		Name: "read_text_name",
		Description: "读取受控工作区指定中文本文件内容的完整内容." +
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

	content, err := readTextFileInDir(r.dir, args)
	if err != nil {
		return "", fmt.Errorf("执行 read_text_file 失败: %w", err)
	}

	return content, err
}

var _ Tool = (*ReadTextFileTool)(nil)
