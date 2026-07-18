package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct { // AIAPI请求变量
	APIKey         string        `json:"apikey"`
	Model          string        `json:"model"`
	LLMAPIURL      string        `json:"apiurl"`
	Method         string        `json:"method"`
	RequestTimeout time.Duration `json:"timeout"`
	Prompt         string        `json:"prompt"`
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
	llmAPIURL, exists := os.LookupEnv("LLM_API_URL")
	if !exists {
		log.Println("未从环境中找到, LLM_API_URL装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		llmAPIURL = "https://api.deepseek.com/chat/completions"
	}
	prompt, exists := os.LookupEnv("SYSTEM_PROMPT")
	if !exists {
		log.Println("未从环境中找到, SYSTEM_PROMPT装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		prompt = `你是智能生命体Noah,中文直译为诺亚,关于编程问题默认使用go语言解答`
	}
	timeoutStr, exists := os.LookupEnv("REQUEST_TIMEOUT") // 局部超时逻辑增加自由度,函数内不使用默认http.Client设置兜底
	if !exists {
		log.Println("未从环境中找到, REQUEST_TIMEOUT装配失败, 缺失变量为非关键值, 后续执行使用默认值")
		timeoutStr = "60s"
	}

	requestTimeout, err := time.ParseDuration(timeoutStr) // 转换超时配置
	if err != nil {
		log.Println("<超时>变量转换失败, 缺失变量为非关键值, 后续执行使用默认值")
		requestTimeout = 60 * time.Second
	}

	return &Config{
		APIKey:         key,
		Model:          model,
		LLMAPIURL:      llmAPIURL,
		Method:         method,
		RequestTimeout: requestTimeout,
		Prompt:         prompt,
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
