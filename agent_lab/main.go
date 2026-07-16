package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if err := loadEnv("agent_lab/.env.local"); err != nil { // 配置环境变量
		log.Fatalf("设置环境变量失败, 后续未执行. 详情: %s", err)
	}
	fmt.Println("环境配置成功") // 提示

	config := LoadConfig("POST")
	httpClientTimeout := 120 * time.Second
	client := NewClient(httpClientTimeout)

	agent := NewAgent(client, config)

	fmt.Println("输入exit退出:")
	scanners := bufio.NewScanner(os.Stdin) // 读取输入内容

	var requestFailureCount int // 允许的错误计数

	for {
		if requestFailureCount == 3 {
			break
		}

		if !scanners.Scan() {
			break
		}

		input := strings.TrimSpace(scanners.Text())

		if input == "exit" {
			break
		}

		resp, err := agent.callLLM(input)
		if err != nil {
			requestFailureCount++
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
