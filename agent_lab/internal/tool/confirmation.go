package tool

import (
	"context"
)
type confirmationRequest struct { 
	Action string
	Summary string
	Details string
}

type Confirmer interface { // 确认接口 | 确认接口不应该放在Tool中让不需要确认的工具赘上一份无用代码
	Confirm(ctx context.Context, request confirmationRequest)(bool, error)
}