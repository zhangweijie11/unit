package global

import (
	"github.com/mitchellh/mapstructure"
	toolGlobal "gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/config"
)

const (
	TimeOut = 3
)

// 爬取字段
const (
	SearchUnitInfo  = "unit_info"
	SearchICP       = "icp"
	SearchWeiBo     = "weibo"
	SearchWeChat    = "wechat"
	SearchAPP       = "app"
	SearchWxApp     = "wx_app"
	SearchJob       = "job"
	SearchCopyright = "copyright"
	SearchSupplier  = "supplier"
	SearchInvest    = "invest"
	SearchBranch    = "branch"
	SearchHolds     = "holds"
	SearchPartner   = "partner"
)

// 爬取源
const (
	SourceTyc     = "tyc"
	SourceQcc     = "qcc"
	SourceAqc     = "aqc"
	SourceXlb     = "xlb"
	SourceAldzs   = "aldzs"
	SourceCoolapk = "coolapk"
	SourceQimai   = "qimai"
	SourceChinaz  = "chinaz"
	SourceAll     = "all"
)

var ToolConf config.ToolConfig
var DefaultAllSource = []string{SourceQcc, SourceAqc, SourceTyc, SourceXlb, SourceAldzs, SourceCoolapk, SourceQimai, SourceChinaz, SourceAll}
var DefaultAllInfos = []string{SearchUnitInfo, SearchICP, SearchWeiBo, SearchWeChat, SearchAPP, SearchWxApp, SearchJob, SearchCopyright}
var CanSearchAllInfos = []string{SearchUnitInfo, SearchICP, SearchWeiBo, SearchWeChat, SearchAPP, SearchWxApp, SearchJob, SearchCopyright, SearchSupplier, SearchInvest, SearchBranch, SearchHolds, SearchPartner}
var ScanTypeKeys = map[string]string{
	SourceAqc:     "爱企查",
	SourceQcc:     "企查查",
	SourceTyc:     "天眼查",
	SourceXlb:     "小蓝本",
	SourceAll:     "全部查询",
	SourceAldzs:   "阿拉丁",
	SourceCoolapk: "酷安市场",
	SourceQimai:   "七麦数据",
	SourceChinaz:  "站长之家",
}
var SourceTypeMap = map[string]string{
	"爱企查":  SourceAqc,
	"企查查":  SourceQcc,
	"天眼查":  SourceTyc,
	"小蓝本":  SourceXlb,
	"阿拉丁":  SourceAldzs,
	"酷安市场": SourceCoolapk,
	"七麦数据": SourceQimai,
	"站长之家": SourceChinaz,
}

func InitToolConf() {
	toolConfig := toolGlobal.Config.Tool
	err := mapstructure.Decode(toolConfig, &ToolConf)
	if err != nil {
		logger.Error("读取工具配置出现错误！", err)
	}
}
