package main

import (
	"encoding/json"
	"github.com/bitly/go-simplejson" // for json get
	"github.com/larspensjo/config"
	"strconv"
	"time"
)

func JsonStrToMap(jsonStr string) map[string]interface{} {
	json, err := simplejson.NewJson([]byte(jsonStr))
	if err != nil {
		panic(err.Error())
	}
	var nodes = make(map[string]interface{})
	nodes, _ = json.Map()
	return nodes
}

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
	var f interface{}
	err := json.Unmarshal([]byte(jsonStr), &f)
	if err != nil {
		panic(err)
		return false
	}
	return float2Int(f)
}

func float2Int(input interface{}) interface{} {
	switch input.(type) {
	case []interface{}:
		for k, v := range input.([]interface{}) {
			switch v.(type) {
			case float64:
				input.([]interface{})[k] = int(v.(float64))
			case []interface{}:
				input.([]interface{})[k] = float2Int(input.([]interface{})[k])
			case map[string]interface{}:
				input.([]interface{})[k] = float2Int(input.([]interface{})[k])
			}
		}
	case map[string]interface{}:
		for k, v := range input.(map[string]interface{}) {
			switch v.(type) {
			case float64:
				input.(map[string]interface{})[k] = int(v.(float64))
			case []interface{}:
				input.(map[string]interface{})[k] = float2Int(input.(map[string]interface{})[k])
			case map[string]interface{}:
				input.(map[string]interface{})[k] = float2Int(input.(map[string]interface{})[k])
			}
		}
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
