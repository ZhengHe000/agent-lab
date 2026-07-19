package main

import (
	"errors"
)

var (
	// 一个string类型 的 内容为""(也就是空)时 应该使用的错误
	ErrContentEmpty = errors.New("内容不能为空")

	// 用于表示内容格式与预设的正则表达式规则不匹配时返回的错误
	ErrInvalidFilename = errors.New("内容格式与正则规则不符")

	// 当 `os.MkdirAll` 因参数错误而失败时，应返回此错误. 若需保留底层原始错误，建议在调用处使用 `fmt.Errorf("%w: %w", ErrMkdirAll, err)` 进行包装。
	ErrMkdirAll = errors.New("调用 os.MkdirAll 失败")

	// 当 `os.WriteFile` 因参数错误而失败时，应返回此错误. 若需保留底层原始错误，建议在调用处使用 `fmt.Errorf("%w: %w", ErrWriteFile, err)` 进行包装。
	ErrWriteFile = errors.New("调用 os.WriteFile 失败")

	// 当用户或逻辑选择拒绝继续写入操作时，返回此错误
	ErrWriteCancelled = errors.New("写入操作已取消")

	// 当输入内容无效（如期望 `y` 或 `n` 却得到其他值）时，返回此错误
	ErrInvalidInput = errors.New("无效输入，请输入 y 或 n")

	// ErrReadFile 表示读取文件时失败
	ErrReadFile = errors.New("读取文件失败")

	// 无效的文件类型
	ErrInvalidFileType = errors.New("文件含有被禁止内容")

	// ErrFileTooLarge 表示待读取文件超过允许的最大大小。
	ErrFileTooLarge = errors.New("文件内容超过允许读取的大小")
)
