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
	if global.ToolConf.ProxyEnable && global.ToolConf.Proxy != "" {
		client.SetProxy(global.ToolConf.Proxy)
	}

	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36 Edg/98.0.1108.43"},
		"Accept":     {"text/html, application/xhtml+xml, image/jxr, */*"},
		"Cookie":     {"BAIDUID_BFESS=9531255470E63C4EA800BF470ED7BD17:FG=1; __bid_n=18bf4f28fcf80649fa20ea; BDUSS=U1MdjFJQW1KOXRreVUyV1dncEU5ejJMVEZMbn5FVVpUcDN2a3ZtT1h1MVQtb1JsRVFBQUFBJCQAAAAAAAAAAAEAAADXSgA31cXS3TUxNgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNtXWVTbV1lO; BDUSS_BFESS=U1MdjFJQW1KOXRreVUyV1dncEU5ejJMVEZMbn5FVVpUcDN2a3ZtT1h1MVQtb1JsRVFBQUFBJCQAAAAAAAAAAAEAAADXSgA31cXS3TUxNgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNtXWVTbV1lO; ZD_ENTRY=google; BDPPN=e079f5145dc992ecf48316ae35bc14fa; login_type=passport; _t4z_qc8_=xlTM-TogKuTwElJqnFDWa4sxdI37OvsFggmd; Hm_lvt_ad52b306e1ae4557f5d3534cce8f8bbf=1700727979; log_guid=2cd54e95b62c4542ad121bcc73d9cc0e; _j47_ka8_=57; Hm_lpvt_ad52b306e1ae4557f5d3534cce8f8bbf=1700728001; log_first_time=1701054493021; _fb537_=xlTM-TogKuTwIpOtlhlVrAWe9dDGsKTjQYWj5S7snEUxk0pZqDk41Usmd; ab170105400=f41442b98406fa6f80278e3bb0f1bb661701054649235; ab_sr=1.0.1_MzkxM2IxZDFlYWQ2MzQwODc0NTcwMTdmOTczM2YyN2ZlZmRlN2I1NWMzMmQxYzUzZGJmZGYwYzg4OGFkNTA3ZGEzODBkMzllODNmMjQxNWVkNWZhODkwNzM3NTljMWZmMWUwNDRmMmI1ZjhkNDEyMTc3ZTI1ODgzYmEyMTA3MGQzNzk3NjU1YTUwMzU5OGY1NTkwNWNiMmY1NDJlZGQwNQ==; _s53_d91_=261acb38790890e894461f82e14e67eb8ee19bd4a189ea49dde7a1ba28be730537112106ef6913791aab3e4f811b396e00f1edc2479f16dce467c58ff675122d0fd708eae193ce3139e9bc3acd1552e45fac5d24a9bbc209d9c5583eecd1e82432a0502fbde9f18fc291efd5d062dbbcf257add88a7d3eb8bbdc61851b35922a7a912858d0d146cfb77173a470b3c10a38c7e026d6d968269b10232f0ea5a1e5d46fb75aff982224432a965463f2dbb24b4b5be5ffb95ea2cfedb003225f7be11d048a548f8fc1ab20e443cde1a4b3a4; _y18_s21_=555e3262; log_last_time=1701054653914; RT=\"z=1&dm=baidu.com&si=a16e98a5-f96f-474f-ab2c-aebd4a9930bf&ss=lpgc0bm9&sl=2&tt=1g3&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=ge2&ul=ht3\""},
		"Referer":    {"https://aiqicha.baidu.com/"},
	}
	resp, err := client.R().Get(url)

	if err != nil {
		if global.ToolConf.Proxy != "" {
			client.RemoveProxy()
		}
		logger.Error(fmt.Sprintf("【AQC】请求发生错误， %s 5秒后重试", url), err)
		time.Sleep(5 * time.Second)
		return GetReq(url)
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
		logger.Warn("【AQC】请求错误 404 %s ")
	} else {
		logger.Warn(fmt.Sprintf("【AQC】未知错误 %d", resp.StatusCode()))
	}
	return ""
}

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
		res = "qcc"
	} else if len(pid) == 14 {
		res = "aqc"
	} else if len(pid) == 8 || len(pid) == 7 || len(pid) == 6 || len(pid) == 9 || len(pid) == 10 {
		res = "tyc"
	} else if len(pid) == 33 || len(pid) == 34 {
		if pid[0] == 'p' {
			logger.Warn("无法查询法人信息")
			res = ""
		} else {
			res = "xlb"
		}
	} else {
		logger.Warn(fmt.Sprintf("pid长度%d不正确，pid: %s", len(pid), pid))
		return ""
	}
	return res
}
