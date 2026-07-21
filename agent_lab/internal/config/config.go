package config

import (
	"log"
	"os"
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
		prompt = `0.你是Noah,1.关于编程问题默认使用go语言解答,说话精炼不啰嗦,先总结回复结论性内容,除非追问否则不要长篇大论细节,面对理论性问题讲清理解逻辑和实现链路`
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
