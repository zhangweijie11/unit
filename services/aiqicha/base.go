package aiqicha

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"net/http"
	"time"
)

type CategoryInfo struct {
	name      string
	total     int64
	available int64
	api       string   // API 地址
	gNum      string   // 判断数量大小的关键词
	field     []string // 获取的字段名称 看JSON
	keyWord   []string // 关键词
}

var UnitDataTypeMapAQC = map[string]string{
	"webRecord":     "icp",
	"appinfo":       "app",
	"wechatoa":      "wechat",
	"enterprisejob": "job",
	"microblog":     "weibo",
	"hold":          "holds",
	"shareholders":  "partner",
}

// 获取单位各类别信息的映射关系
func getUnitCategoryInfoMap() map[string]*CategoryInfo {
	unitCategoryInfoMap := make(map[string]*CategoryInfo)
	unitCategoryInfoMap = map[string]*CategoryInfo{
		"unit_info": {
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
	for k := range unitCategoryInfoMap {
		unitCategoryInfoMap[k].keyWord = append(unitCategoryInfoMap[k].keyWord, "数据关联  ")
		unitCategoryInfoMap[k].field = append(unitCategoryInfoMap[k].field, "inFrom")
	}
	return unitCategoryInfoMap

}

func GetReq(url string, params *schemas.UnitParams) string {
	client := resty.New()
	client.SetTimeout(time.Duration(global.TimeOut) * time.Minute)
	if global.ToolConf.ProxyEnable && global.ToolConf.Proxy != "" {
		client.SetProxy(global.ToolConf.Proxy)
	}

	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36 Edg/98.0.1108.43"},
		"Accept":     {"text/html, application/xhtml+xml, image/jxr, */*"},
		"Cookie":     {params.Cookies.Aiqicha},
		"Referer":    {"https://aiqicha.baidu.com/"},
	}
	resp, err := client.R().Get(url)

	if err != nil {
		if global.ToolConf.Proxy != "" {
			client.RemoveProxy()
		}
		logger.Error(fmt.Sprintf("【AQC】请求发生错误， %s 5秒后重试", url), err)
		time.Sleep(5 * time.Second)
		return GetReq(url, params)
	}
	if resp.StatusCode() == 200 {
		return string(resp.Body())
	} else if resp.StatusCode() == 403 {
		logger.Warn("【AQC】ip被禁止访问网站，请更换ip")
	} else if resp.StatusCode() == 401 {
		logger.Warn("【AQC】Cookie有问题或过期，请重新获取")
	} else if resp.StatusCode() == 302 {
		logger.Warn("【AQC】需要更新Cookie")
	} else if resp.StatusCode() == 404 {
		logger.Warn("【AQC】请求错误 404")
	} else {
		logger.Warn(fmt.Sprintf("【AQC】未知错误 %d", resp.StatusCode()))
	}
	return ""
}
