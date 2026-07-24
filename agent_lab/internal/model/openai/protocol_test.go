package openai

import (
	"encoding/json"
	"testing"
)

func TestChatCompletionResponse(t *testing.T) {
	responseJSON := []byte(`{
	"choices": [
		{"finish_reason": "tool_calls",
		"message":{ 
			"role": "assistant",
			"content": null,
			"tool_calls": [ 
					{
						"id": "abc_tool123",
						"type": "function",
						"function": {
										"name": "read_test_file",
										"arguments":"{\"filename\":\"test\"}"
									}
					}
							]
				}
		}
				]
							}`)

	var response chatCompletionResponse

	if err := json.Unmarshal(responseJSON, &response); err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if len(response.Choices) != 1 {
		t.Fatalf("Choices数量 want: 1, got: %d", len(response.Choices))
	}
	choice := response.Choices[0]

	if len(choice.Message.ToolCalls) != 1 {
		t.Fatalf("ToolCalls数量 want: 1, got: %d", len(choice.Message.ToolCalls))
	}

	call := choice.Message.ToolCalls[0]

	if choice.Message.Content != nil {
		t.Fatalf("Content与期望值不同 want: %v, got: %v", nil, choice.Message.Content)
	}

	if choice.FinishReason != "tool_calls" {
		t.Fatalf("FinishReason与期望值不同 want: %q, got: %q", "tool_calls", choice.FinishReason)
	}
	if call.ID != "abc_tool123" {
		t.Fatalf("call.ID与期望值不同 want: %q, got: %q", "abc_tool123", call.ID)
	}

	if call.Type != "function" {
		t.Fatalf("call.Type与期望值不同 want: %q, got: %q", "function", call.Type)
	}

	if call.Function.Name != "read_test_file" {
		t.Fatalf(" call.Function.name与期望值不同 want: %q, got: %q", "read_test_file", call.Function.Name)
	}
	if call.Function.Arguments != `{"filename":"test"}` {
		t.Fatalf("call.Function.Arguments与期望值不同 want: %q, got: %q", `{"filename":"test"}`, call.Function.Arguments)
	}
}
