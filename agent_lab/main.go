package main

import (
	"log"
)

func main() {
	if err := loadEnv("agent_lab/.env.local"); err != nil { // 配置环境变量
		log.Fatalf("设置环境变量失败, 后续未执行. 详情: %s", err)
	}
	log.Println("环境配置成功") // 提示

}
