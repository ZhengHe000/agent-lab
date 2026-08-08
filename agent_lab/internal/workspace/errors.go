package workspace

import (
	"errors"
)

var (
	// 一个string类型 的 内容为""(也就是空)时 应该使用的错误
	ErrContentEmpty = errors.New("内容不能为空")

	// 用于表示内容格式与预设的正则表达式规则不匹配时返回的错误
	ErrInvalidFilename = errors.New("内容格式与正则规则不符")

	// os.MkdirAll 失败时返回此错误
	ErrMkdirAll = errors.New("创建目录失败")

	// os.WriteFile 失败时返回此错误
	ErrWriteFile = errors.New("文件写入失败")

	// 当用户或逻辑选择拒绝继续写入操作时，返回此错误
	ErrWriteCancelled = errors.New("写入操作已取消")

	// 当输入内容无效（如期望 `y` 或 `n` 却得到其他值）时，返回此错误
	ErrInvalidInput = errors.New("无效输入，请输入 y 或 n")

	// ErrReadFile 表示读取文件时失败
	ErrReadFile = errors.New("读取文件失败")

	// ErrFileTooLarge 表示待读取文件超过允许的最大大小。
	ErrFileTooLarge = errors.New("文件内容超过允许读取的大小")
)
