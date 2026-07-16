package main

import (
	"errors"
)

var (
	// 一个string类型 的内容为""(也就是空)时 应该使用的哨兵错误
	ErrContentEmpty = errors.New("错误, 因为存在内容为空")
)