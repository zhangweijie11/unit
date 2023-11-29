package global

import (
	"github.com/mitchellh/mapstructure"
	toolGlobal "gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/config"
)

const (
	TimeOut         = 3
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

var ToolConf config.ToolConfig
var DefaultAllSource = []string{"qcc", "aqc", "tyc", "xlb", "aldzs", "coolapk", "qimai", "chinaz", "all"}
var DefaultAllInfos = []string{"unit_info", "icp", "weibo", "wechat", "app", "wx_app", "job", "copyright"}
var CanSearchAllInfos = []string{"unit_info", "icp", "weibo", "wechat", "app", "wx_app", "job", "copyright", "supplier", "invest", "branch", "holds", "partner"}
var ScanTypeKeys = map[string]string{
	"aqc":     "爱企查",
	"qcc":     "企查查",
	"tyc":     "天眼查",
	"xlb":     "小蓝本",
	"all":     "全部查询",
	"aldzs":   "阿拉丁",
	"coolapk": "酷安市场",
	"qimai":   "七麦数据",
	"chinaz":  "站长之家",
}
var SourceTypeMap = map[string]string{
	"爱企查":  "aqc",
	"企查查":  "qcc",
	"天眼查":  "tyc",
	"小蓝本":  "xlb",
	"阿拉丁":  "aldzs",
	"酷安市场": "coolapk",
	"七麦数据": "qimai",
	"站长之家": "chinaz",
}

func InitToolConf() {
	toolConfig := toolGlobal.Config.Tool
	err := mapstructure.Decode(toolConfig, &ToolConf)
	if err != nil {
		logger.Error("读取工具配置出现错误！", err)
	}
}
