package tool

import (
	"context"
	"encoding/json"

	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/model"
)

// Tool 表示可以暴露给模型并由程序执行的原子操作。
type Tool interface {
	Definition() model.ToolDefinition
	Execute(ctx context.Context, arguments json.RawMessage) (string, error)
}
