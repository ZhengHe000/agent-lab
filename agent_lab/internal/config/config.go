package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultModelName    = "deepseek-v4-flash"
	defaultEndpoint     = "https://api.deepseek.com/chat/completions"
	defaultTimeout      = "60s"
	defaultSystemPrompt = "你是智能专家Noah, 以专业角度和成熟冷静的逻辑处理问题和对话。"
)

type Config struct {
	Model ModelConfig
	Agent AgentConfig
}

type ModelConfig struct {
	APIKey   string
	Name     string
	Endpoint string
	Timeout  time.Duration
}

type AgentConfig struct {
	SystemPrompt string
}

func envOrDefault(key, defaultValue string) string { // 辅助函数, 当一项[环境配置]可以接受[默认选项]时使用该函数
	envValue := strings.TrimSpace(os.Getenv(key))
	if envValue == "" {
		return defaultValue
	}

	return envValue
}

func LoadConfig() (Config, error) {
	apiKey := strings.TrimSpace(os.Getenv("AI_API_KEY"))
	if apiKey == "" {
		return Config{}, fmt.Errorf("[环境变量] AI_API_KEY 未配置")
	}

	modelName := envOrDefault("MODEL", defaultModelName)
	modelEndpoint := envOrDefault("LLM_API_URL", defaultEndpoint)
	modelTimeoutText := envOrDefault("REQUEST_TIMEOUT", defaultTimeout)
	prompt := envOrDefault("SYSTEM_PROMPT", defaultSystemPrompt)

	timeout, err := time.ParseDuration(modelTimeoutText)
	if err != nil {
		return Config{}, fmt.Errorf("[环境变量] REQUEST_TIMEOUT 解析失败: %w ", err)
	}

	config := Config{
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
