package main

import (
	"encoding/json"
	"github.com/larspensjo/config"
	"strconv"
	"strings"
	"time"
	"fmt"
)

func GenerateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func JsonEncode(nodes interface{}) string {
	body, err := json.Marshal(nodes)
	if err != nil {
		panic(err.Error())
		return "[]"
	}
	return string(body)
}

func JsonDecode(jsonStr string) interface{} {
	jsonStr = strings.Replace(jsonStr, "\n", "", -1)
	var f interface{}
	err := json.Unmarshal([]byte(jsonStr), &f)
	if err != nil {
		fmt.Println(jsonStr)
		panic(err)
		return false
	}
	return float2Int(f)
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
		return targetConfig, Error("unable to open config file or wrong fomart")
	}
	sections := cfg.Sections()
	if len(sections) == 0 {
		return targetConfig, Error("no " + sec + " config")
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
