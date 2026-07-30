package openai

import (
	"context"
	"encoding/json"
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
				if r.Method != http.MethodPost { // 检查请求方法
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				if r.Header.Get("Content-Type") != "application/json" {
					w.WriteHeader(http.StatusUnauthorized)
					t.Errorf("Header的Content-Type信息错误, want: %s,got: %s", "application/json", r.Header.Get("Content-Type"))
					return
				}

				if r.Header.Get("Accept") != "application/json" {
					w.WriteHeader(http.StatusUnauthorized)
					t.Errorf("Header的Accept信息错误, want: %s,got: %s", "application/json", r.Header.Get("Accept"))
					return
				}

				if r.Header.Get("Authorization") != "Bearer test-Key" { // 检查Header信息
					w.WriteHeader(http.StatusUnauthorized)
					t.Errorf("Header的Authorization信息错误, want: %s,got: %s", "Bearer test-Key", r.Header.Get("Authorization"))
					return
				}

				if r.ContentLength == 0 { // 检查请求体是否为空
					w.WriteHeader(http.StatusBadRequest)
					t.Errorf("请求体为空")
					return
				}

				var apiResponse chatCompletionRequest
				if err := json.NewDecoder(r.Body).Decode(&apiResponse); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					t.Errorf("解析请求体失败: %v", err)
					return
				}

				if apiResponse.Model != "test-model" {
					t.Errorf("请求体Model字段错误, want: %s, got %s", "test-model", apiResponse.Model)
				}

				if len(apiResponse.Messages) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					t.Errorf("请求体Messages字段为空")
					return
				}

				msg := apiResponse.Messages[0]

				if msg.Content == nil {
					t.Errorf("请求体 Content 为 nil")
					return
				}

				if msg.Role != "user" {
					t.Errorf("请求体Role字段错误, want: %s, got %s", "user", msg.Role)
				}

				if *msg.Content != "Yes" {
					t.Errorf("请求体Content字段错误, want: %s, got %s", "Yes", *msg.Content)
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
			t.Cleanup(tt.testServer.Close)

			testClient, err := NewClient(tt.testServer.Client(), tt.testServer.URL, "test-Key")
			if err != nil {
				t.Fatalf("测试客户端创建失败, 错误: %v", err)
			}

			var modelClient model.Model = testClient

			modelResponse, err := modelClient.Complete(context.Background(), tt.modelRequest)

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

			if !tt.wantErr && err != nil {
				t.Fatalf("Complete() 出现意外错误: %v", err)
			}

			if modelResponse.FinishReason != "stop" {
				t.Fatalf("want: %s, got: %s", "stop", modelResponse.FinishReason)
			}

			if modelResponse.Message.Role != model.RoleAssistant {
				t.Fatalf("want: %s, got: %s", model.RoleAssistant, modelResponse.Message.Role)
			}

			if tt.wantContent != "" {
				if modelResponse.Message.Content != tt.wantContent {
					t.Fatalf("want: %s, got: %s", tt.wantContent, modelResponse.Message.Content)
				}
			}
		})
	}
}
