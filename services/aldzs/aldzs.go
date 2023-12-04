package aldzs

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/unit/global"
	"gitlab.example.com/zhangweijie/unit/middleware/schemas"
	"gitlab.example.com/zhangweijie/unit/services/result"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getReq(searchType string, data map[string]string) gjson.Result {
	url := fmt.Sprintf("https://zhishuapi.aldwx.com/Main/action/%s", searchType)
	client := resty.New()
	client.SetTimeout(global.TimeOut * time.Second)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.Header = http.Header{
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"},
		"Accept":       {"text/html, application/xhtml+xml, image/jxr, */*"},
		"Content-Type": {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Referer":      {"https://www.aldzs.com"},
	}
	resp, err := client.R().SetFormData(data).Post(url)
	if err != nil {
		fmt.Println(err)
	}
	res := gjson.Parse(string(resp.Body()))
	if res.Get("code").String() != "200" {
		logger.Warn(fmt.Sprintf("【aldzs】似乎出了点问题 %s ", res.Get("msg")))
	}
	return res.Get("data")
}

func GetInfoByKeyword(params *schemas.UnitParams) (unitInfo *result.UnitInfo, unitOutMap map[string]*result.OrgMap) {
	unitInfo = &result.UnitInfo{}
	unitInfo.Infos = make(map[string][]gjson.Result)
	unitOutMap = make(map[string]*result.OrgMap)

	keyword := params.KeyWord
	//拿到Token信息
	//token := params.Cookies.Aldzs
	token := "clT85gC3kcvFPTZDT6lT3P5REd8IEQtBEfx8TA0pCLs%2BZ%2B7mQOuGUzfTNM1bm8DHUbf0zOg2ji3IAC39nK52hTUVSxJg84zpTVeFHGCRTJIz0fc5sKCXapVAiOYMJEA9OqWnGeuOZqfHMcHrZ4zxeLeUy0zOQyIhLoQWGpwLXY86r2zi6fx%2B62ct3JvsZ%2BVOloQs3glijOmG9S73lFicfSty1SrBA4F9cck1ezxyASA%3D"
	logger.Info(fmt.Sprintf("查询关键词 %s 的小程序", keyword))
	appList := getReq("Search/Search/search", map[string]string{
		"appName":    keyword,
		"page":       "1",
		"token":      token,
		"visit_type": "1",
	}).Array()
	if len(appList) == 0 {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NO", "ID", "小程序名称", "所属公司", "描述"})
	for k, v := range appList {
		table.Append([]string{
			strconv.Itoa(k),
			v.Get("id").String(),
			v.Get("name").String(),
			v.Get("company").String(),
			v.Get("desc").String(),
		})
	}
	table.Render()
	//默认取第一个进行查询
	logger.Info(fmt.Sprintf("查询 %s 开发的相关小程序 【默认取100个】", appList[0].Get("company")))
	appKey := appList[0].Get("appKey").String()
	sAppList := getReq("Miniapp/App/sameBodyAppList", map[string]string{
		"appKey": appKey,
		"page":   "1",
		"size":   "100",
		"token":  token,
	}).Array()
	unitInfo.Infos["wx_app"] = sAppList
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NO", "ID", "小程序名称", "描述"})
	for k, v := range sAppList {
		table.Append([]string{
			strconv.Itoa(k),
			v.Get("id").String(),
			v.Get("name").String(),
			v.Get("desc").String(),
		})
	}
	table.Render()

	for k, v := range getUnitCategoryInfoMap() {
		unitOutMap[k] = &result.OrgMap{Name: v.name, Field: v.field, KeyWord: v.keyWord}
	}
	return unitInfo, unitOutMap
}

type CategoryInfo struct {
	name     string
	api      string
	fids     string
	params   map[string]string
	field    []string
	keyWord  []string
	typeInfo []string
}

func getUnitCategoryInfoMap() map[string]*CategoryInfo {
	unitCategoryInfoMap := make(map[string]*CategoryInfo)
	unitCategoryInfoMap = map[string]*CategoryInfo{
		"wx_app": {
			name:    "微信小程序",
			field:   []string{"name", "categoryTitle", "logo", "", ""},
			keyWord: []string{"名称", "分类", "头像", "二维码", "阅读量"},
		},
	}
	for k := range unitCategoryInfoMap {
		unitCategoryInfoMap[k].keyWord = append(unitCategoryInfoMap[k].keyWord, "数据关联  ")
		unitCategoryInfoMap[k].field = append(unitCategoryInfoMap[k].field, "inFrom")
	}
	return unitCategoryInfoMap
}
