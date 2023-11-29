package main

import (
	"context"
	"encoding/json"
	"errors"
	tool "gitlab.example.com/zhangweijie/tool-sdk/cmd"
	toolGlobal "gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	toolSchemas "gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	toolModels "gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/option"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/unit"
)

type executorIns struct {
	toolGlobal.ExecutorIns
}

// ValidWorkCreateParams 验证任务参数
func (ei *executorIns) ValidWorkCreateParams(params map[string]interface{}) (err error) {
	var schema = new(schemas.UnitParams)
	err = toolSchemas.CustomBindSchema(params, schema, schemas.RegisterValidatorRule)
	if err != nil {
		return err
	} else {
		isExist := schemas.ValidParamsExist(schema.CompanyID, schema.KeyWord)
		if !isExist {
			return errors.New(schemas.WorkKeyIDErr)
		}
		// 验证扫描源是否超出范围
		for _, source := range schema.ScanSource {
			if !utils.IsInList(source, global.DefaultAllSource) {
				return errors.New(schemas.WorkScanSourceErr)
			}
		}
		for _, source := range schema.ResultField {
			if !utils.IsInList(source, global.CanSearchAllInfos) {
				return errors.New(schemas.WorkResultFieldErr)
			}
		}
	}
	return err
}

// ExecutorMainFunc 任务执行主函数（可自由发挥）
// params = map[string]interface{}{
// "work": &toolmodels.Work
// }
func (ei *executorIns) ExecutorMainFunc(ctx context.Context, params map[string]interface{}) error {
	global.InitToolConf()
	errChan := make(chan error, 2)
	go func() {
		defer close(errChan)
		work := params["work"].(*toolModels.Work)
		var validParams schemas.UnitParams
		err := json.Unmarshal(work.Params, &validParams)
		if err != nil {
			logger.Error(toolSchemas.JsonParseErr, err)
			errChan <- err
		} else {
			err = unit.UnitMainWorker(ctx, work, &validParams)
			errChan <- err
		}
	}()
	select {
	case <-ctx.Done():
		return errors.New(toolSchemas.WorkCancelErr)
	case err := <-errChan:
		return err
	}
}

func main() {
	defaultOption := option.GetDefaultOption()
	defaultOption.ExecutorIns = &executorIns{}
	tool.Start(defaultOption)
}
