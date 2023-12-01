package aiqicha

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/result"
	"go.uber.org/zap"
	urlTool "net/url"
	"strconv"
	"strings"
)

// getCategoryInfoList 获取信息列表
func getCategoryInfoList(pid string, types string, params *schemas.UnitParams) []gjson.Result {
	url := "https://aiqicha.baidu.com/" + types + "?pid=" + pid
	content := GetReq(url, params)
	var listData []gjson.Result
	if gjson.Get(content, "status").String() == "0" {
		data := gjson.Get(content, "data")
		//判断一个获取的特殊值
		if types == "relations/relationalMapAjax" {
			data = gjson.Get(content, "data.investRecordData")
		}
		//判断是否多页，遍历获取所有数据
		pageCount := data.Get("pageCount").Int()
		if pageCount > 1 {
			for i := 1; int(pageCount) >= i; i++ {
				reqUrls := url + "&p=" + strconv.Itoa(i)
				content = GetReq(reqUrls, params)
				listData = append(listData, gjson.Get(content, "data.list").Array()...)
			}
		} else {
			listData = data.Get("list").Array()
		}
	}
	return listData
}

// getCompanyInfoById 获取公司基本信息
// pid 公司id
// isSearch 是否递归搜索信息【分支机构、对外投资信息】
// options options
func getCompanyInfoById(pid string, deep int, inFrom string, isDetail bool, params *schemas.UnitParams, unitInfo *result.UnitInfo) {

	// 获取各个类别需要提取的字段以及相关的 API 数据
	unitCategoryInfoMap := getUnitCategoryInfoMap()

	// 企业基本信息获取
	detailUrl := "https://aiqicha.baidu.com/company_detail_" + pid
	detailResponse := pageParseJson(GetReq(detailUrl, params))
	unitBaseInfo := "unit_info"
	unitJsonTMP, _ := sjson.Set(detailResponse.Raw, "inFrom", inFrom)
	unitInfo.Infos[unitBaseInfo] = append(unitInfo.Infos[unitBaseInfo], gjson.Parse(unitJsonTMP))
	tmpUnitInfos := make(map[string][]gjson.Result)
	// 获取详情数据
	if isDetail {
		unitInfo.Pid = detailResponse.Get("pid").String()
		unitInfo.Name = detailResponse.Get("entName").String()
		unitInfo.LegalPerson = detailResponse.Get("legalPerson").String()
		unitInfo.OpenStatus = detailResponse.Get("openStatus").String()
		unitInfo.Telephone = detailResponse.Get("telephone").String()
		unitInfo.Email = detailResponse.Get("email").String()
		unitInfo.RegCode = detailResponse.Get("taxNo").String()

		// 命令行展示
		data := [][]string{
			{"PID", unitInfo.Pid},
			{"企业名称", unitInfo.Name},
			{"法人代表", unitInfo.LegalPerson},
			{"开业状态", unitInfo.OpenStatus},
			{"电话", unitInfo.Telephone},
			{"邮箱", unitInfo.Email},
			{"统一社会信用代码", unitInfo.RegCode},
		}
		utils.TableShow([]string{}, data)

	}

	// 判断企业状态，如果是异常情况就可以跳过了
	if unitInfo.OpenStatus == "注销" || unitInfo.OpenStatus == "吊销" {

	}

	// 获取企业各个信息列表（基本信息，重点关注，知识产权，企业发展，经营状况，数据解读，新闻资讯）
	infoListUrl := "https://aiqicha.baidu.com/compdata/navigationListAjax?pid=" + pid
	infoListResponse := GetReq(infoListUrl, params)

	// 获取相关信息列表更详细的分类数据
	if gjson.Get(infoListResponse, "status").String() == "0" {
		for _, categoryInfoList := range gjson.Get(infoListResponse, "data").Array() {
			for _, children := range categoryInfoList.Get("children").Array() {
				resId := children.Get("id").String()
				if _, ok := UnitDataTypeMapAQC[resId]; ok {
					resId = UnitDataTypeMapAQC[resId]
				}
				categoryInfo := unitCategoryInfoMap[resId]
				if categoryInfo == nil {
					categoryInfo = &CategoryInfo{}
				}
				categoryInfo.name = children.Get("name").String()
				categoryInfo.total = children.Get("total").Int()
				categoryInfo.available = children.Get("avaliable").Int()
				unitCategoryInfoMap[children.Get("id").String()] = categoryInfo
			}
		}
	}

	// 根据任务参数获取相关数据
	for _, resultField := range params.ResultField {
		if _, ok := unitCategoryInfoMap[resultField]; ok {
			unitCategoryInfo := unitCategoryInfoMap[resultField]
			// 根据信息列表获取到的信息判定那些类别是有数据的（total>0说明该类别有数据，存在 API 说明能获取到该类别数据）
			if unitCategoryInfo.total > 0 && unitCategoryInfo.api != "" {
				// 投资信息、股东信息、供应商、分支信息、控股企业
				if (resultField == global.SearchInvest || resultField == global.SearchPartner || resultField ==
					global.SearchSupplier || resultField == global.SearchBranch || resultField == global.SearchHolds) && (deep > params.Deep) {
					continue
				}
				categoryInfoList := getCategoryInfoList(pid, unitCategoryInfo.api, params)
				//判断下网站备案，然后提取出来，处理下数据
				if resultField == global.SearchICP {
					var tmp []gjson.Result
					for _, categoryInfo := range categoryInfoList {
						for _, domain := range categoryInfo.Get("domain").Array() {
							valueTmp, _ := sjson.Set(categoryInfo.Raw, "domain", domain.String())
							valueTmp, _ = sjson.Set(valueTmp, "homeSite", categoryInfo.Get("homeSite").Array()[0].String())
							tmp = append(tmp, gjson.Parse(valueTmp))
						}
					}
					categoryInfoList = tmp
				}

				// 添加来源信息，并把信息存储到数据里面
				for _, categoryInfo := range categoryInfoList {
					valueTmp, _ := sjson.Set(categoryInfo.Raw, "inFrom", inFrom)
					unitInfo.Infos[resultField] = append(unitInfo.Infos[resultField], gjson.Parse(valueTmp))
					//存入临时数据
					tmpUnitInfos[resultField] = append(tmpUnitInfos[resultField], gjson.Parse(valueTmp))
				}

				//命令输出展示
				var data [][]string
				for _, categoryInfo := range categoryInfoList {
					results := gjson.GetMany(categoryInfo.Raw, unitCategoryInfoMap[resultField].field...)
					var str []string
					for _, ss := range results {
						str = append(str, ss.String())
					}
					data = append(data, str)
				}
				utils.TableShow(unitCategoryInfoMap[resultField].keyWord, data)
			}
		}
	}
	//判断是否查询层级信息 deep
	if deep <= params.Deep {
		// 查询对外投资详细信息
		// 对外投资>0 && 是否递归 && 参数投资比例大于0
		if unitCategoryInfoMap["invest"].total > 0 && params.InvestNum > 0 {
			for _, tmp := range tmpUnitInfos["invest"] {
				openStatus := tmp.Get("openStatus").String()
				if openStatus == "注销" || openStatus == "吊销" {
					continue
				}
				logger.Info(fmt.Sprintf("企业名称：%s 投资【%d级】占比：%s", tmp.Get("entName"), deep, tmp.Get("regRate")))
				// 计算投资比例信息
				investNum := utils.FormatInvest(tmp.Get("regRate").String())
				// 如果达到设定要求就开始获取信息
				if investNum >= params.InvestNum {
					beReason := fmt.Sprintf("%s 投资【%d级】占比 %s", tmp.Get("entName"), deep, tmp.Get("regRate"))
					getCompanyInfoById(tmp.Get("pid").String(), deep+1, beReason, false, params, unitInfo)
				}
			}
		}

		// 查询分支机构公司详细信息
		// 分支机构大于0 && 是否递归模式 && 参数是否开启查询
		// 不查询分支机构的分支机构信息
		if unitCategoryInfoMap["branch"].total > 0 && params.IsBranch {
			for _, tmp := range tmpUnitInfos["branch"] {
				if tmp.Get("inFrom").String() == "" {
					openStatus := tmp.Get("openStatus").String()
					if openStatus == "注销" || openStatus == "吊销" {
						continue
					}
					logger.Info(fmt.Sprintf("分支名称：%s 状态：%s", tmp.Get("entName"), tmp.Get("openStatus")))
					beReason := fmt.Sprintf("%s 分支机构", tmp.Get("entName"))
					getCompanyInfoById(tmp.Get("pid").String(), deep+1, beReason, false, params, unitInfo)
				}
			}
		}

		//查询控股公司
		// 不查询下层信息
		if unitCategoryInfoMap["holds"].total > 0 && params.IsHold {
			if len(tmpUnitInfos["holds"]) == 0 {
				logger.Info("【无控股信息】，需要账号开通【超级会员】！")
			} else {
				for _, tmp := range tmpUnitInfos["holds"] {
					if tmp.Get("inFrom").String() == "" {
						openStatus := tmp.Get("openStatus").String()
						if openStatus == "注销" || openStatus == "吊销" {
							continue
						}
						logger.Info(fmt.Sprintf("控股公司：%s 状态：%s", tmp.Get("entName"), tmp.Get("openStatus")))
						beReason := fmt.Sprintf("%s 控股公司投资比例 %s", tmp.Get("entName"), tmp.Get("proportion"))
						getCompanyInfoById(tmp.Get("pid").String(), deep+1, beReason, false, params, unitInfo)
					}
				}
			}
		}

		// 查询供应商
		// 不查询下层信息
		if unitCategoryInfoMap["supplier"].total > 0 && params.IsSupplier {
			for _, tmp := range tmpUnitInfos["supplier"] {
				if tmp.Get("inFrom").String() == "" {
					openStatus := tmp.Get("openStatus").String()
					if openStatus == "注销" || openStatus == "吊销" {
						continue
					}
					logger.Info(fmt.Sprintf("供应商：%s 状态：%s", tmp.Get("supplier"), tmp.Get("openStatus")))
					beReason := fmt.Sprintf("%s 供应商", tmp.Get("supplier"))
					getCompanyInfoById(tmp.Get("supplierId").String(), deep+1, beReason, false, params, unitInfo)
				}
			}
		}
	}

}

// GetUnitInfoByPid 根据PID获取公司信息
func GetUnitInfoByPid(params *schemas.UnitParams) (*result.UnitInfo, map[string]*result.OrgMap) {
	pid := ""
	if params.CompanyID == "" {
		_, pid = SearchName(params)
	} else {
		pid = params.CompanyID
	}
	// 获取公司信息
	unitInfo := &result.UnitInfo{}
	outMap := make(map[string]*result.OrgMap)
	unitInfo.ScanSource = global.SourceAqc

	if pid == "" {
		logger.Warn("没有获取到PID")
		return unitInfo, outMap
	}
	logger.Info(fmt.Sprintf("查询PID %s", pid))

	unitInfo.Infos = make(map[string][]gjson.Result)
	getCompanyInfoById(pid, 1, "", true, params, unitInfo)
	params.CompanyName = unitInfo.Name

	for k, v := range getUnitCategoryInfoMap() {
		outMap[k] = &result.OrgMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}

	return unitInfo, outMap

}

// pageParseJson 提取页面中的JSON字段
func pageParseJson(content string) gjson.Result {

	tag1 := "window.pageData ="
	tag2 := "window.isSpider ="
	//tag2 := "/* eslint-enable */</script><script data-app"
	idx1 := strings.Index(content, tag1)
	idx2 := strings.Index(content, tag2)
	if idx2 > idx1 {
		str := content[idx1+len(tag1) : idx2]
		str = strings.Replace(str, "", "", -1)
		str = strings.Replace(str, " ", "", -1)
		str = str[:len(str)-1]
		return gjson.Get(str, "result")
	} else {
		logger.Warn("无法解析信息错误信息%s", zap.String("content", content))
	}
	return gjson.Result{}
}

// SearchName 根据企业名称搜索信息
func SearchName(params *schemas.UnitParams) ([]gjson.Result, string) {
	name := params.KeyWord

	urls := "https://aiqicha.baidu.com/s?q=" + urlTool.QueryEscape(name) + "&t=0"
	content := GetReq(urls, params)
	rq := pageParseJson(content)
	unitList := rq.Get("resultList").Array()
	if len(unitList) == 0 {
		logger.Warn(fmt.Sprintf("没有查询到关键词 “%s” ", name))
		return unitList, ""
	} else {
		logger.Info(fmt.Sprintf("关键词：“%s” 查询到 %d 个结果，默认选择第一个 ", name, len(unitList)))
	}

	return unitList, unitList[0].Get("pid").String()
}
