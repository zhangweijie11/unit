package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetReq(url string) string {
	client := resty.New()
	client.SetTimeout(time.Duration(global.TimeOut) * time.Minute)
	if global.ToolConf.Proxy != "" {
		client.SetProxy(global.ToolConf.Proxy)
	}

	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36 Edg/98.0.1108.43"},
		"Accept":     {"text/html, application/xhtml+xml, image/jxr, */*"},
		"Cookie":     {global.ToolConf.Cookies.Aiqicha},
		"Referer":    {"https://aiqicha.baidu.com/"},
	}
	resp, err := client.R().Get(url)

	if err != nil {
		if global.ToolConf.Proxy != "" {
			client.RemoveProxy()
		}
		logger.Error(fmt.Sprintf("【AQC】请求发生错误， %s 5秒后重试\n", url), err)
		time.Sleep(5 * time.Second)
		return GetReq(url)
	}
	if resp.StatusCode() == 200 {
		return string(resp.Body())
	} else if resp.StatusCode() == 403 {
		logger.Warn("【AQC】ip被禁止访问网站，请更换ip\n")
	} else if resp.StatusCode() == 401 {
		logger.Warn("【AQC】Cookie有问题或过期，请重新获取\n")
	} else if resp.StatusCode() == 302 {
		logger.Warn("【AQC】需要更新Cookie\n")
	} else if resp.StatusCode() == 404 {
		logger.Warn("【AQC】请求错误 404 %s \n")
	} else {
		logger.Warn(fmt.Sprintf("【AQC】未知错误 %d\n", resp.StatusCode()))
	}
	return ""
}

func TableShow(keys []string, values [][]string) {
	if !global.ToolConf.IsApiMode {
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
