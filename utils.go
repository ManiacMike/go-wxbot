package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/larspensjo/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//GenerateID unique Id based on unix nano id
func GenerateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func float2Int(input interface{}) interface{} {
	if m, ok := input.([]interface{}); ok {
		for k, v := range m {
			switch v.(type) {
			case float64:
				m[k] = int(v.(float64))
			case []interface{}:
				m[k] = float2Int(m[k])
			case map[string]interface{}:
				m[k] = float2Int(m[k])
			}
		}
	} else if m, ok := input.(map[string]interface{}); ok {
		for k, v := range m {
			switch v.(type) {
			case float64:
				m[k] = int(v.(float64))
			case []interface{}:
				m[k] = float2Int(m[k])
			case map[string]interface{}:
				m[k] = float2Int(m[k])
			}
		}
	} else {
		return false
	}
	return input
}

func getConfig(sec string) (map[string]string, error) {
	targetConfig := make(map[string]string)
	cfg, err := config.ReadDefault("config.ini")
	if err != nil {
		return targetConfig, NewServiceError("unable to open config file or wrong fomart")
	}
	sections := cfg.Sections()
	if len(sections) == 0 {
		return targetConfig, NewServiceError("no " + sec + " config")
	}
	for _, section := range sections {
		if section != sec {
			continue
		}
		sectionData, _ := cfg.SectionOptions(section)
		for _, key := range sectionData {
			value, err := cfg.String(section, key)
			if err == nil {
				targetConfig[key] = value
			}
		}
		break
	}
	return targetConfig, nil
}

//SimpleHTTPPost simple post json func
func SimpleHTTPPost(urlstr string, params interface{}) ([]byte, error) {
	var (
		err  error
		resp *http.Response
	)
	httpclient := http.Client{
		CheckRedirect: nil,
		Jar:           nil,
	}
	jsonPost, err := json.Marshal(params)
	fmt.Println(string(jsonPost))
	if err != nil {
		return []byte(""), NewServiceError("json encode fail")
	}
	requestBody := bytes.NewBuffer([]byte(jsonPost))
	request, err := http.NewRequest("POST", urlstr, requestBody)
	if err != nil {
		return []byte(""), err
	}

	resp, err = httpclient.Do(request)

	if err != nil || resp == nil {
		return []byte(""), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}
