package config

import (
	"bufio"
	"fmt"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/workspace"
	"os"
	"strings"
)

func LoadEnvFile(fileName string) error { // 配置环境变量
	file, err := os.Open(fileName) //打开env环境配置文件
	if err != nil {
		return fmt.Errorf("%w:%w", workspace.ErrReadFile, err)
	}
	defer file.Close() // 结束前关闭文件

	scanner := bufio.NewScanner(file) // 创建扫描器
	lineNumber := 0                   // 用来计算当前第几行

	for scanner.Scan() { // 逐行扫描
		lineNumber++

		line := strings.TrimSpace(scanner.Text()) // 拿到一行内容

		if line == "" || strings.HasPrefix(line, "#") { // 跳过空行和注释
			continue
		}

		parts := strings.SplitN(line, "=", 2) // 以=分开

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" { // 检查配置名称
			return fmt.Errorf("第 %d 行的环境变量名称为空", lineNumber)
		}

		if _, exists := os.LookupEnv(key); exists { // 查看配置是否存在
			continue
		}

		if err := os.Setenv(key, value); err != nil { // 进行环境配置
			return fmt.Errorf("[环境配置] 第 %d 行失败:%w", lineNumber, err)
		}

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("[环境配置|扫描器异常]:%w", err)
	}

	return nil
}
