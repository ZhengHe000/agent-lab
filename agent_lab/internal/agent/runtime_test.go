package agent

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

const testModelSystemPrompt = "test-SystemPrompt" // 没有特殊[提示词]的测试需求时 使用该变量即可
const testUserInput = "Hi"                        // 没有特殊[用户输入]的测试需求时 使用该变量即可
const testAssistantContent = "你好"                 // 没有特殊[模型响应的Content字段]的测试需求时 使用该变量即可

type fakeModel struct {
	LastRequest model.Request // 测试例中不用填这个
	Response    model.Response
	err         error
}

func (f *fakeModel) Complete(ctx context.Context, request model.Request) (model.Response, error) {
	f.LastRequest = request // 将runturn的candidate这个第一轮的提示词+输入的Messages赋值给fakeModel的字段,后续测试直接比对,用来确认模型收到的上下文是正确的
	return f.Response, f.err
}

var _ model.Model = (*fakeModel)(nil) // 编译器确认fakeModel实现了model.Model接口

func TestRunTurn(t *testing.T) {
	tests := []struct {
		name             string
		testFakeModel    *fakeModel
		testModelName    string
		testSystemPrompt string
		testInput        string
		wantErr          error
		wantMessages     []model.Message
		wantStr          string
	}{
		{
			name: "成功调用",
			testFakeModel: &fakeModel{
				Response: model.Response{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: testAssistantContent,
					},
					FinishReason: "stop",
				},
				err: nil,
			},
			testModelName:    "test-model",
			testSystemPrompt: testModelSystemPrompt,
			testInput:        testUserInput,
			wantErr:          nil,
			wantMessages: []model.Message{
				model.Message{
					Role:    model.RoleSystem,
					Content: testModelSystemPrompt,
				},
				model.Message{
					Role:    model.RoleUser,
					Content: testUserInput,
				},
				model.Message{
					Role:    model.RoleAssistant,
					Content: testAssistantContent,
				},
			},
			wantStr: testAssistantContent,
		},
		{
			name: "失败调用",
			testFakeModel: &fakeModel{
				Response: model.Response{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: testAssistantContent,
					},
					FinishReason: "stop",
				},
				err: fmt.Errorf("[测试-Err]: 调用模型时发生了xxx导致失败"),
			},
			testModelName:    "test-model",
			testSystemPrompt: testModelSystemPrompt,
			testInput:        testUserInput,
			wantErr:          ErrModelInvocationFailed,
			wantMessages: []model.Message{
				model.Message{
					Role:    model.RoleSystem,
					Content: testModelSystemPrompt,
				},
			},
			wantStr: "",
		},
		{
			name: "响应角色不是 assistant",
			testFakeModel: &fakeModel{
				Response: model.Response{
					Message: model.Message{
						Role:    model.RoleTool,
						Content: testAssistantContent,
					},
					FinishReason: "stop",
				},
				err: nil,
			},
			testModelName:    "test-model",
			testSystemPrompt: testModelSystemPrompt,
			testInput:        testUserInput,
			wantErr:          ErrResponseRoleError,
			wantMessages: []model.Message{
				model.Message{
					Role:    model.RoleSystem,
					Content: testModelSystemPrompt,
				},
			},
			wantStr: "",
		},
		{
			name: "响应包含暂不支持的 ToolCalls",
			testFakeModel: &fakeModel{
				Response: model.Response{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: testAssistantContent,
						ToolCalls: []model.ToolCall{
							model.ToolCall{
								ID:   "错误测试使用",
								Name: "错误测试使用",
							},
						},
					},
					FinishReason: "stop",
				},
				err: nil,
			},
			testModelName:    "test-model",
			testSystemPrompt: testModelSystemPrompt,
			testInput:        testUserInput,
			wantErr:          ErrToolCallsUnsupported,
			wantMessages: []model.Message{
				model.Message{
					Role:    model.RoleSystem,
					Content: testModelSystemPrompt,
				},
			},
			wantStr: "",
		},
		{
			name: "响应 Content 为空",
			testFakeModel: &fakeModel{
				Response: model.Response{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "",
					},
					FinishReason: "stop",
				},
				err: nil,
			},
			testModelName:    "test-model",
			testSystemPrompt: testModelSystemPrompt,
			testInput:        testUserInput,
			wantErr:          ErrEmptyContent,
			wantMessages: []model.Message{
				model.Message{
					Role:    model.RoleSystem,
					Content: testModelSystemPrompt,
				},
			},
			wantStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime, err := NewRuntime(tt.testFakeModel, tt.testModelName, tt.testSystemPrompt) // 拿到Runtime结构体
			if err != nil {
				t.Fatalf("调用[NewRuntime]获取[Runtime]结构体失败, 错误: %v", err)
			}

			if len(runtime.messages) != 1 {
				t.Fatalf("调用[NewRuntime]获取的[Runtime]结构体的[messages字段]错误, want Len(): %d, got Len(): %d", 1, len(runtime.messages))
			}

			wantSystemMessage := model.Message{ // 期望Runtime创建时messages只含有一份系统提示词的model.Message
				Role:    model.RoleSystem,
				Content: tt.testSystemPrompt,
			}

			if !reflect.DeepEqual(runtime.messages[0], wantSystemMessage) {
				t.Fatalf("调用[NewRuntime]获取的[Runtime]结构体的[messages字段]错误, want: %v, got: %v", wantSystemMessage, runtime.messages[0])
			}

			gotStr, err := runtime.RunTurn(context.Background(), tt.testInput)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("实际错误信息链中不包含期望的哨兵错误, wantErr: %v, gotErr: %v", tt.wantErr, err)
				}

				if gotStr != "" {
					t.Fatalf("出现错误时应该返回空字符串,但出现了值: %s", gotStr)
				}
			}

			if tt.wantErr == nil && err != nil {
				t.Fatalf("期望正确的测试出现错误, err: %v", err)
			}
			if gotStr != tt.wantStr {
				t.Fatalf("期望调用[RunTurn]得到的string与实际不符, want: %s, got: %s", tt.wantStr, gotStr)
			}

			candidate := make([]model.Message, 0, len(runtime.messages)) // 准备拼接系统提示词和输入 组成固定的第一轮请求信息
			candidate = append(candidate, runtime.messages[0])
			candidate = append(candidate, model.Message{
				Role:    model.RoleUser,
				Content: tt.testInput,
			})

			wantModelRequest := model.Request{
				Model:    runtime.modelName,
				Messages: candidate,
			}

			if !reflect.DeepEqual((*tt.testFakeModel).LastRequest, wantModelRequest) { // 将期望请求和假模型内收到的真实请求对比
				t.Fatalf("期望调用[Complete]时传入的model.Request错误, want: %v, got: %v", wantModelRequest, (*tt.testFakeModel).LastRequest)
			}

			if !reflect.DeepEqual(runtime.messages, tt.wantMessages) {
				t.Fatalf("调用[Complete]一轮结束后 此时的runtime.messages与期望得到的结构不符, want: %v, got: %v", tt.wantMessages, runtime.messages)
			}
		})
	}
}
