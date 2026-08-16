package openai

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/model"
)

func TestToModelResponse(t *testing.T) {
	tests := []struct {
		name         string
		wantErr      bool
		testResponse chatCompletionResponse
		check        func(t *testing.T, resp model.Response)
	}{
		{
			name:    "普通回答",
			wantErr: false,
			testResponse: chatCompletionResponse{
				Choices: []chatChoice{
					chatChoice{
						Message: chatMessage{
							Role:      "assistant",
							Content:   stringPtr("你好"),
							ToolCalls: nil,
						},
						FinishReason: "stop",
					},
				},
			},
			check: func(t *testing.T, resp model.Response) {
				t.Helper()

				if resp.FinishReason != "stop" {
					t.Fatalf("want: %s, got: %s", "stop", resp.FinishReason)
				}

				msg := resp.Message
				if msg.Role != "assistant" {
					t.Fatalf("want: %s, got: %s", "assistant", msg.Role)
				}
				if msg.Content != "你好" {
					t.Fatalf("want: %s, got: %s", "你好", msg.Content)
				}
			},
		},
		{
			name:    "工具调用",
			wantErr: false,
			testResponse: chatCompletionResponse{
				Choices: []chatChoice{
					chatChoice{
						Message: chatMessage{
							Role:    "assistant",
							Content: nil,
							ToolCalls: []chatToolCall{
								chatToolCall{
									ID:   "tool_zero",
									Type: "function",
									Function: chatFunctionCall{
										Name:      "read_test_file",
										Arguments: `{"filename":"note.txt"}`,
									},
								},
							},
						},
						FinishReason: "tool_calls",
					},
				},
			},
			check: func(t *testing.T, resp model.Response) {
				t.Helper()

				if resp.FinishReason != "tool_calls" {
					t.Fatalf("want: %s, got: %s", "tool_calls", resp.FinishReason)
				}

				msg := resp.Message
				if msg.Role != "assistant" {
					t.Fatalf("want: %s, got: %s", "assistant", msg.Role)
				}

				if len(msg.ToolCalls) != 1 {
					t.Fatalf("want: %d, got: %d", 1, len(msg.ToolCalls))
				}

				tcl := msg.ToolCalls[0]
				if tcl.ID != "tool_zero" {
					t.Fatalf("want: %s, got: %s", "tool_zero", tcl.ID)
				}
				if tcl.Name != "read_test_file" {
					t.Fatalf("want: %s, got: %s", "read_test_file", tcl.Name)
				}
				wantArguments := json.RawMessage(`{"filename":"note.txt"}`)

				if !bytes.Equal(tcl.Arguments, wantArguments) {
					t.Fatalf("Arguments: %s, want: %s", tcl.Arguments, wantArguments)
				}
			},
		},
		{
			name:         "空 Choices",
			wantErr:      true,
			testResponse: chatCompletionResponse{},
		},
		{
			name:    "错误的Role身份",
			wantErr: true,
			testResponse: chatCompletionResponse{
				Choices: []chatChoice{
					chatChoice{
						Message: chatMessage{
							Role: "NORole",
						},
					},
				},
			},
		},
		{
			name:    "错误的工具调用类型",
			wantErr: true,
			testResponse: chatCompletionResponse{
				Choices: []chatChoice{
					chatChoice{
						Message: chatMessage{
							Role: "assistant",
							ToolCalls: []chatToolCall{
								chatToolCall{
									Type: "NOType",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := toModelResponse(tt.testResponse) //调用被测函数获取转换结果

			if (err != nil) != tt.wantErr { // 判断Err是否符合预期
				t.Fatalf("want: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr { // 预期存在错误时, 判断得到的resp是否为model.Response{}
				empty := model.Response{}
				if !reflect.DeepEqual(resp, empty) {
					t.Errorf("错误时返回值必须为默认零值, 实际得到: %+v", resp)
				}
			}

			if !tt.wantErr && tt.check != nil { // 当期望不出现错误并且测试函数不为nil时, 调用测试函数对resp进行测试
				tt.check(t, resp)
			}
		})
	}
}

func stringPtr(s string) *string { // stringPtr 为 TestToModelResponse的辅助函数
	return &s
}

func TestToCompletionRequest(t *testing.T) {
	tests := []struct {
		name             string
		wantErr          bool
		testModelRequest model.Request
		check            func(t *testing.T, chatReq chatCompletionRequest)
	}{
		{
			name:    "普通用户消息",
			wantErr: false,
			testModelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					model.Message{
						Role:    model.RoleUser,
						Content: "你好",
					},
				},
			},
			check: func(t *testing.T, chatReq chatCompletionRequest) {
				if chatReq.Model != "test-model" {
					t.Fatalf("结果中[Model]字段与预期不符, want: %s, got: %s", "test-model", chatReq.Model)
				}

				if len(chatReq.Messages) == 0 {
					t.Fatalf("结果中[Message]字段异常, want Len: %d, got Len: %d", 1, len(chatReq.Messages))
				}
				message := chatReq.Messages[0]

				if message.Role != "user" {
					t.Fatalf("结果中[Role]字段异常, want: %s, got Len: %v", "user", message.Role)
				}
				if message.Content == nil {
					t.Fatalf("结果中[Content]字段异常 当前值为nil, 期望值为非nil")
				}
				if *message.Content != "你好" {
					t.Fatalf("结果中[Content]字段异常, want: %s, got Len: %v", "你好", *message.Content)
				}
			},
		},
		{
			name:    "Assistant 工具调用",
			wantErr: false,
			testModelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					model.Message{
						Role:    model.RoleAssistant,
						Content: "",
						ToolCalls: []model.ToolCall{
							model.ToolCall{
								ID:        "tool_read",
								Name:      "read_test_file",
								Arguments: json.RawMessage(`{"filename":"note.txt"}`),
							},
						},
					},
				},
			},
			check: func(t *testing.T, chatReq chatCompletionRequest) {
				if chatReq.Model != "test-model" {
					t.Fatalf("结果中[Model]字段与预期不符, want: %s, got: %s", "test-model", chatReq.Model)
				}

				if len(chatReq.Messages) == 0 {
					t.Fatalf("结果中[Message]字段异常, want Len: %d, got Len: %d", 1, len(chatReq.Messages))
				}

				message := chatReq.Messages[0]

				if message.Role != "assistant" {
					t.Fatalf("结果中[Role]字段异常, want: %s, got: %s", "assistant", message.Role)
				}

				if message.Content != nil {
					t.Fatalf("结果中[Content]字段异常, 期望该字段为nil, 但得到: %v", message.Content)
				}

				if len(message.ToolCalls) == 0 {
					t.Fatalf("结果中[ToolCalls]字段异常, want Len: %d, got Len: %d", 1, len(message.ToolCalls))
				}

				chatTCl := message.ToolCalls[0]

				if chatTCl.ID != "tool_read" {
					t.Fatalf("结果中[ID]字段异常, want %s, got: %s", "tool_read", chatTCl.ID)
				}

				if chatTCl.Type != "function" {
					t.Fatalf("结果中[Type]字段异常, want %s, got: %s", "function", chatTCl.Type)
				}

				if chatTCl.Function.Name != "read_test_file" {
					t.Fatalf("结果中[Name]字段异常, want %s, got: %s", "read_test_file", chatTCl.Function.Name)
				}

				if chatTCl.Function.Arguments != `{"filename":"note.txt"}` {
					t.Fatalf("结果中[Arguments]字段异常, want %s, got: %s", `{"filename":"note.txt"}`, chatTCl.Function.Arguments)
				}
			},
		},
		{
			name:    "工具结果消息",
			wantErr: false,
			testModelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					model.Message{
						Role:       model.RoleTool,
						Content:    "got-testnote...",
						ToolCallID: "tool-read",
					},
				},
			},
			check: func(t *testing.T, chatReq chatCompletionRequest) {
				if chatReq.Model != "test-model" {
					t.Fatalf("结果中[Model]字段与预期不符, want: %s, got: %s", "test-model", chatReq.Model)
				}

				if len(chatReq.Messages) == 0 {
					t.Fatalf("结果中[Message]字段异常, want Len: %d, got Len: %d", 1, len(chatReq.Messages))
				}

				message := chatReq.Messages[0]

				if message.Role != "tool" {
					t.Fatalf("结果中[Role]字段异常, want: %s, got: %s", "tool", message.Role)
				}

				if message.Content == nil {
					t.Fatalf("结果中[Content]字段异常, 期望字段为非nil, 但得到nil")
				}

				if *message.Content != "got-testnote..." {
					t.Fatalf("结果中[Content]字段异常, want: %s, got: %s", "got-testnote...", *message.Content)
				}

				if message.ToolCallID != "tool-read" {
					t.Fatalf("结果中[ToolCallID]字段异常, want: %s, got: %s", "tool-read", message.ToolCallID)
				}
			},
		},
		{
			name:    "非法内部角色",
			wantErr: true,
			testModelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					model.Message{
						Role:    "我是非法角色-暗黑雷龙暴风雨-我要来攻击程序了!",
						Content: "恶意内容",
					},
				},
			},
		},
		{
			name:    "包含工具定义",
			wantErr: false,
			testModelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					model.Message{
						Role:    model.RoleUser,
						Content: "查看文件",
					},
				},
				Tools: []model.ToolDefinition{
					model.ToolDefinition{
						Name:        "test_read_text_file",
						Description: "测试-用户需要查看文件内容时使用",
						Parameters: json.RawMessage(`{
						"type":"object",
					"properties":{
						"filename":{"type":"string"}
					},
					"required":["filename"]
					}`),
					},
				},
			},
			check: func(t *testing.T, chatReq chatCompletionRequest) {
				t.Helper()

				if len(chatReq.Tools) != 1 {
					t.Fatalf("Tools数量错误, want: 1, got: %d", len(chatReq.Tools))
				}

				got := chatReq.Tools[0]
				if got.Type != "function" {
					t.Fatalf("工具类型错误, want: function, got: %q", got.Type)
				}

				if got.Function.Name != "test_read_text_file" {
					t.Fatalf("工具名称错误, got: %q", got.Function.Name)
				}

				if got.Function.Description != "测试-用户需要查看文件内容时使用" {
					t.Fatalf("工具描述错误, got: %q", got.Function.Description)
				}

				if !json.Valid(got.Function.Parameters) {
					t.Fatal("转换后的工具参数不是合法JSON")
				}
			},
		},
		{
			name:    "工具参数Schema非法",
			wantErr: true,
			testModelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					model.Message{
						Role:    model.RoleUser,
						Content: "test",
					},
				},
				Tools: []model.ToolDefinition{
					model.ToolDefinition{
						Name:        "test_read_text_file",
						Description: "测试-用户需要查看文件内容时使用",
						Parameters:  json.RawMessage(`{"type":`),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := toChatCompletionRequest(tt.testModelRequest)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望err不为空, 但实际没有发生错误")
				}

				empty := chatCompletionRequest{}
				if !reflect.DeepEqual(resp, empty) {
					t.Fatalf("错误时返回值必须为默认零值, 实际得到: %+v", resp)
				}

				return
			}

			tt.check(t, resp)
		})
	}
}
