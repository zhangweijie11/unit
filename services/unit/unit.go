package unit

import (
	"context"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	toolModels "gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/aiqicha"
	"gitlab.example.com/zhangweijie/unit/services/tianyancha"
	"sync"
)

func UnitMainWorker(ctx context.Context, work *toolModels.Work, validParams *schemas.UnitParams) error {
	logger.Info(fmt.Sprintf("关键词:【%s|%s】数据源：%s 数据字段：%s", validParams.KeyWord, validParams.CompanyID, validParams.ScanSource, validParams.ResultField))

	//validProxy, err := proxy.GetProxy()
	//if err != nil {
	//	logger.Info(err.Error())
	//}

	// 统一处理任务参数
	if validParams.ScanSource == nil {
		validParams.ScanSource = global.DefaultAllSource
	}
	if validParams.ResultField == nil {
		validParams.ResultField = global.CanSearchAllInfos
	}
	if validParams.Deep == 0 {
		validParams.Deep = 1
	}

	var wg sync.WaitGroup

	//爱企查
	if utils.IsInList(global.SourceAqc, validParams.ScanSource) {
		if validParams.CompanyID == "" || (validParams.CompanyID != "" && utils.CheckPid(validParams.CompanyID) == global.SourceAqc) {
			wg.Add(1)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Warn(fmt.Sprintf("[QCC] ERROR: %v", err))
						wg.Done()
					}
				}()
				//查询企业信息
				aiqicha.GetUnitInfoByPid(validParams)
				//res, ensOutMap := aiqicha.GetUnitInfoByPid(validParams)
				//fmt.Println("------------>", res, ensOutMap)
				wg.Done()
			}()
		}
	}

	//天眼查
	if utils.IsInList(global.SourceTyc, validParams.ScanSource) {
		if validParams.CompanyID == "" || (validParams.CompanyID != "" && utils.CheckPid(validParams.CompanyID) == global.SourceTyc) {
			wg.Add(1)
			if validParams.Cookies.Tianyancha == "" || validParams.Cookies.TianyanchaTycid == "" {
				logger.Warn("【TYC】MUST LOGIN 请补充天眼查COOKIE和tycId")
			}
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Warn(fmt.Sprintf("[TYC] ERROR: %v", err))
						wg.Done()
					}
				}()
				tianyancha.GetUnitInfoByPid(validParams)
				//res, ensOutMap := tianyancha.GetUnitInfoByPid(validParams)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	return nil
}
