package model

import (
	"encoding/json"
	"context"
)

type Role string

const (
	RoleSystem Role = "system" // 系统级
	RoleUser Role = "user" // 用户级
	RoleAssistant Role = "assistant" // 模型级
	RoleTool Role = "tool" // 工具级
)

type ToolCall struct { // 模型使用工具时需要生成的信息
	ID string // 具体工具的调用顺序 例如:(查目录ID:1, 调浏览器ID:2, 读取链接文档ID:3)
	Name string // 工具名
	Arguments json.RawMessage // 模型对具体工具按要求生成的操作参数, 使用json.RawMessage方便使用具体工具前 解析到当前工具的参数结构体
}

type Message struct { // 通用消息
	Role Role // 身份
	Content string // 正文, 不使用*string, 因为是否为null应该在上层做判断, 这里直接判断是否为""
	ToolCalls []ToolCall // 单个或多工具信息
	ToolCallID string // 工具函数返回时携带的自身信息, 告诉模型一段工具调用结果来自什么具体工具
}

type Request struct {
	Model string // 模型名
	Messages []Message 
}

type Response struct {
	Messages Message
	FinishReason string // 模型结束循环的理由
}

type Model interface { // 定义Model接口
	Complete(ctx context.Context, request Request)(Response, error) // 只要实现Complete函数作为方法, 就可以当作Model接口使用
}