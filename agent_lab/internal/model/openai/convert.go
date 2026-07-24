package openai

import (
	"encoding/json"
	"fmt"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

func toModelResponse(response chatCompletionResponse) (model.Response, error) {
	if len(response.Choices) == 0 {
		return model.Response{}, fmt.Errorf("模型响应Choices为空")
	}

	choices := response.Choices[0] // 减少后续重复

	modelMessage, err := toModelMessage(choices.Message) // 获取转换后model包类型的Message
	if err != nil {
		return model.Response{}, err
	}

	modelResponse := model.Response{ // 组装, 而非直接return, 为后续可能的逻辑增加灵活性
		Message:      modelMessage,
		FinishReason: choices.FinishReason,
	}

	return modelResponse, nil
}

func toModelMessage(chatMessage chatMessage) (model.Message, error) {
	content := ""                   // 定义文本容器, 默认为空""
	if chatMessage.Content != nil { // 当content字段指向目标不为nil, 说明它指向一个地址, 下面使用*解引用拿到该地址的值
		content = *chatMessage.Content
	}

	var role model.Role // 定义角色容器, 通过白名单校验
	switch chatMessage.Role {
	case "system":
		role = model.RoleSystem
	case "user":
		role = model.RoleUser
	case "assistant":
		role = model.RoleAssistant
	case "tool":
		role = model.RoleTool
	default:
		return model.Message{}, fmt.Errorf("Role类型不符合预设: %q", chatMessage.Role)
	}

	modelMessage := model.Message{ // 填充基础响应格式
		Role:       role,
		Content:    content,
		ToolCallID: chatMessage.ToolCallID,
	}

	if len(chatMessage.ToolCalls) == 0 { // 当工具调用信息长度为0时, 上面两个字段就是模型全部回复, 直接返回
		return modelMessage, nil
	}

	toolCalls := make([]model.ToolCall, 0, len(chatMessage.ToolCalls))

	for _, chatToolCall := range chatMessage.ToolCalls { // 遍历工具信息切片, 对每份信息校验后append加入model切片
		switch chatToolCall.Type { // 使用白名单, 方便日后拓展新工具类型
		case "function":

		default:
			return model.Message{}, fmt.Errorf("错误的工具类型: %q", chatToolCall.Type)
		}

		modelID := chatToolCall.ID
		modelName := chatToolCall.Function.Name
		modelArguments := json.RawMessage(chatToolCall.Function.Arguments)

		toolCalls = append(toolCalls, model.ToolCall{ // 向切片追加转换后的工具信息
			ID:        modelID,
			Name:      modelName,
			Arguments: modelArguments,
		})
	}

	modelMessage.ToolCalls = toolCalls // 将结果赋值给ToolCalls字段

	return modelMessage, nil
}
