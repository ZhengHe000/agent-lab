package agent

import (
	"errors"
)

var (
	// 调用模型时出错
	ErrModelInvocationFailed = errors.New("调用模型失败")
)
