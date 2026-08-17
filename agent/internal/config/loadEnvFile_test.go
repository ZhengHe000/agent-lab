package config

import (
	"os"
	"path/filepath"
	"testing"
)

func createEnvFile(t *testing.T, dir string, fileName string, envContent string) string {
	t.Helper()

	filePath := filepath.Join(dir, fileName)
	if err := os.WriteFile(filePath, []byte(envContent), 0644); err != nil {
		t.Fatalf("创建配置测试文件失败:%v", err)
	}

	return filePath
}
func TestLoadEnvFile(t *testing.T) {
	tests := []struct {
		name             string
		fileName         string
		envContent       string
		wantErr          bool
		notExists        string
		fixed_allocation string
		fixed_value      string
	}{
		{
			name:             "使用错误文件路径",
			fileName:         ".envGG",
			envContent:       "NUMBER=0",
			wantErr:          true,
			notExists:        "NUMBER",
			fixed_allocation: "ZERO",
			fixed_value:      "0",
		},
		{
			name:     "是否跳过空行/注释",
			fileName: ".env.text",
			envContent: `
			#这是=注释

			#这是=注释
						
			`,
			wantErr:          false,
			notExists:        "#这是",
			fixed_allocation: "ZERO",
			fixed_value:      "0",
		},
		{
			name:             "使用空白key",
			fileName:         ".env.text",
			envContent:       "  =0",
			wantErr:          true,
			notExists:        "  ",
			fixed_allocation: "ZERO",
			fixed_value:      "0",
		},
		{
			name:             "是否更改旧配置",
			fileName:         ".env.text",
			envContent:       "ZERO=1",
			wantErr:          false,
			notExists:        "",
			fixed_allocation: "ZERO",
			fixed_value:      "0",
		},
		{
			name:             "缺少等于号",
			fileName:         ".env.text",
			envContent:       "ZERO1",
			wantErr:          true,
			notExists:        "",
			fixed_allocation: "ZERO",
			fixed_value:      "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Setenv(tt.fixed_allocation, tt.fixed_value)

			testDir := t.TempDir()

			var filePath string
			if tt.name == "使用错误文件路径" {
				filePath = filepath.Join(testDir, "nonexistent", tt.fileName)
			} else {
				filePath = createEnvFile(t, testDir, tt.fileName, tt.envContent)
			}

			err := LoadEnvFile(filePath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("错误与期望不同")
			}

			if tt.notExists != "" {
				if _, exists := os.LookupEnv(tt.notExists); exists {
					t.Fatalf("不期望配置的变量意外配置成功, notExists: %s", tt.notExists)
				}
			}
			if v, _ := os.LookupEnv("ZERO"); v != "0" {
				t.Fatalf("默认环境配置在测试中被修改, 默认值: 0,当前值: %s", v)
			}
		})
	}
}
