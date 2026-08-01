package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/agent"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/config"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model/openai"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := config.LoadEnvFile("agent_lab/.env.local"); err != nil {
		return fmt.Errorf("加载环境文件失败: %w", err)
	}

	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载程序配置失败: %w", err)
	}

	httpClient := &http.Client{
		Timeout: config.Model.Timeout,
	}

	client, err := openai.NewClient(httpClient, config.Model.Endpoint, config.Model.APIKey)
	if err != nil {
		return fmt.Errorf("客户端配置失败: %w", err)
	}

	runtime, err := agent.NewRuntime(client, config.Model.Name, config.Agent.SystemPrompt)
	if err != nil {
		return fmt.Errorf("创建Agent运行器失败: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Agent Lab 已启动, 输入 exit 退出")
	for {
		fmt.Print("狰和: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("读取终端失败: %w", err)
			}

			return nil
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			continue
		}

		reply, err := runtime.RunTurn(context.Background(), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "本轮执行失败: %v\n", err)
			continue
		}

		fmt.Printf("Noah: %s\n", reply)
	}
}
