package workspace

import (
	"testing"
)

func TestValidToolPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "允许包含子目录",
			input:   "internal/config/config.go",
			want:    "internal/config/config.go",
			wantErr: false,
		},
		{
			name:    "允许包含..",
			input:   "version..backup.go",
			want:    "version..backup.go",
			wantErr: false,
		},
		{
			name:    "禁止路径穿越",
			input:   "../secret.go",
			want:    "",
			wantErr: true,
		},
		{
			name:    "禁止静默规范化",
			input:   "internal//config.go",
			want:    "",
			wantErr: true,
		},
		{
			name:    "项目规划不同于fs.ValidPath",
			input:   ".",
			want:    "",
			wantErr: true,
		},
		{
			name:    "锁定跨平台协议边界",
			input:   `C:\project\main.go`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateToolPath(tt.input)
			if got != tt.want {
				t.Fatalf("want: %v, got: %v", tt.want, got)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("结果异常: %v", err)
			}
		})
	}
}
