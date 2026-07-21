package main

import (
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/config"
	"log"
)

func main() {
	if err := config.LoadEnv("agent_lab/.env.local"); err != nil {
		log.Fatalf("环境配置失败:%v", err)
	}
	log.Println("环境配置成功")
}
