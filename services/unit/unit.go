package unit

import (
	"context"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	toolModels "gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/aiqicha"
	"sync"
)

func UnitMainWorker(ctx context.Context, work *toolModels.Work, validParams *schemas.UnitParams) error {
	logger.Info(fmt.Sprintf("关键词:【%s|%s】数据源：%s 数据字段：%s\n", validParams.KeyWord, validParams.CompanyID, validParams.ScanSource, validParams.SearchField))

	//validProxy, err := proxy.GetProxy()
	//if err != nil {
	//	logger.Info(err.Error())
	//}

	var wg sync.WaitGroup

	//爱企查
	if utils.IsInList("aqc", validParams.ScanSource) {
		if validParams.CompanyID == "" || (validParams.CompanyID != "" && utils.CheckPid(validParams.CompanyID) == "aqc") {
			wg.Add(1)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Warn(fmt.Sprintf("[QCC] ERROR: %v", err))
						wg.Done()
					}
				}()
				//查询企业信息
				res, ensOutMap := aiqicha.GetUnitInfoByPid(validParams)
				fmt.Println("------------>", res, ensOutMap)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	return nil
}
