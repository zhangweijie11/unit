package tianyancha

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	toolGlobal "gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/result"
	"golang.org/x/net/html"
	"strconv"
	"strings"
	"time"
)

/* Tianyancha By Gungnir,Keac
 * admin@wgpsec.org
 */

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
	unitInfo.ScanSource = global.SourceTyc

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

// getCompanyInfoById 获取公司基本信息
// pid 公司id
// isSearch 是否递归搜索信息【分支机构、对外投资信息】
// options options
func getCompanyInfoById(pid string, deep int, inFrom string, isDetail bool, params *schemas.UnitParams, unitInfo *result.UnitInfo) {
	// 获取各个类别需要提取的字段以及相关的 API 数据
	unitCategoryInfoMap := getUnitCategoryInfoMap()

	var unitCategoryCount gjson.Result
	var detailResponse gjson.Result
	tmpUnitInfos := make(map[string][]gjson.Result)
	//切换获取企业信息和统计的方式
	tds := false
	//基本信息

	//var res map[string]string
	//res, unitCategoryCount = SearchBaseDetail(pid, ensInfoMap, options)
	//urls := "https://www.tianyancha.com/company/" + pid
	//body := GetReqReturnPage(urls, options)
	//提取页面的JS数据
	detailResponse, unitCategoryCount = SearchBaseDetail(pid, tds, params)
	unitJsonTMP, _ := sjson.Set(detailResponse.Raw, "inFrom", inFrom)
	//修复成立日期信息
	ts := time.UnixMilli(detailResponse.Get("fromTime").Int())
	unitJsonTMP, _ = sjson.Set(unitJsonTMP, "fromTime", ts.Format(toolGlobal.TimeFormatDay))
	unitBaseInfo := "unit_info"
	unitInfo.Infos[unitBaseInfo] = append(unitInfo.Infos[unitBaseInfo], gjson.Parse(unitJsonTMP))
	//数量统计 API base_count
	//unitCategoryCount = enBaseInfo.Get("props.pageProps.dehydratedState.queries").Array()[16].Get("state.data")
	if isDetail {
		unitInfo.Pid = detailResponse.Get("name").String()
		unitInfo.Name = detailResponse.Get("name").String()
		unitInfo.LegalPerson = detailResponse.Get("legalPersonName").String()
		unitInfo.OpenStatus = detailResponse.Get("regStatus").String()
		unitInfo.Telephone = detailResponse.Get("phoneNumber").String()
		unitInfo.Email = detailResponse.Get("email").String()
		unitInfo.RegCode = detailResponse.Get("creditCode").String()
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

	//获取数据
	for _, resultField := range params.ResultField {
		if _, ok := unitCategoryInfoMap[resultField]; ok {
			// 投资信息、股东信息、供应商、分支信息、控股企业
			if (resultField == global.SearchInvest || resultField == global.SearchPartner || resultField ==
				global.SearchSupplier || resultField == global.SearchBranch || resultField == global.SearchHolds) && (deep > params.Deep) {
				continue
			}
			// 获取企业详情时已经顺便把各个类别的数量一并获取，通过 gNum 可以直接提取
			if unitCategoryCount.Get(unitCategoryInfoMap[resultField].gNum).Int() > 0 {
				unitCategoryInfo := unitCategoryInfoMap[resultField]
				categoryInfoList := getCategoryInfoList(pid, unitCategoryInfo.api, unitCategoryInfo, params)
				for _, categoryInfo := range categoryInfoList {
					valueTmp, _ := sjson.Set(categoryInfo.Raw, "inFrom", inFrom)
					unitInfo.Infos[resultField] = append(unitInfo.Infos[resultField], gjson.Parse(valueTmp))
					//存入临时数据
					tmpUnitInfos[resultField] = append(tmpUnitInfos[resultField], gjson.Parse(valueTmp))
				}

				if len(categoryInfoList) > 0 {
					//命令输出展示
					var data [][]string
					for _, categoryInfo := range categoryInfoList {
						results := gjson.GetMany(categoryInfo.Raw, unitCategoryInfoMap[resultField].field...)
						var str []string
						for _, s := range results {
							str = append(str, s.String())
						}
						data = append(data, str)
					}
					utils.TableShow(unitCategoryInfoMap[resultField].keyWord, data)
				}
			}

		}
	}

	//判断是否查询层级信息 deep
	if deep <= params.Deep {
		// 查询分支机构公司详细信息
		// 分支机构大于0 && 是否递归模式 && 参数是否开启查询
		if params.InvestNum > 0 {
			for _, tmp := range tmpUnitInfos["invest"] {
				openStatus := tmp.Get("regStatus").String()
				if openStatus == "注销" || openStatus == "吊销" {
					continue
				}
				logger.Info(fmt.Sprintf("企业名称：%s 投资占比：%s", tmp.Get("name"), tmp.Get("percent")))
				// 计算投资比例信息
				investNum := utils.FormatInvest(tmp.Get("percent").String())
				// 如果达到设定要求就开始获取信息
				if investNum >= params.InvestNum {
					beReason := fmt.Sprintf("%s 投资【%d级】占比 %s", tmp.Get("name"), deep, tmp.Get("percent"))
					getCompanyInfoById(tmp.Get("id").String(), deep+1, beReason, false, params, unitInfo)
				}
			}

		}
		// 查询分支机构公司详细信息
		// 分支机构大于0 && 是否递归模式 && 参数是否开启查询
		// 不查询分支机构的分支机构信息
		if params.IsBranch {
			for _, tmp := range tmpUnitInfos["branch"] {
				if tmp.Get("inFrom").String() == "" {
					openStatus := tmp.Get("regStatus").String()
					if openStatus == "注销" || openStatus == "吊销" {
						continue
					}
					logger.Info(fmt.Sprintf("分支名称：%s 状态：%s\n", tmp.Get("name"), tmp.Get("regStatus")))
					beReason := fmt.Sprintf("%s 分支机构", tmp.Get("entName"))
					getCompanyInfoById(tmp.Get("id").String(), deep+1, beReason, false, params, unitInfo)
				}
			}
		}

		//查询控股公司
		// 不查询下层信息
		if params.IsHold {
			if len(tmpUnitInfos["holds"]) == 0 {
				logger.Info("需要登陆才能查询控股公司！")
			} else {
				for _, tmp := range tmpUnitInfos["holds"] {
					if tmp.Get("inFrom").String() == "" {
						openStatus := tmp.Get("regStatus").String()
						if openStatus == "注销" || openStatus == "吊销" {
							continue
						}
						logger.Info(fmt.Sprintf("控股公司：%s 状态：%s", tmp.Get("name"), tmp.Get("regStatus")))
						beReason := fmt.Sprintf("%s 控股公司投资比例 %s", tmp.Get("name"), tmp.Get("percent"))
						getCompanyInfoById(tmp.Get("cid").String(), deep+1, beReason, false, params, unitInfo)
					}
				}
			}
		}
		// 查询供应商
		// 不查询下层信息
		if params.IsSupplier {
			for _, tmp := range tmpUnitInfos["supplier"] {
				if tmp.Get("inFrom").String() == "" {
					openStatus := tmp.Get("regStatus").String()
					if openStatus == "注销" || openStatus == "吊销" {
						continue
					}
					logger.Info(fmt.Sprintf("供应商：%s 状态：%s", tmp.Get("supplier_name"), tmp.Get("regStatus")))
					beReason := fmt.Sprintf("%s 供应商", tmp.Get("supplier_name"))
					getCompanyInfoById(tmp.Get("supplier_graphId").String(), deep+1, beReason, false, params, unitInfo)
				}
			}
		}

	}
}

func pageParseJson(content string) (res gjson.Result) {
	content = strings.ReplaceAll(content, "var aa = ", "")
	return gjson.Parse(content)
}

func SearchBaseDetail(pid string, tds bool, params *schemas.UnitParams) (result gjson.Result, unitBaseInfo gjson.Result) {
	url := "https://www.tianyancha.com/company/" + pid

	if tds {
		//htmlInfo := htmlquery.FindOne(body, "//*[@class=\"position-rel company-header-container\"]//script")
		//unitBaseInfo = pageParseJson(htmlquery.InnerText(htmlInfo))
		result = gjson.Get(GetReq("https://capi.tianyancha.com/cloud-other-information/companyinfo/baseinfo/web?id="+pid, "", params), "data")
		fmt.Println(result.String())
	} else {
		body := GetReqReturnPage(url, params)
		htmlInfos := htmlquery.FindOne(body, "//*[@id=\"__NEXT_DATA__\"]")
		unitInfo := gjson.Parse(htmlquery.InnerText(htmlInfos))
		unitInfoData := unitInfo.Get("props.pageProps.dehydratedState.queries").Array()
		result = unitInfoData[0].Get("state.data.data")
		//数量统计 API base_count
		for i := 0; i < len(unitInfoData); i++ {
			if unitInfoData[i].Get("queryKey").String() == "base_count" {
				unitBaseInfo = unitInfoData[i].Get("state.data")
			}
		}
		//unitBaseInfo = unitInfo.Get("props.pageProps.dehydratedState.queries").Array()[11].Get("state.data")
	}

	return result, unitBaseInfo
}

func SearchBaseInfoByTables(pid string, ensMap map[string]*CategoryInfo, params *schemas.UnitParams) (result map[string]string, unitInfoCount gjson.Result) {
	//var re = regexp.MustCompile(`(?m)placeholder="请输入公司名称、老板姓名、品牌名称等\"\s*value="(.*?)\"\/>`)
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("[TYC] SearchBaseDetail panic: %v", err))
		}
	}()
	url := "https://www.tianyancha.com/company/" + pid
	body := GetReqReturnPage(url, params)
	htmlInfo := htmlquery.FindOne(body, "//*[@class=\"position-rel company-header-container\"]//script")
	unitInfoCount = pageParseJson(htmlquery.InnerText(htmlInfo))
	htmlAll := htmlquery.Find(body, "//*[@id=\"_container_baseInfo\"]/table")
	result = make(map[string]string)
	isOrg := len(htmlquery.Find(htmlAll[0], "//tr"))
	var orgPs [][]int
	if isOrg == 5 {
		//兼容事业单位 社会组织
		orgPs = [][]int{{0}, {1, 2}, {1, 6}, {1}, {2}, {1, 4}, {3, 6}, {4, 2}, {5, 2}, {3, 4}, {}}
	} else {
		orgPs = ensMap["unit_info"].PosiToTaeS
	}

	for k, v := range orgPs {
		esf := ensMap["unit_info"].field[k]
		if esf == "pid" {
			result[esf] = pid
		}
		if len(v) == 1 {
			res := htmlquery.Find(body, "//*[@data-clipboard-target=\"#copyCompanyInfoThroughThisTag\"]")
			if len(res) > v[0] {
				result[esf] = htmlquery.InnerText(res[v[0]])
			}

		}
		if len(v) == 2 {
			expr := fmt.Sprintf("//tr[%d]/td[%d]", v[0], v[1])
			if esf == "legalPerson" && isOrg != 5 {
				expr += "//a"
			}
			htmlResTmp := htmlquery.Find(htmlAll[0], expr)
			if len(htmlResTmp) > 0 {
				result[esf] = htmlquery.InnerText(htmlResTmp[0])
			}
		}
	}

	return result, unitInfoCount
}

func SearchName(params *schemas.UnitParams) ([]gjson.Result, string) {
	name := params.KeyWord
	//使用关键词推荐方法进行检索，会出现信息不对的情况
	//urls := "https://sp0.tianyancha.com/search/suggestV3?_=" + url.QueryEscape(name)
	urls := "https://capi.tianyancha.com/cloud-tempest/web/searchCompanyV3"
	searchData := map[string]string{
		"key":      name,
		"pageNum":  "1",
		"pageSize": "20",
		"referer":  "search",
		"sortType": "0",
		"word":     name,
	}
	marshal, err := json.Marshal(searchData)
	if err != nil {
		return nil, ""
	}
	content := GetReq(urls, string(marshal), params)
	unitList := gjson.Get(content, "data.companyList").Array()

	if len(unitList) == 0 {
		logger.Info(fmt.Sprintf("没有查询到关键词 “%s” ", name))
		return unitList, ""
	} else {
		logger.Info(fmt.Sprintf("关键词：“%s” 查询到 %d 个结果，默认选择第一个", name, len(unitList)))
	}
	return unitList, unitList[0].Get("id").String()
}

func JudgePageNumWithCookie(page *html.Node) int {
	list := htmlquery.Find(page, "//li")
	return len(list) - 1
}

func getCategoryInfoList(pid string, types string, s *CategoryInfo, params *schemas.UnitParams) (listData []gjson.Result) {
	data := ""
	if len(s.sData) != 0 {
		dataTmp, _ := json.Marshal(s.sData)
		data = string(dataTmp)
	}
	url := "https://capi.tianyancha.com/" + types + "?_=" + strconv.Itoa(int(time.Now().Unix()))

	if data == "" {
		url += "&pageSize=100&graphId=" + pid + "&id=" + pid + "&gid=" + pid + "&pageNum=1" + s.gsData
	} else {
		data, _ = sjson.Set(data, "gid", pid)
		data, _ = sjson.Set(data, "pageSize", 100)
		data, _ = sjson.Set(data, "pageNum", 1)
	}
	content := GetReq(url, data, params)
	if gjson.Get(content, "state").String() != "ok" {
		return listData
	}
	pageCount := 0
	pList := []string{"itemTotal", "count", "total", "pageBean.total"}
	for _, k := range gjson.GetMany(gjson.Get(content, "data").Raw, pList...) {
		if k.Int() != 0 {
			pageCount = int(k.Int())
		}
	}
	pats := "data." + s.rf

	listData = gjson.Get(content, pats).Array()
	if pageCount > 100 {
		url = strings.ReplaceAll(url, "&pageNum=1", "")
		for i := 2; pageCount/100 >= i-1; i++ {
			reqUrls := url
			if data == "" {
				reqUrls = url + "&pageNum=" + strconv.Itoa(i)
			} else {
				data, _ = sjson.Set(data, "pageNum", i)
			}
			content = GetReq(reqUrls, data, params)
			listData = append(listData, gjson.Get(content, pats).Array()...)
		}
	}

	return listData
}

func getInfoListByTable(pid string, ensInfoMap *CategoryInfo, params *schemas.UnitParams) []map[string]string {
	urls := "https://www.tianyancha.com/" + ensInfoMap.api + "?ps=30&id=" + pid
	page := GetReqReturnPage(urls, params)

	List := getTb(page, ensInfoMap, 1)
	page_num := JudgePageNumWithCookie(page)
	if page_num > 1 {
		for i := 2; i <= page_num; i++ {
			urls = "https://www.tianyancha.com/" + ensInfoMap.api + "?ps=30&id=" + pid + "&pn=" + strconv.Itoa(i)
			page = GetReqReturnPage(urls, params)
			tmp_List := getTb(page, ensInfoMap, i)
			List = append(List, tmp_List...)
		}
	}
	return List

}

func SearchByName(params *schemas.UnitParams) (enName string) {
	res, _ := SearchName(params)
	if len(res) > 0 {
		enName = res[0].Get("comName").String()
	}
	return enName
}

func getTb(page *html.Node, categoryInfo *CategoryInfo, page_num int) []map[string]string {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("[TYC] getTb panic: %v", err))
		}
	}()
	var infoss []map[string]string
	exrps := "//body/table/tbody/tr"
	htmlAll, err := htmlquery.QueryAll(page, exrps)
	flag := false
	if len(htmlAll) == 0 {
		flag = true
		exrps = "//body/div/table/tbody/tr"
		htmlAll, err = htmlquery.QueryAll(page, exrps)
	}
	//doc := goquery.NewDocumentFromNode(page)
	if err != nil {
		panic(`not a valid XPath expression.`)
	}
	//doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
	//	s.Find("td").Each(func(ii int, t *goquery.Selection) {
	//		fmt.Println(t.Text())
	//	})
	//})

	for i := 0; i < len(htmlAll); i++ {
		results := make(map[string]string)
		for tmpNum := 0; tmpNum < len(categoryInfo.field)-1; tmpNum++ {
			if (categoryInfo.field[tmpNum] == "") || categoryInfo.PosiToTake[tmpNum] == 0 {
				continue
			}
			expr := "//td"
			if flag {
				expr = "/td"
			}
			htmls, _ := htmlquery.QueryAll(htmlAll[i], expr)
			//fmt.Println(htmlquery.InnerText(htmlAll[i]))
			//if len(htmls) == 2 {
			//	fmt.Println(htmlquery.InnerText(htmls[0]))
			//}

			taskNum := categoryInfo.PosiToTake[tmpNum] - 1
			if len(htmls) < taskNum {
				continue
			}
			htmlAllS := htmls[taskNum]

			if categoryInfo.field[tmpNum] == "logo" || categoryInfo.field[tmpNum] == "qrcode" {
				htmlA := htmlquery.Find(htmlAllS, "//img")
				if len(htmlA) > 0 {
					results[categoryInfo.field[tmpNum]] = htmlquery.SelectAttr(htmlA[0], "data-src")
				}
			} else if categoryInfo.field[tmpNum] == "pid" || categoryInfo.field[tmpNum] == "StockName" || categoryInfo.field[tmpNum] == "legalPerson" || categoryInfo.field[tmpNum] == "href" {
				htmlA := htmlquery.Find(htmlAllS, "//a")
				if len(htmlA) > 0 {
					if categoryInfo.field[tmpNum] == "pid" || categoryInfo.field[tmpNum] == "href" {
						if categoryInfo.name == "股东信息" && strings.Contains(htmlquery.SelectAttr(htmlA[0], "href"), "human") {
							results[categoryInfo.field[tmpNum]] = ""
							//results[categoryInfo.field[tmpNum]] = strings.ReplaceAll(htmlquery.SelectAttr(htmlA[0], "href"), " https://www.tianyancha.com/human/", "")

						} else {
							results[categoryInfo.field[tmpNum]] = strings.ReplaceAll(htmlquery.SelectAttr(htmlA[0], "href"), "https://www.tianyancha.com/company/", "")
						}

					} else {
						results[categoryInfo.field[tmpNum]] = htmlquery.InnerText(htmlA[0])
					}
				}
			} else if categoryInfo.name == "供应商" && (categoryInfo.field[tmpNum] == "entName" || categoryInfo.field[tmpNum] == "source") {

				htmlA := htmlquery.Find(htmlAllS, "//a")
				if categoryInfo.field[tmpNum] == "entName" {
					if len(htmlA) > 0 {
						results[categoryInfo.field[tmpNum]] = htmlquery.InnerText(htmlA[0])
					} else {
						str := htmlquery.InnerText(htmlAllS)
						str = strings.ReplaceAll(str, "查看全部", "")
						str = strings.ReplaceAll(str, "条采购数据", "")
						results[categoryInfo.field[tmpNum]] = str
					}
				}
				if categoryInfo.field[tmpNum] == "source" {
					results[categoryInfo.field[tmpNum]] = htmlquery.InnerText(htmlAllS)
					if len(htmlA) > 0 {
						results[categoryInfo.field[tmpNum]] += "https://www.tianyancha.com/" + htmlquery.SelectAttr(htmlA[0], "href")
					}
				}

			} else if categoryInfo.name == "投资信息" && (categoryInfo.field[tmpNum] == "entName") {
				htmlA := htmlquery.Find(htmlAllS, "//a")
				if len(htmlA) > 0 {
					results[categoryInfo.field[tmpNum]] = htmlquery.InnerText(htmlA[0])
				} else {
					results[categoryInfo.field[tmpNum]] = strings.ReplaceAll(htmlquery.InnerText(htmlAllS), "股权结构", "")
				}
			} else if categoryInfo.name == "ICP备案" {
				results[categoryInfo.field[tmpNum]] = strings.ReplaceAll(htmlquery.InnerText(htmlAllS), "该网站与ICP或年报备案网站一致", "")
			} else {
				txt := htmlquery.InnerText(htmlAllS)
				txt = strings.ReplaceAll(txt, "... 更多", "")
				results[categoryInfo.field[tmpNum]] = txt
			}

			//results[categoryInfo.field[tmp_num]] = htmlquery.InnerText(a[tmp_num+i*categoryInfo.NumOfEachGroup])
		}
		infoss = append(infoss, results)
	}
	return infoss
}
