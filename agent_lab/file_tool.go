package main

import (
	"fmt"
	"unicode/utf8"
)

type ByteTooLongError struct { // 字节过长错误
	limit int
	actual int 
}
func (t *TooLongError) Error()string {
	return fmt.Sprintf("内容过长, 限制 %d 字节, 当前 %d 字节", t.limit, t.actual)
}

type RuneTooLongError struct { // 字符过长错误
	limit int
	actual int 
}
func (r *RuneTooLongError) Error()string {
	return fmt.Sprintf("内容过长, 限制 %d 字符, 当前 %d 字符", r.limit, r.actual)
}

func validateContent(content string)error { // 内容检验

	if content = "" {
		return ErrContentEmpty
	}

	if charcter := utf8.RuneCountInString(content); charcter > 10000 { // 判断字符
		return &RuneTooLongError{
			limit: 10000,
			actual: charcter,
		}
	}

	if size := len(content); size > 40000 { // 判断字符
		return &ByteTooLongError{
			limit: 40000,
			actual: size,
		}
	}

	return nil
}