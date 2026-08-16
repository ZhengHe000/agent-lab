package workspace

import (
	"testing"
)

func TestIsAllowedTextFile(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid_.go_suffix",
			input: "internal/config/config.go",
			want:  true,
		},
		{
			name:  "valid_basePath go.mod",
			input: "nested/go.mod",
			want:  true,
		},
		{
			name:  "valid_basePath go.sum",
			input: "go.sum",
			want:  true,
		},
		{
			name:  "unsupported_majuscule",
			input: "main.GO",
			want:  false,
		},
		{
			name:  "unsupported_suffix",
			input: "assets/logo.png",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedTextFile(tt.input)
			if result != tt.want {
				t.Fatalf("want: %v, got: %v", tt.want, result)
			}
		})
	}

	t.Run("Effective_expansion_name_whitelist", func(t *testing.T) {
		whitelistTests := []string{
			"main.go",
			"README.md",
			"note.txt",
			"config.json",
			"config.yaml",
			"config.yml",
			"config.toml",
		}

		for _, toolPath := range whitelistTests {
			if !isAllowedTextFile(toolPath) {
				t.Fatalf("want true, but %v is false", toolPath)
			}
		}
	})
}
