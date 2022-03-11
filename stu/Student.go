package stu

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type Student struct {
	Userid    string
	ReqHeader map[string]string
}

func NewStudent() *Student {
	countryCapitalMap := make(map[string]string)
	countryCapitalMap["Host"] = "wljy.whut.edu.cn"
	countryCapitalMap["Connection"] = "keep-alive"
	countryCapitalMap["Accept"] = "application/json, text/javascript, */*; q=0.01"
	countryCapitalMap["X-Requested-With"] = "XMLHttpRequest"
	countryCapitalMap["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
	countryCapitalMap["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	countryCapitalMap["Origin"] = "http://wljy.whut.edu.cn"
	countryCapitalMap["Cookie"] = "sspm_proctorUrl=http://wljy.whut.edu.cn:80/; sspm_orgid=4406; sspm_appid=157438568781088; JSESSIONID=7F1F8C548B6C0086EB550B8F945EA999"
	return &Student{
		Userid:    "201693412300013",
		ReqHeader: countryCapitalMap,
	}
}

func (s *Student) SendHttp(par string, method string, url string) io.ReadCloser {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(par))
	defer req.Body.Close() // 函数结束时关闭Body
	if err != nil {
		log.Fatal(err)
		return nil
	}
	header := s.ReqHeader
	//循环设置
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Body
}
