package unit

import (
	"context"
	toolModels "gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
)

func UnitMainWorker(ctx context.Context, work *toolModels.Work, validParams *schemas.UnitParams) error {
	return nil
}
