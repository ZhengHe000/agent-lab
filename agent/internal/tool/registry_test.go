package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/model"
)

var testParameters json.RawMessage = json.RawMessage(`{"test":"arguments"}`)

type fakeTool struct {
	definition model.ToolDefinition
}

func (f *fakeTool) Definition() model.ToolDefinition {
	return f.definition
}

func (f *fakeTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	return "test-result", nil
}

var toolA Tool = &fakeTool{ // 正常 A
	definition: model.ToolDefinition{
		Name:        "test_read_file",
		Description: "测试工具描述A",
		Parameters:  testParameters,
	},
}

var toolB Tool = &fakeTool{ // 正常 B
	definition: model.ToolDefinition{
		Name:        "test_list_file",
		Description: "测试工具描述B",
		Parameters:  testParameters,
	},
}
var toolC Tool = &fakeTool{ // Parameters字段 存在JSON错误 C
	definition: model.ToolDefinition{
		Name:        "test_list_file",
		Description: "测试工具描述B",
		Parameters:  json.RawMessage(`{"test":"arguments"`),
	},
}

func TestRegistry(t *testing.T) {
	tests := []struct {
		name                string
		tools               []Tool
		wantErr             bool
		wantBool            bool
		setName             string
		wantTool            Tool
		wantToolName        string
		wantToolDescription string
		wantToolParameters  json.RawMessage
		wantDefinitions     []model.ToolDefinition
	}{
		{
			name: "全部正常",
			tools: []Tool{
				toolA,
				toolB,
			},
			wantErr:  false,
			wantBool: true,
			setName:  "test_read_file",
			wantTool: toolA,
			wantDefinitions: []model.ToolDefinition{
				toolA.Definition(),
				toolB.Definition(),
			},
		},
		{
			name: "Get不存在的名称返回false",
			tools: []Tool{
				toolA,
				toolB,
			},
			wantErr:  false,
			wantBool: false,
			setName:  "test_read_999",
		},
		{
			name: "拒绝重复名称",
			tools: []Tool{
				toolA,
				toolA,
			},
			wantErr: true,
		},
		{
			name: "Parameters不是合法JSON",
			tools: []Tool{
				toolA,
				toolC,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry, err := NewRegistry(tt.tools...)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望出错, 但Err实际为nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("不期望出错, 但: %v", err)
			}

			gotTool, gotBool := registry.Get(tt.setName)

			if !tt.wantBool {
				if gotBool == true {
					t.Fatalf("期望Get()返回 false, 但实际为 true")
				}
				return
			}

			if !gotBool {
				t.Fatalf("期望Get()返回 true, 但实际为 false")
			}

			definition := gotTool.Definition()
			wantdefinition := tt.wantTool.Definition()

			if definition.Name != wantdefinition.Name ||
				definition.Description != wantdefinition.Description ||
				!bytes.Equal(definition.Parameters, wantdefinition.Parameters) {
				t.Fatalf("Get() 得到的Tool 调用Definition() 后获得的信息与期望不符, got: %+v, wantName: %s, wantDescription: %s, wantParameters: %s", definition, wantdefinition.Name, wantdefinition.Description, string(wantdefinition.Parameters))
			}

			definitions := registry.Definitions()
			if !reflect.DeepEqual(definitions, tt.wantDefinitions) {
				t.Fatalf("得到的[]Tool: %v, 与期望的: %v 不符", definitions, tt.wantDefinitions)
			}
		})
	}
}
