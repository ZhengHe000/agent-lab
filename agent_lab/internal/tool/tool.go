package tool

import (
	"context"
	"encoding/json"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

type Tool interface { // Tool 表示一个可以暴露给模型并由程序执行的原子操作
	Definition() model.ToolDefinition                                       // 返回模型选择和调用工具所需的说明
	Execute(ctx context.Context, arguments json.RawMessage) (string, error) // 使用模型返回的JSON作为参数执行工具
}
