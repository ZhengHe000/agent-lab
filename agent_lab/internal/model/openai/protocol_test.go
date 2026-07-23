package openai

import (
	"encoding/json"
	"testing"
)

func TestChatCompletionResponse(t *testing.T) {
	testResponse := chatCompletionResponse{
		Choices: []chatChoice{
			{
				Message: chatMessage{
					Role:    "assistant",
					Content: nil,
					ToolCalls: []chatToolCall{
						{
							ID:   "tool_abc123",
							Type: "function",
							Function: chatFunctionCall{
								Name:      "test_tool",
								Arguments: `"filename":"test_filename"`,
							},
						},
					},
				},
				FinishReason: "调用test_tool",
			},
		},
	}

	data, err := json.Marshal(testResponse)
	if err != nil {
		t.Fatalf("测试响应序列化失败: %v", err)
	}

	var r chatCompletionResponse

	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("测试json反序列化失败: %v", err)
	}

	if len(r.Choices) == 0 {
		t.Fatal("解析后得到内容为空")
	}

	if r.Choices[0].Message.Content != nil {
		t.Fatalf("解析后得到Content与期望的 nil 不同: %q", *r.Choices[0].Message.Content)
	}

	if r.Choices[0].Message.ToolCalls[0].ID != "tool_abc123" {
		t.Fatalf("解析后得到ToolCalls[0].ID与期望的 tool_abc123 不同: %q", r.Choices[0].Message.ToolCalls[0].ID)
	}
	if r.Choices[0].Message.ToolCalls[0].Function.Name != "test_tool" {
		t.Fatalf("解析后得到Function.Name与期望的 test_tool 不同: %q", r.Choices[0].Message.ToolCalls[0].Function.Name)
	}
	if r.Choices[0].Message.ToolCalls[0].Function.Arguments != `"filename":"test_filename"` {
		t.Fatalf("解析后得到Function.Arguments与期望的 filename:test_filename 不同: %q", r.Choices[0].Message.ToolCalls[0].Function.Arguments)
	}

	if r.Choices[0].FinishReason != "调用test_tool" {
		t.Fatalf("解析后得到FinishReason与期望的 调用test_tool 不同: %q", r.Choices[0].FinishReason)
	}
}
