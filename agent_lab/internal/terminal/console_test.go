package terminal

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name         string
		retry        bool
		testInput    string
		wantDecision bool
	}{
		{
			name:         "允许授权",
			retry:        false,
			testInput:    "y\n",
			wantDecision: true,
		},
		{
			name:         "拒绝授权",
			retry:        false,
			testInput:    "n\n",
			wantDecision: false,
		},
		{
			name:         "错误输入后允许授权",
			retry:        true,
			testInput:    "invalid\ny\n",
			wantDecision: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.testInput)
			var output bytes.Buffer

			console, err := NewConsole(reader, &output)
			if err != nil {
				t.Fatalf("创建终端交互对象失败: %v", err)
			}

			request := tool.ConfirmationRequest{
				Action:  "test_write_file",
				Summary: "写入测试文本",
				Details: "package...main(){...}",
			}

			decision, err := console.Confirm(context.Background(), request)
			if err != nil {
				t.Fatalf("模拟交互行为失败: %v", err)
			}

			if decision != tt.wantDecision {
				t.Fatalf("模拟授权实际结果 %t 与期望结果 %t 不同", decision, tt.wantDecision)
			}

			got := output.String()

			if !strings.Contains(got, request.Action) {
				t.Fatalf("授权信息中缺少 Action 或该信息异常, 期望包含: %s", request.Action)
			}

			if !strings.Contains(got, request.Summary) {
				t.Fatalf("授权信息中缺少 Summary 或该信息异常, 期望包含: %s", request.Summary)
			}

			if !strings.Contains(got, request.Details) {
				t.Fatalf("授权信息中缺少 Details 或该信息异常, 期望包含: %s", request.Details)
			}

			if tt.retry {
				if !strings.Contains(got, "输入无效, 请输入 y 或 n: ") {
					t.Fatalf("期望存在重试提醒, 但未找到")
				}

				return 
			}

			if strings.Contains(got, "输入无效, 请输入 y 或 n: ") {
				t.Fatalf("授权过程中错误的出现了 重试提示")
			}
		})
	}
}
