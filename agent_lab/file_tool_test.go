package main

import (
	"testing"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		filename string // 输入的文件名
		wantErr  bool   // 结果是否错误
	}{
		{filename: "note.txt", wantErr: false},
		{filename: "study-note_01.txt", wantErr: false},
		{filename: "../secret.txt", wantErr: true},
		{filename: "a/b.txt", wantErr: true},
		{filename: "my file.txt", wantErr: true},
		{filename: "report.md", wantErr: true},
	}

	for i, tt := range tests {
		err := validateFilename(tt.filename)

		if (err != nil) != tt.wantErr {
			t.Errorf("case: %d, Filename: %s, Want: %t, Got: %t, err: %s", i, tt.filename, tt.wantErr, (err != nil), err)
		}
	}
}
