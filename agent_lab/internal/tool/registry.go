package tool

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

// Registry 保存工具实现，并提供工具定义列表和名称查找能力。
type Registry struct {
	tools       map[string]Tool
	definitions []model.ToolDefinition
}

// NewRegistry 校验并注册工具。
func NewRegistry(tools ...Tool) (*Registry, error) {
	registry := &Registry{
		tools:       make(map[string]Tool, len(tools)),
		definitions: make([]model.ToolDefinition, 0, len(tools)),
	}

	for _, candidate := range tools {
		if candidate == nil {
			return nil, fmt.Errorf("工具不能为空")
		}

		definition := candidate.Definition()

		name := strings.TrimSpace(definition.Name)

		if name == "" {
			return nil, fmt.Errorf("工具名不能为空")
		}

		if _, exists := registry.tools[name]; exists {
			return nil, fmt.Errorf("工具名 %q 已存在", name)
		}

		description := strings.TrimSpace(definition.Description)
		if description == "" {
			return nil, fmt.Errorf("工具 %q 的描述不能为空", name)
		}

		if !json.Valid(definition.Parameters) {
			return nil, fmt.Errorf("工具 %q 的参数 Schema 不是合法JSON", name)
		}

		definition.Name = name
		definition.Description = description
		registry.tools[name] = candidate
		registry.definitions = append(registry.definitions, definition)
	}

	return registry, nil
}

// Get 按工具名称查找已注册工具。
func (r *Registry) Get(name string) (Tool, bool) {
	if r == nil {
		return nil, false
	}

	registeredTool, exists := r.tools[strings.TrimSpace(name)]
	return registeredTool, exists
}

// Definitions 返回工具定义的副本。
func (r *Registry) Definitions() []model.ToolDefinition {
	if r == nil {
		return nil
	}

	definitions := make([]model.ToolDefinition, len(r.definitions))
	copy(definitions, r.definitions)

	return definitions
}
