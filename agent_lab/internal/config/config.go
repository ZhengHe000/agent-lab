package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Model ModelConfig
	Agent AgentConfig
}

type ModelConfig struct {
	APIKey   string        // 模型API密钥
	Name     string        // 模型名称
	Endpoint string        // 模型API的完整URL路径
	Timeout  time.Duration // 客户端超时
}

type AgentConfig struct {
	SystemPrompt string // 模型系统提示词
}

func envOrDefault(key, defaultValue string) string { // 辅助函数, 当一项[环境配置]可以接受[默认选项]时使用该函数
	envValue := strings.TrimSpace(os.Getenv(key)) // 读取环境变量的key参数的值
	if envValue == "" {                           // 如果[环境变量]中该值为空, 使用传入的[defaultValue参数]作为默认值
		return defaultValue
	}

	return envValue
}

func LoadConfig() (Config, error) {
	apiKey := strings.TrimSpace(os.Getenv("AI_API_KEY")) // 获取环境变量的[AI_API_KEY]值
	if apiKey == "" {
		return Config{}, fmt.Errorf("[环境变量] AI_API_KEY 未配置")
	}

	modelName := envOrDefault("MODEL", "deepseek-v4-pro")                                                                                                                                                       // 使用存在默认值的辅助函数获取[MODEL]值
	modelEndpoint := envOrDefault("LLM_API_URL", "https://api.deepseek.com/chat/completions")                                                                                                                   // 使用存在默认值的辅助函数获取[LLM_API_URL]值
	modelTimeoutText := envOrDefault("REQUEST_TIMEOUT", "60s")                                                                                                                                                  // 使用存在默认值的辅助函数获取[REQUEST_TIMEOUTL]值
	prompt := envOrDefault("SYSTEM_PROMPT", `你是智能生命体:Noah|小名:诺诺,你的思维逻辑清晰,不受情绪影响,对人类各领域技能实践经验为满级,1.提问模糊时直接反问最多3个解决问题需要的关键变量,2.编程默认Golang,优先回复结论性内容(可直接运行的代码片段或具体方案),不主动发散解释底层原理,3.被追问细节时按照工程拓展性最优选择讲清设计逻辑和代码实现`) // 使用存在默认值的辅助函数获取[SYSTEM_PROMPT]值

	timeout, err := time.ParseDuration(modelTimeoutText) // 将得到的string类型的Timeout 解析为time.Duration类型
	if err != nil {
		return Config{}, fmt.Errorf("[环境变量] REQUEST_TIMEOUT 解析失败, 错误: %w ", err)
	}

	config := Config{ // 组装Config
		Model: ModelConfig{
			APIKey:   apiKey,
			Name:     modelName,
			Endpoint: modelEndpoint,
			Timeout:  timeout,
		},
		Agent: AgentConfig{
			SystemPrompt: prompt,
		},
	}

	return config, nil
}
