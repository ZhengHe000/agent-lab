package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
)

type fakeConfirmer struct {
	calls    int
	err      error
	request  tool.ConfirmationRequest
	decision bool
}

func (f *fakeConfirmer) Confirm(ctx context.Context, request tool.ConfirmationRequest) (bool, error) {
	f.calls++
	f.request = request
	return f.decision, f.err
}

// 下面的 "test.txt" 作为固定文件名, 后续新增测试例中请勿使用其他文件名作为 Arguments 中的"filename"值
var arguments = json.RawMessage(`{"filename":"test.txt","content":"Hi"}`)

func TestWriteTextFileTool(t *testing.T) {
	tests := []struct {
		name              string
		confirmerErr      error
		confirmerDecision bool
		wantErr           bool
		wantExecuteErr    error
		wantCallCount     int
		wantFile          bool
		testArguments     json.RawMessage
		wantContent       string
	}{
		{
			name:              "允许写入",
			confirmerErr:      nil,
			confirmerDecision: true,
			wantErr:           false,
			wantCallCount:     1,
			wantFile:          true,
			testArguments:     arguments,
			wantContent:       "Hi",
		},
		{
			name:              "拒绝写入",
			confirmerErr:      nil,
			confirmerDecision: false,
			wantErr:           false,
			wantCallCount:     1,
			wantFile:          false,
			testArguments:     arguments,
			wantContent:       "",
		},
		{
			name:           "参数非法时不请求确认",
			wantErr:        true,
			wantExecuteErr: ErrInvalidFilename,
			wantCallCount:  0,
			wantFile:       false,
			testArguments:  json.RawMessage(`{"filename":"../secret","content":"Hi"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			testPath := filepath.Join(testDir, "test.txt")

			var testConfirmer = &fakeConfirmer{
				err:      tt.confirmerErr,
				decision: tt.confirmerDecision,
			}

			writeTool, err := newWriteTextFileTool(testDir, testConfirmer)
			if err != nil {
				t.Fatalf("创建 write_text_file 工具失败: %v", err)
			}

			got, err := writeTool.Execute(context.Background(), tt.testArguments)

			if testConfirmer.calls != tt.wantCallCount {
				t.Fatalf("确认器调用次数为 %d次 与期望的 %d次 不符", testConfirmer.calls, tt.wantCallCount)
			}

			if tt.wantFile {
				_, err := os.Stat(testPath)

				if os.IsNotExist(err) {
					t.Fatalf("实际未在测试目录下创建期望的 test.txt 文件")
				}

				data, err := os.ReadFile(testPath)
				if err != nil {
					t.Fatalf("读取测试文件失败")
				}

				if string(data) != tt.wantContent {
					t.Fatalf("写入实际内容为: %s, 期待内容为: %s", string(data), tt.wantContent)
				}
			}

			if tt.name == "拒绝写入" { // 单独分出这个情况避免影响整体协调性
				if err != nil {
					t.Fatalf("期望nil, 但实际err: %v", err)
				}

				if tt.wantExecuteErr != nil {
					if !errors.Is(err, tt.wantExecuteErr) {
						t.Fatalf("[错误链] 期望包含错误: %v, 实际得到: %v", tt.wantExecuteErr, err)
					}
				}

				_, err = os.Stat(testPath)
				if !os.IsNotExist(err) {
					t.Fatalf("程序错误的创建 test.txt 文件")
				}
			}

			if tt.wantErr {

				if err == nil {
					t.Fatalf("期望错误, 但实际err为空")
				}

				if tt.wantExecuteErr != nil {
					if !errors.Is(err, tt.wantExecuteErr) {
						t.Fatalf("[错误链] 期望包含错误: %v, 实际得到: %v", tt.wantExecuteErr, err)
					}
				}

				_, err = os.Stat(testPath)
				if !os.IsNotExist(err) {
					t.Fatalf("程序错误的创建 test.txt 文件")
				}

				return
			}

			if testConfirmer.request.Action != "write_text_file" ||
				!strings.Contains(testConfirmer.request.Summary, "创建或覆盖文件") ||
				!strings.Contains(testConfirmer.request.Details, "Hi") {
				t.Fatalf("确认请求没有准确展示写入操作")
			}
			
			if !strings.Contains(got, "test.txt") {
				t.Fatalf("期望 write_text_file 的调用结果中包换test.txt, 但实际: %s", got)
			}

		})
	}
}
