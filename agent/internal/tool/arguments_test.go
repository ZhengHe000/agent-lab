package tool

import (
	"encoding/json"
	"testing"
)

func TestDecodeObjectArguments_rejects_invalid_object_arguments(
	t *testing.T,
) {
	type arguments struct {
		Path string `json:"path"`
	}

	tests := []struct {
		name      string
		arguments json.RawMessage
		wantError bool
	}{
		{
			name:      "valid object",
			arguments: json.RawMessage(`{"path":"README.md"}`),
			wantError: false,
		},
		{
			name:      "unknown field",
			arguments: json.RawMessage(`{"path":"README.md","unknown":true}`),
			wantError: true,
		},
		{
			name:      "non-object",
			arguments: json.RawMessage(`null`),
			wantError: true,
		},
		{
			name:      "multiple values",
			arguments: json.RawMessage(`{"path":"README.md"} {}`),
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := DecodeObjectArguments[arguments](
				test.arguments,
			)

			if (err != nil) != test.wantError {
				t.Fatalf(
					"DecodeObjectArguments() error = %v, wantError %v",
					err,
					test.wantError,
				)
			}

			if !test.wantError && got.Path != "README.md" {
				t.Fatalf(
					"DecodeObjectArguments().Path = %q, want %q",
					got.Path,
					"README.md",
				)
			}
		})
	}
}
