package tool

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

type Registry struct {
	tools       map[string]Tool
	definitions []model.ToolDefinition
}

func NewRegistry(tools ...Tool) (*Registry, error) {
	registry := &Registry{ // 创建容器并为字段追加分配空间
		tools:       make(map[string]Tool, len(tools)),
		definitions: make([]model.ToolDefinition, 0, len(tools)),
	}

	for _, candidate := range tools { // 遍历工具切片, 得到 *xxx_tool
		if candidate == nil {
			return nil, fmt.Errorf("工具不能为空")
		}

		definition := candidate.Definition() // 调用工具的Definition方法拿到 model.ToolDefinition

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

func (r *Registry) Get(name string) (Tool, bool) {
	if r == nil {
		return nil, false
	}

	registeredTool, exists := r.tools[strings.TrimSpace(name)]
	return registeredTool, exists
}

func (r *Registry) Definitions() []model.ToolDefinition {
	if r == nil {
		return nil
	}

	definitions := make([]model.ToolDefinition, len(r.definitions))
	copy(definitions, r.definitions)

	return definitions
}
