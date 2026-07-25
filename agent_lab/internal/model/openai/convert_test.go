package openai

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
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
										Name:      "read_file",
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
				if tcl.Name != "read_file" {
					t.Fatalf("want: %s, got: %s", "read_file", tcl.Name)
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

func stringPtr(s string) *string {
	return &s
}
