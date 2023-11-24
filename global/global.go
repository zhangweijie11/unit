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

var ToolConf config.ToolConfig

// DefaultAllInfos 默认收集信息列表
var DefaultAllInfos = []string{"icp", "weibo", "wechat", "app", "weibo", "job", "wx_app", "copyright"}
var DefaultInfos = []string{"icp", "weibo", "wechat", "app", "wx_app"}
var CanSearchAllInfos = []string{"enterprise_info", "icp", "weibo", "wechat", "app", "weibo", "job", "wx_app", "copyright", "supplier", "invest", "branch", "holds", "partner"}

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

func init() {
	toolConfig := toolGlobal.Config.Tool
	err := mapstructure.Decode(toolConfig, &ToolConf)
	if err != nil {
		logger.Error("读取工具配置出现错误！", err)
	}
}
