package main

import (
	"errors"
)

var (
	// 一个string类型 的 内容为""(也就是空)时 应该使用的错误
	ErrContentEmpty = errors.New("错误, 因为存在内容为空")

	// 一个内容格式 与 设定的正则表达式规则 不符 应该使用的错误
	ErrInvalidFilename = errors.New("错误, 因为存在与设定的正则表达式规则不符的内容")

	// 使用os.MkdirAll时 应该返回的错误, 如果想看原始错误信息应该在编写函数时 使用其他办法收集 例如使用日志 或 fmt.Errorf("%w %w", ErrMkdirAll, err)
	ErrMkdirAll = errors.New("错误, 因为过程中用错误的参数调用了os.MkdirAll函数")

	// 使用os.WriteFile时 应该返回的错误, 如果想看原始错误信息应该在编写函数时 使用其他办法收集 例如使用日志 或 fmt.Errorf("%w %w",ErrWriteFile, err)
	ErrWriteFile = errors.New("错误, 因为过程中用错误的参数调用了os.WriteFile函数")

	// 选择拒绝继续时 返回此错误
	ErrWriteCancelled = errors.New("取消写入操作")

	// 使用无效的输入内容时 返回次错误
	ErrInvalidInput = errors.New("无效输入，请输入 y 或 n")
)
