package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/agent"
	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/config"
	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/model/openai"
	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/terminal"
	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/tool"
	"github.com/ZhengHeOwo/agent-lab/agent_lab/internal/workspace"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	console, err := terminal.NewConsole(os.Stdin, os.Stdout)
	if err != nil {
		return fmt.Errorf("创建终端交互对象失败: %w", err)
	}

	if err := config.LoadEnvFile("agent_lab/.env.local"); err != nil {
		return fmt.Errorf("加载环境文件失败: %w", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载程序配置失败: %w", err)
	}

	httpClient := &http.Client{
		Timeout: cfg.Model.Timeout,
	}

	client, err := openai.NewClient(httpClient, cfg.Model.Endpoint, cfg.Model.APIKey)
	if err != nil {
		return fmt.Errorf("客户端配置失败: %w", err)
	}

	writeTextFileTool, err := workspace.NewWriteTextFileTool(console)
	if err != nil {
		return fmt.Errorf("创建文本写入工具失败: %w", err)
	}

	toolsRegistry, err := tool.NewRegistry(
		workspace.NewReadTextFileTool(),
		workspace.NewListTextFilesTool(),
		writeTextFileTool,
	)

	if err != nil {
		return fmt.Errorf("工具注册表创建失败: %w", err)
	}

	runtime, err := agent.NewRuntime(client, cfg.Model.Name, cfg.Agent.SystemPrompt, toolsRegistry)
	if err != nil {
		return fmt.Errorf("创建Agent运行器失败: %w", err)
	}

	fmt.Println("Agent Lab 已启动, 输入 exit 退出")
	for {
		input, err := console.ReadLine(": ")
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("读取终端失败: %w", err)
		}

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

		fmt.Printf("安悬: %s\n", reply)
	}
}
