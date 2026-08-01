package agent

import (
	"errors"
)

var (
	// 调用模型时出错
	ErrModelInvocationFailed = errors.New("调用模型失败")

	// 响应角色错误
	ErrResponseRoleError = errors.New("响应信息角色错误")

	// 暂不支持工具调用
	ErrToolCallsUnsupported = errors.New("工具调用暂不支持")

	//响应Content为空
	ErrEmptyContent = errors.New("响应内容为空")
)
