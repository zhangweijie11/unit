package result

import (
	"github.com/tidwall/gjson"
	"time"
)

type OrgMap struct {
	Name    string
	Field   []string
	JField  []string
	KeyWord []string
	Only    string
}

var UnitInfoList = make(map[string][][]interface{})
var OrgMapList = make(map[string]*OrgMap)
var OrgMapLN = map[string]*OrgMap{
	"unit_info": {
		Name:    "企业信息",
		JField:  []string{"name", "legal_person", "status", "phone", "email", "registered_capital", "incorporation_date", "address", "scope", "reg_code", "pid"},
		KeyWord: []string{"企业名称", "法人代表", "经营状态", "电话", "邮箱", "注册资本", "成立日期", "注册地址", "经营范围", "统一社会信用代码", "PID"},
	},
	"icp": {
		Name:    "ICP信息",
		Only:    "domain",
		JField:  []string{"wesbite_name", "website", "domain", "icp", "company_name"},
		KeyWord: []string{"网站名称", "网址", "域名", "网站备案/许可证号", "公司名称"},
	},
	"wx_app": {
		Name:    "微信小程序",
		JField:  []string{"name", "category", "logo", "qrcode", "read_num"},
		KeyWord: []string{"名称", "分类", "头像", "二维码", "阅读量"},
	},
	"wechat": {
		Name:    "微信公众号",
		JField:  []string{"name", "wechat_id", "description", "qrcode", "avatar"},
		KeyWord: []string{"名称", "ID", "简介", "二维码", "头像"},
	},
	"weibo": {
		Name:    "微博",
		JField:  []string{"name", "profile_url", "description", "avatar"},
		KeyWord: []string{"微博昵称", "链接", "简介", "头像"},
	},
	"supplier": {
		Name:    "供应商",
		JField:  []string{"name", "scale", "amount", "report_time", "data_source", "relation", "pid"},
		KeyWord: []string{"名称", "金额占比", "金额", "报告期/公开时间", "数据来源", "关联关系", "PID"},
	},
	"job": {
		Name:    "招聘",
		JField:  []string{"name", "education", "location", "publish_time", "salary"},
		KeyWord: []string{"招聘职位", "学历", "办公地点", "发布日期", "薪资"},
	},
	"invest": {
		Name:    "投资",
		JField:  []string{"name", "legal_person", "status", "scale", "pid"},
		KeyWord: []string{"企业名称", "法人", "状态", "投资比例", "PID"},
	},
	"branch": {
		Name:    "分支机构",
		JField:  []string{"name", "legal_person", "status", "pid"},
		KeyWord: []string{"企业名称", "法人", "状态", "PID"},
	},
	"holds": {
		Name:    "控股企业",
		JField:  []string{"name", "legal_person", "status", "scale", "level", "pid"},
		KeyWord: []string{"企业名称", "法人", "状态", "投资比例", "持股层级", "PID"},
	},
	"app": {
		Name:    "应用",
		JField:  []string{"name", "category", "version", "update_at", "description", "logo", "bundle_id", "link", "market"},
		KeyWord: []string{"名称", "分类", "当前版本", "更新时间", "简介", "logo", "Bundle ID", "链接", "market"},
	},
	"copyright": {
		Name:    "软件著作权",
		JField:  []string{"name", "short_name", "category", "reg_num", "pub_type"},
		KeyWord: []string{"软件全称", "软件简称", "分类", "登记号", "权利取得方式"},
	},
	"partner": {
		Name:    "股东信息",
		JField:  []string{"name", "scale", "reg_cap", "pid"},
		KeyWord: []string{"股东名称", "持股比例", "认缴出资金额", "PID"},
	},
}

type UnitInfo struct {
	Name        string                              // 单位名称
	Pid         string                              // 唯一标识
	LegalPerson string                              // 法人代表
	OpenStatus  string                              // 开业状态
	Email       string                              // 邮箱
	Telephone   string                              // 电话
	ScanSource  string                              // 扫描源
	RegCode     string                              // 统一社会信用代码
	BranchNum   int64                               // 分支结构数量
	InvestNum   int64                               // 投资机构数量
	InTime      time.Time                           //
	PidS        map[string]string                   //
	Infos       map[string][]gjson.Result           // 网页信息
	EnInfos     map[string][]map[string]interface{} //
	EnInfo      []map[string]interface{}            //
}
