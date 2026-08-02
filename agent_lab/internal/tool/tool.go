package tool

import (
	"context"
	"encoding/json"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

type Tool interface {
	Definition() model.ToolDefinition
	Execute(ctx context.Context, arguments json.RawMessage) (string, error)
}
