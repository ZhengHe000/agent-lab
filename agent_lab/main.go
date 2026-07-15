package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Agent struct {
	client *http.Client
	config *Config
}

func NewAgent(client *http.Client, config *Config) *Agent {
	return &Agent{
		client: client,
		config: config,
	}
}

func loadEnv(fileName string) error { // 配置环境变量
	file, err := os.Open(fileName) //打开env环境配置文件
	if err != nil {
		return fmt.Errorf("打开 %v 文件失败, 错误: %w", fileName, err)
	}

	defer file.Close() // 结束前关闭文件

	scanner := bufio.NewScanner(file) // 创建扫描器

	for scanner.Scan() { // 逐行扫描

		line := strings.TrimSpace(scanner.Text()) // 拿到一行内容
		parts := strings.SplitN(line, "=", 2)     // 以=分开

		if len(parts) == 2 { // 有键值对便当作环境变量使用
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])) // 设置环境变量
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

type Config struct { // AIAPI请求变量
	APIKey   string        `json:"apikey"`
	Model    string        `json:"model"`
	AIAPIurl string        `json:"apiurl`
	Method   string        `json:"method"`
	Timeout  time.Duration `json:"timeout"`
	Prompt   string        `json:"prompt"`
}

func LoadConfig(method string) *Config { // 装配AIAPI请求变量
	key, exists := os.LookupEnv("AI_API_KEY")
	if !exists {
		log.Fatal("未从环境中找到, AI_API_KEY装配失败, 缺失变量为关键值, 中断全局")
	}
	if key == "" {
		log.Fatal("环境配置异常, AI_API_KEY环境变量名被配置, 但装配时得到空值, 中断全局")
	}
	model, exists := os.LookupEnv("MODEL")
	if !exists {
		log.Println("未从环境中找到, MODELY装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		model = "deepseek-v4-pro"
	}
	aiAPIurl, exists := os.LookupEnv("AI_API_URL")
	if !exists {
		log.Println("未从环境中找到, AI_API_URL装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		aiAPIurl = "https://api.deepseek.com/chat/completions"
	}
	prompt, exists := os.LookupEnv("SYSTEM_PROMPT")
	if !exists {
		log.Println("未从环境中找到, SYSTEM_PROMPT装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		prompt = "你的名字叫AnoX,与你对话的是狰和,对话者只有你和他.你的性格清冷,对待专业问题会用专业术语的回复的同时细心讲解每个知识,对待日常问题使用人类口语化回复,不使用第三人称描述自己"
	}
	timeoutStr, exists := os.LookupEnv("TIMEOUT")
	if !exists {
		log.Println("未从环境中找到, TIMEOUT装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		timeoutStr = "60s"
	}

	timeout, err := time.ParseDuration(timeoutStr) // 转换超时配置
	if err != nil {
		log.Println("<超时>变量转换失败, 缺失变量为非关键值, 后续执行使用默认值")
		timeout = 60 * time.Second
	}

	return &Config{
		APIKey:   key,
		Model:    model,
		AIAPIurl: aiAPIurl,
		Method:   method,
		Timeout:  timeout,
		Prompt:   prompt,
	}
}

func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

type ResponseJSON struct { // 先只拿响应文本信息
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (a *Agent) call_llm(input string) (ResponseJSON, error) { //发送请求并处理响应
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Timeout) // 创建当前请求超时
	defer cancel()

	input = strings.TrimSpace(input) // 对输入简单处理

	message := []any{ // 组装信息
		map[string]string{
			"content": a.config.Prompt,
			"role":    "system",
		},
		map[string]string{
			"content": input,
			"role":    "user",
		},
	}

	reqBody := map[string]any{ // 组装请求体
		"messages":    message,
		"model":       a.config.Model,
		"temperature": 0.5,
	}

	data, err := json.Marshal(reqBody) //将请求体编码
	if err != nil {
		log.Println("call_llm内请求体编码失败, 返回原始err, 中断当前函数")
		return ResponseJSON{}, err
	}

	req, err := http.NewRequestWithContext(ctx, a.config.Method, a.config.AIAPIurl, bytes.NewReader(data)) // 创建完整请求
	if err != nil {
		log.Println("call_llm内请求创建失败, 返回原始err, 中断当前函数")
		return ResponseJSON{}, err
	}
	req.Header.Set("Content-Type", "application/json")         // 告诉对方我发送的是json
	req.Header.Set("Accept", "application/json")               // 希望对方响应json
	req.Header.Set("Authorization", "Bearer "+a.config.APIKey) // 拼接密钥

	res, err := a.client.Do(req) // 发送请求
	if err != nil {
		log.Println("call_llm内请求发送失败, 返回原始err, 中断当前函数")
		return ResponseJSON{}, err
	}
	defer res.Body.Close()

	if res.Body == nil {
		return ResponseJSON{}, fmt.Errorf("call_llm内响应体为空")
	}

	respBody, err := io.ReadAll(res.Body) // 读取响应请求体
	if err != nil {
		log.Println("call_llm内响应读取失败, 返回原始err, 中断当前函数")
		return ResponseJSON{}, err
	}

	if res.StatusCode != http.StatusOK {
		log.Println("call_llm内响应状态非正确, 返回err, 中断当前函数")
		return ResponseJSON{}, fmt.Errorf("响应状态异常, 得到状态: %d", res.StatusCode)
	}

	var responseJson ResponseJSON
	if err := json.Unmarshal(respBody, &responseJson); err != nil {
		log.Println("call_llm内请求体解析失败, 返回err, 中断当前函数")
		return ResponseJSON{}, err
	}

	return responseJson, nil
}
func main() {
	if err := loadEnv("agent_lab/.env.local"); err != nil { // 配置环境变量
		log.Fatalf("设置环境变量失败, 后续未执行. 详情: %s", err)
	}
	fmt.Println("环境配置成功") // 提示

	config := LoadConfig("POST")
	timeout := 120 * time.Second
	client := NewClient(timeout)

	agent := NewAgent(client, config)

	fmt.Println("输入exit退出:")
	scanners := bufio.NewScanner(os.Stdin) // 读取输入内容

	var scannerChance int // 允许的错误计数

	for {
		if scannerChance == 3 {
			break
		}

		if !scanners.Scan() {
			break
		}

		input := strings.TrimSpace(scanners.Text())

		if input == "exit" {
			break
		}

		resp, err := agent.call_llm(input)
		if err != nil {
			scannerChance++
			log.Println(err)
			continue
		}

		if len(resp.Choices) > 0 {
			fmt.Printf("AI: %s\n", resp.Choices[0].Message.Content)
		}
	}

	if err := scanners.Err(); err != nil {
		log.Fatalf("循环程序异常退出, 详情: %s", err)
	} else {
		log.Println("循环程序正常退出")
	}
}
