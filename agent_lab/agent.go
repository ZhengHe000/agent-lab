package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Agent struct {
	client   *http.Client
	config   *Config
	messages []Message
}

type Message struct { // 消息结构体
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct { // 先只拿响应文本信息
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewAgent(client *http.Client, config *Config) *Agent {
	return &Agent{
		client:   client,
		config:   config,
		messages: []Message{{Role: "system", Content: config.Prompt}},
	}
}

func NewClient(httpClientTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: httpClientTimeout,
	}
}

func (a *Agent) callLLM(input string) (AIResponse, error) { //发送请求并处理响应
	ctx, cancel := context.WithTimeout(context.Background(), a.config.RequestTimeout) // 创建当前请求超时
	defer cancel()

	input = strings.TrimSpace(input) // 对输入简单处理

	a.messages = append(a.messages, Message{Role: "user", Content: input}) // 将输入加入上下文

	reqBody := map[string]any{ // 组装请求体
		"messages":    a.messages,
		"model":       a.config.Model,
		"temperature": 1,
	}

	data, err := json.Marshal(reqBody) //将请求体编码
	if err != nil {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		log.Println("callLLM内请求体编码失败, 返回原始err, 中断当前函数")
		return AIResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, a.config.Method, a.config.LLMAPIURL, bytes.NewReader(data)) // 创建完整请求
	if err != nil {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		log.Println("callLLM内请求创建失败, 返回原始err, 中断当前函数")
		return AIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")         // 告诉对方我发送的是json
	req.Header.Set("Accept", "application/json")               // 希望对方响应json
	req.Header.Set("Authorization", "Bearer "+a.config.APIKey) // 拼接密钥

	res, err := a.client.Do(req) // 发送请求
	if err != nil {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		log.Println("callLLM内请求发送失败, 返回原始err, 中断当前函数")
		return AIResponse{}, err
	}
	defer res.Body.Close()

	if res.Body == nil {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		return AIResponse{}, fmt.Errorf("call_llm内响应体为空")
	}

	respBody, err := io.ReadAll(res.Body) // 读取响应请求体
	if err != nil {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		log.Println("callLLM内响应读取失败, 返回原始err, 中断当前函数")
		return AIResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		log.Println("callLLM内响应状态非正确, 返回err, 中断当前函数")
		return AIResponse{}, fmt.Errorf("响应状态异常, 得到状态: %d, 得到响应: %s", res.StatusCode, string(respBody))
	}

	var assistantSaid AIResponse
	if err := json.Unmarshal(respBody, &assistantSaid); err != nil {
		a.messages = a.messages[:len(a.messages)-1] // 错误时删除本次输入
		log.Println("callLLM内请求体解析失败, 返回err, 中断当前函数")
		return AIResponse{}, err
	}

	if len(assistantSaid.Choices) == 0 {
		a.messages = a.messages[:len(a.messages)-1]
		return AIResponse{}, fmt.Errorf("API 响应成功, 但Choices为空")
	}

	if content := strings.TrimSpace(assistantSaid.Choices[0].Message.Content); content == "" {
		a.messages = a.messages[:len(a.messages)-1]
		return AIResponse{}, fmt.Errorf("API 响应成功，但回复内容为空")
	}

	a.messages = append(a.messages, Message{
		Role:    "assistant",
		Content: assistantSaid.Choices[0].Message.Content,
	}) // 将ai回复加入上下文
	return assistantSaid, nil
}
