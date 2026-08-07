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
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/workspace"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := config.LoadEnvFile("agent_lab/.env.local"); err != nil { //加载环境变量
		return fmt.Errorf("加载环境文件失败: %w", err)
	}

	cfg, err := config.LoadConfig() // 提取环境变量
	if err != nil {
		return fmt.Errorf("加载程序配置失败: %w", err)
	}

	httpClient := &http.Client{ // 设置客户端
		Timeout: cfg.Model.Timeout,
	}

	client, err := openai.NewClient(httpClient, cfg.Model.Endpoint, cfg.Model.APIKey) // 组装客户端
	if err != nil {
		return fmt.Errorf("客户端配置失败: %w", err)
	}

	toolsRegistry, err := tool.NewRegistry( // 创建工具注册表
		workspace.NewReadTextFileTool(),
		workspace.NewListTextFilesTool(),
	)

	if err != nil {
		return fmt.Errorf("工具注册表创建失败: %w", err)
	}

	runtime, err := agent.NewRuntime(client, cfg.Model.Name, cfg.Agent.SystemPrompt, toolsRegistry) // 组装单轮运行器
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
			return nil
		}

		reply, err := runtime.RunTurn(context.Background(), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "本轮执行失败: %v\n", err)
			continue
		}

		fmt.Printf("Noah: %s\n", reply)
	}
}
