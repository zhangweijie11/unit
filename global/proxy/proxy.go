package proxy

import (
	"fmt"
	"io"
	"net/http"
)

func GetProxy() (string, error) {

	// api链接
	api_url := ""

	// 请求api链接
	req, _ := http.NewRequest("GET", api_url, nil)
	client := &http.Client{}
	res, err := client.Do(req)

	// 处理返回结果
	if err != nil {
		// 请求发生异常
		fmt.Println(err.Error())
		return "", err
	} else {
		defer res.Body.Close() //保证最后关闭Body
		// 无gzip压缩, 读取返回内容
		body, _ := io.ReadAll(res.Body)
		return string(body), err
	}
}
