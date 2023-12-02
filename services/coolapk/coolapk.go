package coolapk

import (
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/global/utils"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/result"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetReq(params *schemas.UnitParams) (unitInfo *result.UnitInfo, unitOutMap map[string]*result.OrgMap) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("[COOLAPK] ERROR: %v", err))
		}
	}()
	unitInfo = &result.UnitInfo{}
	unitInfo.Infos = make(map[string][]gjson.Result)
	unitOutMap = make(map[string]*result.OrgMap)
	field := []string{"title", "catName", "apkversionname", "lastupdate", "shorttitle", "logo", "apkname", "", "", "inFrom"}
	keyWord := []string{"名称", "分类", "当前版本", "更新时间", "简介", "logo", "Bundle ID", "链接", "market", "数据关联"}
	unitOutMap["app"] = &result.OrgMap{Name: "app", Field: field, KeyWord: keyWord}
	developer := params.KeyWord
	deviceId := "34de7eef-8400-3300-8922-a1a34e7b9b4f"
	ctime := time.Now().Unix()
	md5Timestamp := utils.Md5(strconv.FormatInt(ctime, 10))
	arg1 := "token://com.coolapk.market/c67ef5943784d09750dcfbb31020f0ab?" + md5Timestamp + "$" + deviceId + "&com.coolapk.market"
	md5Str := utils.Md5(base64.StdEncoding.EncodeToString([]byte(arg1)))
	token := md5Str + deviceId + "0x" + strconv.FormatInt(ctime, 16)
	fmt.Println(token)
	url := "https://api.coolapk.com/v6/apk/search?searchType=developer&developer=" +
		developer +
		"&page=1&firstLaunch=0&installTime=" +
		strconv.FormatInt(ctime, 10) +
		"&lastItem=13988"
	client := resty.New()
	client.SetTimeout(global.TimeOut)
	client.Header = http.Header{
		"X-App-Token":      {token},
		"X-App-Version":    {"10.5.3"},
		"User-Agent":       {"Dalvik/2.1.0 (Linux; U; Android 6.0.1; Nexus 6P Build/MMB29M) (#Build; google; Nexus 6P; MMB29M; 6.0.1) +CoolMarket/10.5.3-2009271"},
		"X-Api-Version":    {"10"},
		"X-App-Device":     {"QZDIzVHel5EI7UGbn92bnByOpV2dhVHSgszQyoTMzoDM2oTQCpDMwoDNyAyOsxWduByO2ADO4kjNxIDM2gjN3YDOgsDZiBTYykzYkZDNlBzY0ITZ"},
		"X-Dark-Mode":      {"0"},
		"X-Requested-With": {"XMLHttpRequest"},
		"X-App-Code":       {"2009271"},
		"X-App-Id":         {"com.coolapk.market"},
	}

	resp, err := client.R().Get(url)

	if err != nil {
		logger.Error("coolapk 请求发生错误\n %s", err)
	}

	appList := gjson.Get(string(resp.Body()), "data").Array()
	unitInfo.Infos["app"] = appList
	unitInfo.Name = params.KeyWord
	logger.Info(fmt.Sprintf("酷安API 查询到 %d 条数据", len(appList)))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(keyWord)
	for _, v := range appList {
		res := gjson.GetMany(v.Raw, field...)
		var str []string
		for k, vv := range res {
			if field[k] == "lastupdate" {
				str = append(str, time.Unix(vv.Int(), 0).Format("2006-01-02 15:04:05"))
			} else {
				str = append(str, vv.String())
			}

		}
		table.Append(str)
	}
	table.Render()
	return unitInfo, unitOutMap
}
