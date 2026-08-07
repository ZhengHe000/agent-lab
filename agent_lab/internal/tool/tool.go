package tool

import (
	"context"
	"encoding/json"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

type Tool interface { // 表示一个可以暴露给模型并由程序执行的原子操作
	Definition() model.ToolDefinition
	Execute(ctx context.Context, arguments json.RawMessage) (string, error)
}
