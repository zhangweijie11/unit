package services

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/result"
	"go.uber.org/zap"
	urlTool "net/url"
	"os"
	"strconv"
	"strings"
)

type EnsGo struct {
	name      string
	total     int64
	available int64
	api       string   //API 地址
	gNum      string   //判断数量大小的关键词
	field     []string //获取的字段名称 看JSON
	keyWord   []string //关键词
}

var ENSMapAQC = map[string]string{
	"webRecord":     "icp",
	"appinfo":       "app",
	"wechatoa":      "wechat",
	"enterprisejob": "job",
	"microblog":     "weibo",
	"hold":          "holds",
	"shareholders":  "partner",
}

func getENMap() map[string]*EnsGo {
	ensInfoMap := make(map[string]*EnsGo)
	ensInfoMap = map[string]*EnsGo{
		"enterprise_info": {
			name:    "企业信息",
			field:   []string{"entName", "legalPerson", "openStatus", "telephone", "email", "regCapital", "startDate", "regAddr", "scope", "taxNo", "pid"},
			keyWord: []string{"企业名称", "法人代表", "经营状态", "电话", "邮箱", "注册资本", "成立日期", "注册地址", "经营范围", "统一社会信用代码", "PID"},
		},
		"icp": {
			name:    "ICP备案",
			api:     "detail/icpinfoAjax",
			field:   []string{"siteName", "homeSite", "domain", "icpNo", ""},
			keyWord: []string{"网站名称", "网址", "域名", "网站备案/许可证号", "公司名称"},
		},
		"app": {
			name:    "APP",
			api:     "c/appinfoAjax",
			field:   []string{"name", "classify", "", "", "logoBrief", "logo", "", "", ""},
			keyWord: []string{"名称", "分类", "当前版本", "更新时间", "简介", "logo", "Bundle ID", "链接", "market"},
		},
		"weibo": {
			name:    "微博",
			api:     "c/microblogAjax",
			field:   []string{"nickname", "weiboLink", "brief", "logo"},
			keyWord: []string{"微博昵称", "链接", "简介", "LOGO"},
		},
		"wechat": {
			name:    "微信公众号",
			api:     "c/wechatoaAjax",
			field:   []string{"wechatName", "wechatId", "wechatIntruduction", "qrcode", "wechatLogo"},
			keyWord: []string{"名称", "ID", "描述", "二维码", "LOGO"},
		},
		"job": {
			name:    "招聘信息",
			api:     "c/enterprisejobAjax",
			field:   []string{"jobTitle", "education", "location", "publishDate", "desc"},
			keyWord: []string{"招聘职位", "学历要求", "工作地点", "发布日期", "招聘描述"},
		},
		"copyright": {
			name:    "软件著作权",
			api:     "detail/copyrightAjax",
			field:   []string{"softwareName", "shortName", "softwareType", "PubType", ""},
			keyWord: []string{"软件名称", "软件简介", "分类", "登记号", "权利取得方式"},
		},
		"supplier": {
			name:    "供应商",
			api:     "c/supplierAjax",
			field:   []string{"supplier", "", "", "cooperationDate", "source", "", "supplierId"},
			keyWord: []string{"名称", "金额占比", "金额", "报告期/公开时间", "数据来源", "关联关系", "PID"},
		},
		"invest": {
			name:    "投资信息",
			api:     "detail/investajax",
			field:   []string{"entName", "legalPerson", "openStatus", "regRate", "pid"},
			keyWord: []string{"企业名称", "法人", "状态", "投资比例", "PID"},
		},
		"holds": {
			name:    "控股企业",
			api:     "detail/holdsAjax",
			field:   []string{"entName", "", "", "proportion", "", "pid"},
			keyWord: []string{"企业名称", "法人", "状态", "投资比例", "持股层级", "PID"},
		},
		"branch": {
			name:    "分支信息",
			api:     "detail/branchajax",
			field:   []string{"entName", "legalPerson", "openStatus", "pid"},
			keyWord: []string{"企业名称", "法人", "状态", "PID"},
		},
		"partner": {
			name:    "股东信息",
			api:     "detail/sharesAjax",
			field:   []string{"name", "subRate", "subMoney", "pid"},
			keyWord: []string{"股东名称", "持股比例", "认缴出资金额", "PID"},
		},
	}
	for k := range ensInfoMap {
		ensInfoMap[k].keyWord = append(ensInfoMap[k].keyWord, "数据关联  ")
		ensInfoMap[k].field = append(ensInfoMap[k].field, "inFrom")
	}
	return ensInfoMap

}

// getInfoList 获取信息列表
func getInfoList(pid string, types string) []gjson.Result {
	urls := "https://aiqicha.baidu.com/" + types + "?pid=" + pid
	content := utils.GetReq(urls)
	var listData []gjson.Result
	if gjson.Get(string(content), "status").String() == "0" {
		data := gjson.Get(string(content), "data")
		//判断一个获取的特殊值
		if types == "relations/relationalMapAjax" {
			data = gjson.Get(string(content), "data.investRecordData")
		}
		//判断是否多页，遍历获取所有数据
		pageCount := data.Get("pageCount").Int()
		if pageCount > 1 {
			for i := 1; int(pageCount) >= i; i++ {
				logger.Info(fmt.Sprintf("当前：%s,%d\n", types, i))
				reqUrls := urls + "&p=" + strconv.Itoa(i)
				content = utils.GetReq(reqUrls)
				listData = append(listData, gjson.Get(string(content), "data.list").Array()...)
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
func getCompanyInfoById(pid string, deep int, inFrom string, params *schemas.UnitParams, unitInfo *result.UnitInfo) {

	// 获取初始化API数据
	ensInfoMap := getENMap()
	// 企业基本信息获取

	urls := "https://aiqicha.baidu.com/company_detail_" + pid
	res := pageParseJson(utils.GetReq(urls))
	//获取企业基本信息情况
	enDes := "enterprise_info"
	enJsonTMP, _ := sjson.Set(res.Raw, "inFrom", inFrom)
	unitInfo.Infos[enDes] = append(unitInfo.Infos[enDes], gjson.Parse(enJsonTMP))
	tmpEIS := make(map[string][]gjson.Result)
	if params.IsEnDetail {
		unitInfo.Pid = res.Get("pid").String()
		unitInfo.Name = res.Get("entName").String()
		unitInfo.LegalPerson = res.Get("legalPerson").String()
		unitInfo.OpenStatus = res.Get("openStatus").String()
		unitInfo.Telephone = res.Get("telephone").String()
		unitInfo.Email = res.Get("email").String()
		unitInfo.RegCode = res.Get("taxNo").String()

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

	// 获取企业信息列表
	enInfoUrl := "https://aiqicha.baidu.com/compdata/navigationListAjax?pid=" + pid
	enInfoRes := utils.GetReq(enInfoUrl)

	// 初始化数量数据
	if gjson.Get(enInfoRes, "status").String() == "0" {
		for _, s := range gjson.Get(enInfoRes, "data").Array() {
			for _, t := range s.Get("children").Array() {
				resId := t.Get("id").String()
				if _, ok := ENSMapAQC[resId]; ok {
					resId = ENSMapAQC[resId]
				}
				es := ensInfoMap[resId]
				if es == nil {
					es = &EnsGo{}
				}
				//fmt.Println(t.Get("name").String() + "|" + t.Get("id").String())
				es.name = t.Get("name").String()
				es.total = t.Get("total").Int()
				es.available = t.Get("avaliable").Int()
				ensInfoMap[t.Get("id").String()] = es
			}
		}
	}

	//获取数据
	for _, k := range params.SearchField {
		if _, ok := ensInfoMap[k]; ok {
			s := ensInfoMap[k]
			if s.total > 0 && s.api != "" {
				if k == "branch" && !params.IsGetBranch {
					continue
				}
				if (k == "invest" || k == "partner" || k == "supplier" || k == "branch" || k == "holds") && (deep > params.Deep) {
					continue
				}
				t := getInfoList(pid, s.api)
				//判断下网站备案，然后提取出来，处理下数据
				if k == "icp" {
					var tmp []gjson.Result
					for _, y := range t {
						for _, o := range y.Get("domain").Array() {
							valueTmp, _ := sjson.Set(y.Raw, "domain", o.String())
							valueTmp, _ = sjson.Set(valueTmp, "homeSite", y.Get("homeSite").Array()[0].String())
							tmp = append(tmp, gjson.Parse(valueTmp))
						}
					}
					t = tmp
				}

				// 添加来源信息，并把信息存储到数据里面
				for _, y := range t {
					valueTmp, _ := sjson.Set(y.Raw, "inFrom", inFrom)
					unitInfo.Infos[k] = append(unitInfo.Infos[k], gjson.Parse(valueTmp))
					//存入临时数据
					tmpEIS[k] = append(tmpEIS[k], gjson.Parse(valueTmp))
				}

				//命令输出展示
				var data [][]string
				for _, y := range t {
					results := gjson.GetMany(y.Raw, ensInfoMap[k].field...)
					var str []string
					for _, ss := range results {
						str = append(str, ss.String())
					}
					data = append(data, str)
				}
				utils.TableShow(ensInfoMap[k].keyWord, data)
			}
		}
	}
	//判断是否查询层级信息 deep
	if deep <= params.Deep {
		// 查询对外投资详细信息
		// 对外投资>0 && 是否递归 && 参数投资信息大于0
		if ensInfoMap["invest"].total > 0 && params.InvestNum > 0 {
			for _, t := range tmpEIS["invest"] {
				logger.Info(fmt.Sprintf("企业名称：%s 投资【%d级】占比：%s\n", t.Get("entName"), deep, t.Get("regRate")))
				openStatus := t.Get("openStatus").String()
				if openStatus == "注销" || openStatus == "吊销" {
					continue
				}
				// 计算投资比例信息
				investNum := utils.FormatInvest(t.Get("regRate").String())
				// 如果达到设定要求就开始获取信息
				if investNum >= params.InvestNum {
					beReason := fmt.Sprintf("%s 投资【%d级】占比 %s", t.Get("entName"), deep, t.Get("regRate"))
					getCompanyInfoById(t.Get("pid").String(), deep+1, beReason, params, unitInfo)
				}
			}
		}

		// 查询分支机构公司详细信息
		// 分支机构大于0 && 是否递归模式 && 参数是否开启查询
		// 不查询分支机构的分支机构信息
		if ensInfoMap["branch"].total > 0 && params.IsGetBranch && params.IsSearchBranch {
			for _, t := range tmpEIS["branch"] {
				if t.Get("inFrom").String() == "" {
					openStatus := t.Get("openStatus").String()
					if openStatus == "注销" || openStatus == "吊销" {
						continue
					}
					logger.Info(fmt.Sprintf("分支名称：%s 状态：%s\n", t.Get("entName"), t.Get("openStatus")))
					beReason := fmt.Sprintf("%s 分支机构", t.Get("entName"))
					getCompanyInfoById(t.Get("pid").String(), -1, beReason, params, unitInfo)
				}
			}
		}

		//查询控股公司
		// 不查询下层信息
		if ensInfoMap["holds"].total > 0 && params.IsHold {
			if len(tmpEIS["holds"]) == 0 {
				logger.Info(fmt.Sprintf("【无控股信息】，需要账号开通【超级会员】！\n"))
			} else {
				for _, t := range tmpEIS["holds"] {
					if t.Get("inFrom").String() == "" {
						openStatus := t.Get("openStatus").String()
						logger.Info(fmt.Sprintf("控股公司：%s 状态：%s\n", t.Get("entName"), t.Get("openStatus")))
						if openStatus == "注销" || openStatus == "吊销" {
							continue
						}
						beReason := fmt.Sprintf("%s 控股公司投资比例 %s", t.Get("entName"), t.Get("proportion"))
						getCompanyInfoById(t.Get("pid").String(), -1, beReason, params, unitInfo)
					}
				}
			}
		}

		// 查询供应商
		// 不查询下层信息
		if ensInfoMap["supplier"].total > 0 && params.IsSupplier {
			for _, t := range tmpEIS["supplier"] {
				if t.Get("inFrom").String() == "" {
					openStatus := t.Get("openStatus").String()
					logger.Info(fmt.Sprintf("供应商：%s 状态：%s\n", t.Get("supplier"), t.Get("openStatus")))
					if openStatus == "注销" || openStatus == "吊销" {
						continue
					}
					beReason := fmt.Sprintf("%s 供应商", t.Get("supplier"))
					getCompanyInfoById(t.Get("supplierId").String(), -1, beReason, params, unitInfo)
				}
			}
		}
	}

}

// GetEnInfoByPid 根据PID获取公司信息
func GetEnInfoByPid(params *schemas.UnitParams) (*result.UnitInfo, map[string]*result.OrgMap) {
	pid := ""
	if params.CompanyID == "" {
		_, pid = SearchName(params)
	} else {
		pid = params.CompanyID
	}
	//获取公司信息
	unitInfo := &result.UnitInfo{}
	outMap := make(map[string]*result.OrgMap)

	if params.PID == "" {
		logger.Warn("没有获取到PID\n")
		return unitInfo, outMap
	}
	logger.Info(fmt.Sprintf("查询PID %s\n", pid))

	unitInfo.Infos = make(map[string][]gjson.Result)
	getCompanyInfoById(pid, 1, "", params, unitInfo)
	params.CompanyName = unitInfo.Name

	for k, v := range getENMap() {
		outMap[k] = &result.OrgMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}

	//outputfile.OutPutExcelByEnInfo(ensInfos, ensOutMap, options)
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
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, " ", "", -1)
		str = str[:len(str)-1]
		return gjson.Get(str, "result")
	} else {
		logger.Warn("无法解析信息错误信息%s\n", zap.String("content", content))
	}
	return gjson.Result{}
}

// SearchName 根据企业名称搜索信息
func SearchName(params *schemas.UnitParams) ([]gjson.Result, string) {
	name := params.KeyWord

	urls := "https://aiqicha.baidu.com/s?q=" + urlTool.QueryEscape(name) + "&t=0"
	content := utils.GetReq(urls)
	rq := pageParseJson(content)
	enList := rq.Get("resultList").Array()
	if len(enList) == 0 {
		logger.Warn(fmt.Sprintf("没有查询到关键词 “%s” \n", name))
		return enList, ""
	} else {
		logger.Info(fmt.Sprintf("关键词：“%s” 查询到 %d 个结果，默认选择第一个 \n", name, len(enList)))
	}
	if global.ToolConf.IsTableShow {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"PID", "企业名称", "法人代表", "社会统一信用代码"})
		for _, v := range enList {
			table.Append([]string{v.Get("pid").String(), v.Get("titleName").String(), v.Get("titleLegal").String(), v.Get("regNo").String()})
		}
		table.Render()
	}
	return enList, enList[0].Get("pid").String()
}
