package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	//EmoticonQuest 颜文字指令关键词
	EmoticonQuest = "颜文字"
	//LoveWordsQuest 情话指令关键词
	LoveWordsQuest = "lovewords"
)

type turingInputText struct {
	Text string `json:"text"`
}

type turingPerception struct {
	InputText turingInputText `json:"inputText"`
}

type turingUserInfo struct {
	//APIKey 图灵123 api key
	APIKey string `json:"apiKey"`

	//UserID 用户标识符
	UserID string `json:"userId"`

	//GroupID 群聊标识符
	GroupID string `json:"groupId"`

	//UserIDName 群内昵称
	UserIDName string `json:"userIdName"`
}

type turingAPIIntent struct {
	Code int16 `json:"code"`

	IntentName string `json:"intentName"`

	ActionName string `json:"actionName"`
}

type turingAPIResult struct {
	ResultType string `json:"resultType"`

	Values turingInputText `json:"values"`

	GroupType int16 `json:"groupType"`
}

//TuringAPIRequest 图灵api请求结构体
type TuringAPIRequest struct {
	ReqType int16 `json:"reqType"`

	Perception *turingPerception `json:"perception"`

	UserInfo *turingUserInfo `json:"userInfo"`
}

//TuringAPIResponse 图灵api返回结构体
type TuringAPIResponse struct {
	Intent *turingAPIIntent `json:"intent"`

	Results []turingAPIResult `json:"results"`
}

func randomEmoticon() string {
	happy := []string{"_(┐「ε:)_", "_(:3 」∠)_", "(￣y▽￣)~*捂嘴偷笑",
		"・゜・(PД`q｡)・゜・", "(ง •̀_•́)ง", "(•̀ᴗ•́)و ̑̑", "ヽ(•̀ω•́ )ゝ", "(,,• ₃ •,,)",
		"(｡˘•ε•˘｡)", "(=ﾟωﾟ)ﾉ", "(○’ω’○)", "(´・ω・`)", "ヽ(･ω･｡)ﾉ", "(。-`ω´-)",
		"(´・ω・`)", "(´・ω・)ﾉ", "(ﾉ･ω･)", "  (♥ó㉨ò)ﾉ♡", "(ó㉨ò)", "・㉨・", "( ・◇・)？",
		"ヽ(*´Д｀*)ﾉ", "(´°̥̥̥̥̥̥̥̥ω°̥̥̥̥̥̥̥̥｀)", "(╭￣3￣)╭♡", "(☆ﾟ∀ﾟ)", "⁄(⁄ ⁄•⁄ω⁄•⁄ ⁄)⁄.",
		"(´-ι_-｀)", "ಠ౪ಠ", "ಥ_ಥ", "(/≥▽≤/)", "ヾ(o◕∀◕)ﾉ ヾ(o◕∀◕)ﾉ ヾ(o◕∀◕)ﾉ", "*★,°*:.☆\\(￣▽￣)/$:*.°★*",
		"ヾ (o ° ω ° O ) ノ゙", "╰(*°▽°*)╯ ", " (｡◕ˇ∀ˇ◕）", "o(*≧▽≦)ツ", "≖‿≖✧", ">ㅂ<", "ˋ▽ˊ", "\\(•ㅂ•)/♥",
		"✪ε✪", "✪υ✪", "✪ω✪", "눈_눈", ",,Ծ‸Ծ,,", "π__π", "（/TДT)/", "ʅ（´◔౪◔）ʃ", "(｡☉౪ ⊙｡)", "o(*≧▽≦)ツ┏━┓拍桌狂笑",
		" (●'◡'●)ﾉ♥", "<(▰˘◡˘▰)>", "｡◕‿◕｡", "(｡・`ω´･)", "(♥◠‿◠)ﾉ", "ʅ(‾◡◝) ", " (≖ ‿ ≖)✧", "（´∀｀*)",
		"（＾∀＾）", "(o^∇^o)ﾉ", "ヾ(=^▽^=)ノ", "(*￣∇￣*)", " (*´∇｀*)", "(*ﾟ▽ﾟ*)", "(｡･ω･)ﾉﾞ", "(≡ω≡．)",
		"(｀･ω･´)", "(´･ω･｀)", "(●´ω｀●)φ)"}
	rand.Seed(time.Now().Unix())
	return happy[rand.Intn(len(happy))]
}

func randomLoveWord() string {
	words := []string{"我唯一害怕的就是失去你，我感觉你快要消失了",
		"我养你吧。",
		"我能遇见你已经不可思议了",
		"我在未来等你",
		"我闭着眼睛也看不见自己，但是我却可以看见你",
		"You are like everything to me now",
		"这世界上，除了你我别无所求",
		"有时候我会想你想到无法承受",
		"多希望我自己知道怎么放弃你",
		"我已经爱上你一个星期了，记得吗",
		"送花啊，牵手啊，生日礼物这些我都不擅长，但要是说到结婚，我只希望可以娶你",
		"我真是大笨蛋，除了喜欢你什么都不知道",
		"该死，我记不起其他没有你的地方了",
		"眼睛，是你的眼睛，我知道为什么喜欢你了",
		"我喜欢你的一切，特别是你呆呆傻傻的样子",
		"为什么喜欢你，我好像也不知道，但我知道没有别人会比我更喜欢你"}
	rand.Seed(time.Now().Unix())
	return words[rand.Intn(len(words))]
}

func getAnswer(msg, uid, groupID, userIDName, robotName string) (string, error) {
	fmt.Println(msg)
	if strings.Contains(msg, EmoticonQuest) {
		e := randomEmoticon()
		return e, nil
	}
	if msg == LoveWordsQuest {
		e := randomLoveWord()
		return e, nil
	}
	turingConfig, _ := getConfig("turing")
	robotConfig, err := getConfig(robotName)
	if err != nil {
		return "", NewServiceError("robot配置未定义")
	} else if _, exist := robotConfig["key"]; exist == false {
		return "", NewServiceError("robot配置未定义")
	}

	text := turingInputText{
		Text: msg,
	}

	perception := &turingPerception{
		InputText: text,
	}

	userInfo := &turingUserInfo{
		APIKey:     robotConfig["key"],
		UserID:     uid,
		GroupID:    groupID,
		UserIDName: userIDName,
	}

	apiReq := &TuringAPIRequest{
		ReqType:    0,
		Perception: perception,
		UserInfo:   userInfo,
	}

	jsonBody, err := SimpleHTTPPost(turingConfig["base_url"], apiReq)
	if err == nil {
		fmt.Println(string(jsonBody))
		var response *TuringAPIResponse
		err := json.Unmarshal(jsonBody, &response)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		if response.Intent.Code == 10020 {
			fmt.Println(1)
			result := response.Results[0]
			text := result.Values
			return "翻译：" + text.Text, nil
		}
		if response.Intent.Code != 10004 {
			fmt.Println(1)
			return "我累了", nil
		}
		resultLen := len(response.Results)
		if resultLen < 1 {
			fmt.Println(2)
			return "", NewServiceError("")
		}
		result := response.Results[0]
		text := result.Values
		if text.Text == "" {
			fmt.Println(4)
			return "我累了", nil
		}
		fmt.Println(5)
		return text.Text, nil
	}
	return "", err
}
