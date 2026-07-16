package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type ByteTooLongError struct { // 字节过长错误
	limit  int
	actual int
}

func (t *ByteTooLongError) Error() string {
	return fmt.Sprintf("内容过长, 限制 %d 字节, 当前 %d 字节", t.limit, t.actual)
}

type RuneTooLongError struct { // 字符过长错误
	limit  int
	actual int
}

func (r *RuneTooLongError) Error() string {
	return fmt.Sprintf("内容过长, 限制 %d 字符, 当前 %d 字符", r.limit, r.actual)
}

func validateContent(content string) error { // Content内容校验

	if strings.TrimSpace(content) == "" { // 判断是否为空
		return ErrContentEmpty
	}

	if charcter := utf8.RuneCountInString(content); charcter > 10000 { // 判断字符
		return &RuneTooLongError{
			limit:  10000,
			actual: charcter,
		}
	}

	if size := len(content); size > 40000 { // 判断字符
		return &ByteTooLongError{
			limit:  40000,
			actual: size,
		}
	}

	return nil
}

var filenameRule = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,95}\.txt$`) // 编译可复用的规则

func validateFilename(filename string) error { // 文件名校验
	if !filenameRule.MatchString(filename) { //用正则表达式规则判断 传入值是否合规
		return ErrInvalidFilename
	}

	return nil
}
