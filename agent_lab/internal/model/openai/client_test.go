package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

func TestClientComplete(t *testing.T) {
	tests := []struct {
		name         string
		wantErr      bool
		wantContent  string
		testServer   *httptest.Server
		modelRequest model.Request
	}{
		{
			name:        "成功响应",
			wantErr:     false,
			wantContent: "你好",
			testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				if r.Header.Get("Authorization") != "Bearer test-Key" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
				"choices":[{
				"message": {"role":"assistant","content":"你好"},
				"finish_reason":"stop"}]}`))
			})),
			modelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					{Role: model.RoleUser,
						Content: "Yes",
					},
				},
			},
		},
		{
			name:        "非2xx",
			wantErr:     true,
			wantContent: "",
			testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"[测试]响应错误信息..."}`))
			})),
			modelRequest: model.Request{
				Model: "test-model",
				Messages: []model.Message{
					{Role: model.RoleUser,
						Content: "No",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.testServer.Close()

			testClient, err := NewClient(tt.testServer.Client(), tt.testServer.URL, "test-Key")
			if err != nil {
				t.Fatalf("测试客户端创建失败, 错误: %v", err)
			}

			modelResponse, err := testClient.Complete(context.Background(), tt.modelRequest)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望出错的实例没有出错, 得到的modelResponse: %v", modelResponse)
				}

				empty := model.Response{}
				if !reflect.DeepEqual(modelResponse, empty) {
					t.Fatalf("期望错误时返回空结构信息, 但实际出现了值: %v", modelResponse)
				}

				return
			}

			if tt.wantContent != "" {
				if modelResponse.Message.Content != tt.wantContent {
					t.Fatalf("want: %s, got: %s", tt.wantContent, modelResponse.Message.Content)
				}
			}
		})
	}
}
