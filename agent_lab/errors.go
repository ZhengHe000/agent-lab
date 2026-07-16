package main

import (
	"errors"
)

var (
	// 一个string类型 的 内容为""(也就是空)时 应该使用的错误
	ErrContentEmpty = errors.New("错误, 因为存在内容为空")

	// 一个内容格式 与 设定的正则表达式规则 不符 应该使用的错误
	ErrInvalidFilename = errors.New("错误, 因为存在与设定的正则表达式规则不符的内容")
)