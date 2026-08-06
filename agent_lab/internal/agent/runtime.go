package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZhengHe000/agent-lab/agent_lab/internal/model"
	"github.com/ZhengHe000/agent-lab/agent_lab/internal/tool"
)

const maxModelSteps = 8

type Runtime struct {
	llm       model.Model     // 使用模型的能力
	modelName string          // 模型名称
	messages  []model.Message // 上下文历史
	tools     *tool.Registry  // 持有注册表
}

func NewRuntime(llm model.Model, modelName string, systemPrompt string, tools *tool.Registry) (*Runtime, error) {
	if llm == nil { // 检查模型
		return nil, fmt.Errorf("model 不能为空")
	}

	if strings.TrimSpace(modelName) == "" { // 检查模型名称
		return nil, fmt.Errorf("模型名 不能为空")
	}

	if strings.TrimSpace(systemPrompt) == "" { // 检查系统提示词
		return nil, fmt.Errorf("系统提示词 不能为空")
	}

	if tools == nil {
		return nil, fmt.Errorf("工具注册表不能为空")
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
		tools: tools,
	}

	return runtime, nil
}

func (r *Runtime) RunTurn(ctx context.Context, input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" { // 检查用户输入
		return "", fmt.Errorf("input 不能为空")
	}

	workingMessages := make([]model.Message, 0, len(r.messages)+1) // 创建容器
	workingMessages = append(workingMessages, r.messages...)       // 将正式历史当作上下文
	workingMessages = append(workingMessages, model.Message{       // 将本轮输入追加至末尾
		Role:    model.RoleUser,
		Content: input,
	})

	for stop := 0; stop < maxModelSteps; stop++ { // 进入轮次受控循环
		response, err := r.llm.Complete(ctx, model.Request{ // 调用模型回复
			Model:    r.modelName,
			Messages: workingMessages,
			Tools:    r.tools.Definitions(),
		})
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrModelInvocationFailed, err)
		}

		message := response.Message
		if message.Role != model.RoleAssistant { // 检查响应角色, 必须是 model.RoleAssistant
			return "", fmt.Errorf("%w: want %q, got: %q", ErrResponseRoleError, model.RoleAssistant, message.Role)
		}

		if len(message.ToolCalls) == 0 {
			if strings.TrimSpace(message.Content) == "" {
				return "", ErrEmptyContent
			}

			workingMessages = append(workingMessages, message)
			r.messages = workingMessages
			return message.Content, nil
		}

		workingMessages = append(workingMessages, message)

		for _, call := range message.ToolCalls {
			if strings.TrimSpace(call.Name) == "" {
				return "", fmt.Errorf("%w: 工具名为空", ErrInvalidToolCall)
			}

			if strings.TrimSpace(call.ID) == "" {
				return "", fmt.Errorf("%w: 工具调用ID为空", ErrInvalidToolCall)
			}

			result := r.executeToolCall(ctx, call)

			workingMessages = append(workingMessages, model.Message{
				Role:       model.RoleTool,
				Content:    result,
				ToolCallID: call.ID,
			})
		}
	}

	return "", ErrMaxStepsExceeded
}

func (r *Runtime) executeToolCall(ctx context.Context, call model.ToolCall) string {
	registeredTool, exists := r.tools.Get(call.Name)
	if !exists {
		return fmt.Sprintf("工具执行失败, 未注册工具: %q", call.Name)
	}

	result, err := registeredTool.Execute(ctx, call.Arguments)
	if err != nil {
		return fmt.Sprintf("工具 %q 执行失败: %v", call.Name, err)
	}

	return result
}
