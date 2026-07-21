package workspace

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestConfirmWrite(t *testing.T) {
	tests := []struct {
		testName     string
		readerTest   io.Reader
		writerTest   *bytes.Buffer
		filenameTest string
		contentTest  string
		wantBool     bool
		wantErr      bool
	}{
		{testName: "y确认",
			readerTest:   strings.NewReader("y\n"),
			writerTest:   &bytes.Buffer{},
			filenameTest: "note.txt",
			contentTest:  "hello",
			wantBool:     true,
			wantErr:      false,
		},
		{testName: "n拒绝",
			readerTest:   strings.NewReader("n\n"),
			writerTest:   &bytes.Buffer{},
			filenameTest: "note.txt",
			contentTest:  "hello",
			wantBool:     false,
			wantErr:      false,
		},
		{testName: "maybe返回错误",
			readerTest:   strings.NewReader("maybe\n"),
			writerTest:   &bytes.Buffer{},
			filenameTest: "note.txt",
			contentTest:  "hello",
			wantBool:     false,
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			gotBool, gotErr := confirmWrite(tc.readerTest, tc.writerTest, tc.filenameTest, tc.contentTest)

			if gotBool != tc.wantBool {
				t.Errorf("want: %v, got: %v", tc.wantBool, gotBool)
			}

			if (gotErr != nil) != tc.wantErr {
				t.Errorf("want: %v, got: %v", tc.wantErr, gotErr)
			}

			if !strings.Contains(tc.writerTest.String(), "note.txt") {
				t.Errorf("提示信息未包含文件名")
			}

			if !strings.Contains(tc.writerTest.String(), "hello") {
				t.Errorf("提示信息未包含内容预览")
			}
			if !strings.Contains(tc.writerTest.String(), "确认写入") {
				t.Errorf("提示信息未包含确认")
			}
		})
	}
}
