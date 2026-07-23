package model

import (
	"context"
	"testing"
)

type fakeModel struct {
	response Response
	err      error
}

func (f *fakeModel) Complete(ctx context.Context, request Request) (Response, error) {
	return f.response, f.err
}

var _ Model = (*fakeModel)(nil)

func TestModelInterface(t *testing.T) {
	want := Response{
		Message: Message{
			Role:       RoleAssistant,
			Content:    "测试",
			ToolCalls:  nil,
			ToolCallID: "",
		},
		FinishReason: "stop",
	}

	var m Model = &fakeModel{response: want}

	got, err := m.Complete(context.Background(), Request{
		Model: "zero-Model",
		Messages: []Message{
			{
				Role:    RoleUser,
				Content: "你好",
			},
		},
	})

	if err != nil {
		t.Fatalf("期望类型与实际不符, gotErr: %v", err)
	}

	if got.Message.Content != want.Message.Content {
		t.Fatalf("回答: %q, want: %q", got.Message.Content, want.Message.Content)
	}
}
