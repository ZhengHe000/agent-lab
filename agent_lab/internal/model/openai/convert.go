package openai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
)

func toModelResponse(response chatCompletionResponse) (model.Response, error) { // 将外部响应转换成内部结构
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

func toChatCompletionRequest(request model.Request) (chatCompletionRequest, error) { // 将内部请求转换成外部请求
	if len(request.Messages) == 0 {
		return chatCompletionRequest{}, fmt.Errorf("内部请求信息的Messages为空, 无法转换为外部请求")
	}

	reqMessages := make([]chatMessage, 0, len(request.Messages)) // 创建变量容器
	for _, modelReqMessage := range request.Messages {           // 遍历内部信息切片, 每轮使用toChatMessage函数处理转换
		chatMessage, err := toChatMessage(modelReqMessage)
		if err != nil {
			return chatCompletionRequest{}, err
		}
		reqMessages = append(reqMessages, chatMessage) // 转换成功追加进容器
	}

	var chatTools []chatToolDefinition
	if len(request.Tools) > 0 {
		chatTools = make([]chatToolDefinition, 0, len(request.Tools))

		for _, chatTool := range request.Tools {
			chatDefinition, err := toChatToolDefinition(chatTool)
			if err != nil {
				return chatCompletionRequest{}, err
			}

			chatTools = append(chatTools, chatDefinition)
		}
	}

	chatCompletionRequest := chatCompletionRequest{ // 拼接目标字段
		Model:    request.Model,
		Messages: reqMessages,
		Tools:    chatTools,
	}
	return chatCompletionRequest, nil
}

func toChatMessage(modelMessage model.Message) (chatMessage, error) {
	var content *string                                                 // 定义content字段的容器, 默认为nil
	if modelMessage.Content != "" || len(modelMessage.ToolCalls) == 0 { // 当内部此字段有值 或 不存在工具调用(此时空值就是模型输出) 使用内部content的地址作为值
		content = &modelMessage.Content
	}

	var role string            // 定义一个role容器
	switch modelMessage.Role { //使用白名单得到role结果 或 退出错误
	case model.RoleSystem:
		role = "system"
	case model.RoleUser:
		role = "user"
	case model.RoleAssistant:
		role = "assistant"
	case model.RoleTool:
		role = "tool"
	default:
		return chatMessage{}, fmt.Errorf("Role类型不符合预设: %q", modelMessage.Role)
	}

	chatToolCalls := make([]chatToolCall, 0, len(modelMessage.ToolCalls)) // 创建chatToolCalls容器, 储存和追加从内部信息遍历出来进行转换后的chatToolCall值
	for _, toolCall := range modelMessage.ToolCalls {
		chatToolCall := chatToolCall{
			ID:   toolCall.ID,
			Type: "function", // 默认写死function
			Function: chatFunctionCall{
				Name:      toolCall.Name,
				Arguments: string(toolCall.Arguments), //将内部信息强转换为string
			},
		}

		chatToolCalls = append(chatToolCalls, chatToolCall) // 将一个转换结果追加进chatToolCalls容器
	}
	chatMessage := chatMessage{ // 组装chatMessage结构体
		Role:       role,
		Content:    content,
		ToolCalls:  chatToolCalls,
		ToolCallID: modelMessage.ToolCallID,
	}

	return chatMessage, nil
}

func toChatToolDefinition(definition model.ToolDefinition) (chatToolDefinition, error) {
	name := strings.TrimSpace(definition.Name)
	if name == "" {
		return chatToolDefinition{}, fmt.Errorf("工具名称不能为空")
	}

	description := strings.TrimSpace(definition.Description)
	if description == "" {
		return chatToolDefinition{}, fmt.Errorf("工具 %q 的描述不能为空", name)
	}

	if !json.Valid(definition.Parameters) {
		return chatToolDefinition{}, fmt.Errorf("工具 %q 的参数 Schema 不是合法JSON", name)
	}

	return chatToolDefinition{
		Type: "function",
		Function: chatFunctionDefinition{
			Name:        name,
			Description: description,
			Parameters:  definition.Parameters,
		},
	}, nil
}
