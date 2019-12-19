package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	EMOTICON_QUEST  = "颜文字"
	LOVEWORDS_QUEST = "lovewords"
)

type TuringInputText struct {
	Text string `json:"text"`
}

type TuringPerception struct {
	InputText TuringInputText `json:"inputText"`
}

type TuringUserInfo struct {
	ApiKey string `json:"apiKey"`

	UserId string `json:"userId"`

	GroupId string `json:"groupId"`

	UserIdName string `json:"userIdName"`
}

type TuringApiIntent struct {
	Code int16 `json:"code"`

	IntentName string `json:"intentName"`

	ActionName string `json:"actionName"`
}

type TuringApiResult struct {
	ResultType string `json:"resultType"`

	Values []TuringInputText `json:"values"`

	GroupType int16 `json:"groupType"`
}

type TuringApiRequest struct {
	ReqType int16 `json:"reqType"`

	Perception *TuringPerception `json:"perception"`

	UserInfo *TuringUserInfo `json:"userInfo"`
}

type TuringApiResponse struct {
	Intent *TuringApiIntent `json:"intent"`

	Results []TuringApiResult `json:"results"`
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

func getAnswer(msg string, uid string, groupId string, userIdName string, robotName string) (string, error) {
	fmt.Println(msg)
	if strings.Contains(msg, EMOTICON_QUEST) {
		e := randomEmoticon()
		return e, nil
	}
	if msg == LOVEWORDS_QUEST {
		e := randomLoveWord()
		return e, nil
	}
	turingConfig, _ := getConfig("turing")
	robotConfig, err := getConfig(robotName)
	if err != nil {
		return "", Error("robot配置未定义")
	} else if _, exist := robotConfig["key"]; exist == false {
		return "", Error("robot配置未定义")
	}

	text := TuringInputText{
		Text: msg,
	}

	perception := &TuringPerception{
		InputText: text,
	}

	userInfo := &TuringUserInfo{
		ApiKey: robotConfig["key"],
		UserId: uid,
		GroupId: groupId,
		UserIdName: userIdName,
	}

	apiReq := &TuringApiRequest{
		ReqType: 0,
		Perception: perception,
		UserInfo: userInfo,
	}

	jsonBody, err := SimpleHttpPost(turingConfig["base_url"], apiReq)
	if err == nil {
		fmt.Println(string(jsonBody))
		var response *TuringApiResponse
		err := json.Unmarshal(jsonBody, &response)
		if err != nil {
			return "", err
		} else {
			if response.Intent.Code != 10005{
				return "我累了", nil
			}
			resultLen := len(response.Results)
			if resultLen < 1{
				return "", Error("")
			}
			result := response.Results[0]
			if len(result.Values) < 1{
				return "", Error("")
			}
			text := result.Values[0]
			if text.Text != ""{
				return "我累了", nil
			}else{
				return text.Text, nil
			}
		}
	} else {
		return "", err
	}
}
