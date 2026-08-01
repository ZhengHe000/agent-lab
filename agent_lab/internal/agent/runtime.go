package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

type Runtime struct {
	llm       model.Model     // 使用模型的能力
	modelName string          // 模型名称
	messages  []model.Message // 上下文历史
}

func NewRuntime(llm model.Model, modelName string, systemPrompt string) (*Runtime, error) {
	if llm == nil { // 检查模型
		return nil, fmt.Errorf("model 不能为空")
	}

	if strings.TrimSpace(modelName) == "" { // 检查模型名称
		return nil, fmt.Errorf("模型名 不能为空")
	}

	if strings.TrimSpace(systemPrompt) == "" { // 检查系统提示词
		return nil, fmt.Errorf("系统提示词 不能为空")
	}

	runtime := &Runtime{ // 组装Runtime并取地址
		llm:       llm,
		modelName: strings.TrimSpace(modelName),
		messages: []model.Message{
			model.Message{
				Role:    model.RoleSystem,
				Content: strings.TrimSpace(systemPrompt),
			},
		},
	}

	return runtime, nil
}

func (r *Runtime) RunTurn(ctx context.Context, input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" { // 检查用户输入
		return "", fmt.Errorf("input 不能为空")
	}

	candidate := make([]model.Message, 0, len(r.messages)+1) // 创建容器 存放model.Message
	candidate = append(candidate, r.messages...)             // 将messages的所有model.Message依次取出append追加进容器中
	candidate = append(candidate, model.Message{             // 将输入的问题[input]进行拼接, 以model.Message类型 追加进容器末尾
		Role:    model.RoleUser, // 身份为 用户
		Content: input,          // 输入的信息文本
	})

	modelRequest := model.Request{ // 使用模型名和candidate组成完整上下文且末尾是新输入
		Model:    r.modelName,
		Messages: candidate,
	}
	modelResponse, err := r.llm.Complete(ctx, modelRequest) // 调用Complete方法 得到model.Response和err信息
	if err != nil {
		return "", fmt.Errorf("%w:%w", ErrModelInvocationFailed, err)
	}

	msg := modelResponse.Message // 将modelResponse.Message简化为msg

	if msg.Role != model.RoleAssistant { // 检查回复身份是否是model.RoleAssistant[模型]
		return "", fmt.Errorf("%w want: %s, got: %s", ErrResponseRoleError, model.RoleAssistant, msg.Role)
	}

	if len(msg.ToolCalls) != 0 { // 当存在工具调用时直接返回 [暂时使用!]
		return "", ErrToolCallsUnsupported
	}

	if msg.Content == "" {
		return "", ErrEmptyContent
	}

	r.messages = append(candidate, msg) // 将响应model.Message追加进r.messages
	return msg.Content, nil             // 返回响应的Content文本部分和nil
}
