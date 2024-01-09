package general

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

/*
Curl Get 请求
url：请求地址
header：请求头
body：返回体
err：错误信息
*/
func CurlGet(url string, header map[string]string) (body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	// 开始请求
	if req, err = http.NewRequest("GET", url, bytes.NewReader([]byte{})); err != nil {
		return
	}
	// 设置请求头
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	// 发送请求
	if resp, err = httpClient.Do(req); err != nil {
		return
	}
	// 关闭请求
	defer resp.Body.Close()
	// 获取响应体
	body, err = io.ReadAll(resp.Body)
	return
}

/*
Curl Post Json请求
url：请求地址
data：请求参数
header：请求头
body：返回体
err：错误信息
*/
func CurlPostJson(uri string, data interface{}, header map[string]string) (body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	bytesData, _ := json.Marshal(data)
	// 开始请求
	if req, err = http.NewRequest("POST", uri, bytes.NewReader(bytesData)); err != nil {
		return
	}
	// 设置请求头
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// 发送请求
	if resp, err = httpClient.Do(req); err != nil {
		return
	}
	// 关闭请求
	defer resp.Body.Close()
	// 获取响应体
	body, err = io.ReadAll(resp.Body)
	return
}

/*
Curl Post Form请求
url：请求地址
data：请求参数
header：请求头
body：返回体
err：错误信息
*/
func CurlPostForm(uri string, data map[string]string, header map[string]string) (body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	values := url.Values{}
	for key, value := range data {
		values.Add(key, value)
	}
	// 开始请求
	if req, err = http.NewRequest("POST", uri, strings.NewReader(values.Encode())); err != nil {
		return
	}
	// 设置请求头
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 发送请求
	if resp, err = httpClient.Do(req); err != nil {
		return
	}
	// 关闭请求
	defer resp.Body.Close()
	// 获取响应体
	body, err = io.ReadAll(resp.Body)
	return
}
