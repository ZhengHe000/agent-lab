package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func loadEnv(fileName string) error { // 配置环境变量
	file, err := os.Open(fileName) //打开env环境配置文件
	if err != nil {
		return log.Errorf("打开 %v 文件失败, 错误: %w", fileName, err)
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
