package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

type Client struct { // 定义客户端结构体
	httpClient *http.Client
	endpoint   string
	apiKey     string
}

func NewClient(
	httpClient *http.Client,
	endpoint string,
	apiKey string,
) (*Client, error) {
	if httpClient == nil { // 拒绝客户端配置为nil
		return nil, fmt.Errorf("httpClient 参数不能为 nil")
	}

	if strings.TrimSpace(endpoint) == "" { // 拒绝字段空值
		return nil, fmt.Errorf("endpoint值不可为空")
	}

	if strings.TrimSpace(apiKey) == "" { // 拒绝字段空值
		return nil, fmt.Errorf("apiKey值不可为空")
	}

	return &Client{ // 返回组装结构体的地址和nil
		httpClient: httpClient,
		endpoint:   endpoint,
		apiKey:     apiKey,
	}, nil
}

func (c *Client) Complete(ctx context.Context, request model.Request) (model.Response, error) {
	externalRequest, err := toChatCompletionRequest(request) // 将内部请求转换外部请求
	if err != nil {
		return model.Response{}, err
	}

	requestJSON, err := json.Marshal(externalRequest) // 将外部请求编码为JSON
	if err != nil {
		return model.Response{}, fmt.Errorf("外部请求编码失败, 请求: %v, 错误: %w", externalRequest, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(requestJSON)) // 创建https请求
	if err != nil {
		return model.Response{}, fmt.Errorf("请求创建失败, 错误: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")  // 说明我方请求使用JSON发送请求
	req.Header.Set("Accept", "application/json")        // 告知对方使用JSON返回响应
	req.Header.Set("Authorization", "Bearer "+c.apiKey) // 拼接apiKey, 获得访问能力

	resp, err := c.httpClient.Do(req) // 发送请求
	if err != nil {
		return model.Response{}, fmt.Errorf("请求发送失败, 错误: %w", err)
	}
	defer resp.Body.Close() // 链接成功, 退出前释放文件描述符

	respBytes, err := io.ReadAll(resp.Body) // 先读取请求体内容, 读完后 即使后续Header错误 该链接可以复用
	if err != nil {
		return model.Response{}, fmt.Errorf("读取响应体失败, 错误: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices { //检查响应状态
		return model.Response{}, fmt.Errorf("响应码错误, 响应体: %s, 状态码: %d", string(respBytes), resp.StatusCode)
	}

	var chatCompletionResponse chatCompletionResponse                          // 创建容器
	if err := json.Unmarshal(respBytes, &chatCompletionResponse); err != nil { // 将响应内容解析到外部响应结构体中
		return model.Response{}, fmt.Errorf("响应解析失败, 错误: %w", err)
	}

	modelResponse, err := toModelResponse(chatCompletionResponse) // 将外部响应结构体转换为内部响应结构体
	if err != nil {
		return model.Response{}, err
	}

	return modelResponse, nil // 返回最终数据和空错误
}

var _ model.Model = (*Client)(nil) // 编译器检查Client是否实现model包的Model接口
