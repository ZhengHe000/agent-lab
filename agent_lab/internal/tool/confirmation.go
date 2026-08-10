package tool

import "context"

// ConfirmationRequest 描述需要用户明确授权的副作用操作。
type ConfirmationRequest struct {
	Action  string
	Summary string
	Details string
}

// Confirmer 从模型之外获取用户授权。
//
// 实现不得允许模型自行构造确认结果，授权必须来自用户或其他可信主体。
type Confirmer interface {
	Confirm(ctx context.Context, request ConfirmationRequest) (bool, error)
}
