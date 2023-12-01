package utils

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"os"
	"strconv"
	"strings"
)

func TableShow(keys []string, values [][]string) {
	if global.ToolConf.IsTableShow {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetHeader(keys)
		table.AppendBulk(values)
		table.Render()
	}
}

func FormatInvest(scale string) float64 {
	if scale == "-" || scale == "" || scale == " " {
		return -1
	} else {
		scale = strings.ReplaceAll(scale, "%", "")
	}

	num, err := strconv.ParseFloat(scale, 64)
	if err != nil {
		return -1
	}
	return num
}

func IsInList(target string, list []string) bool {
	if len(list) == 0 {
		return false
	}
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func DelInList(target string, list []string) []string {
	var result []string
	for _, v := range list {
		if v != target {
			result = append(result, v)
		}
	}
	return result
}

// CheckPid 检查pid是哪家单位
func CheckPid(pid string) (res string) {
	if len(pid) == 32 {
		res = global.SourceQcc
	} else if len(pid) == 14 {
		res = global.SourceAqc
	} else if len(pid) == 8 || len(pid) == 7 || len(pid) == 6 || len(pid) == 9 || len(pid) == 10 {
		res = global.SourceTyc
	} else if len(pid) == 33 || len(pid) == 34 {
		if pid[0] == 'p' {
			logger.Warn("无法查询法人信息")
			res = ""
		} else {
			res = global.SourceXlb
		}
	} else {
		logger.Warn(fmt.Sprintf("pid长度%d不正确，pid: %s", len(pid), pid))
		return ""
	}
	return res
}
