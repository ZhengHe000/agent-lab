package model

import (
	"context"
	"encoding/json"
)

type Role string

const (
	RoleSystem    Role = "system"    // 系统
	RoleUser      Role = "user"      // 用户
	RoleAssistant Role = "assistant" // 模型
	RoleTool      Role = "tool"      // 工具
)

type ToolCall struct { // 模型使用工具时需要生成的信息
	ID        string          // 模型为一个工具调用分配的标识, 在工具调用结果封装时原封不动返回给模型, 模型看到这个标识就知道当前这份数据来自什么
	Name      string          // 工具名
	Arguments json.RawMessage // 模型对具体工具按要求生成的操作参数, 使用json.RawMessage方便使用具体工具前 解析到当前工具的参数结构体
}

type Message struct { // 通用消息
	Role       Role       // 身份
	Content    string     // 正文, 不使用*string, 因为是否为null应该在上层做判断, 这里直接判断是否为""
	ToolCalls  []ToolCall // 单个或多工具信息
	ToolCallID string     // 一份工具执行结果所对应的模型调用ID, 装的是模型决定调用工具的响应中 模型为一个工具操作分配的唯一ID
}

type Request struct {
	Model    string // 模型名
	Messages []Message
}

type Response struct {
	Message      Message
	FinishReason string // 模型决定结束循环时填写的理由
}

type Model interface { // 定义Model接口
	Complete(ctx context.Context, request Request) (Response, error) // 只要实现Complete函数作为方法, 就可以当作Model类型使用
}
