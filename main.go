package main

import (
	"fmt"
	"net/http"
  "log"
	"net/url"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
)

type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func Error(msg string) error {
	return &ServiceError{msg}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

type ApiServer struct {
	ApiName string
}

func randomEmoticon() string {
    happy := []string{"_(┐「ε:)_","_(:3 」∠)_","(￣y▽￣)~*捂嘴偷笑",
			"・゜・(PД`q｡)・゜・","(ง •̀_•́)ง","(•̀ᴗ•́)و ̑̑","ヽ(•̀ω•́ )ゝ","(,,• ₃ •,,)",
			"(｡˘•ε•˘｡)","(=ﾟωﾟ)ﾉ","(○’ω’○)","(´・ω・`)","ヽ(･ω･｡)ﾉ","(。-`ω´-)",
			"(´・ω・`)","(´・ω・)ﾉ","(ﾉ･ω･)","  (♥ó㉨ò)ﾉ♡","(ó㉨ò)","・㉨・","( ・◇・)？",
			"ヽ(*´Д｀*)ﾉ","(´°̥̥̥̥̥̥̥̥ω°̥̥̥̥̥̥̥̥｀)","(╭￣3￣)╭♡","(☆ﾟ∀ﾟ)","⁄(⁄ ⁄•⁄ω⁄•⁄ ⁄)⁄.",
			"(´-ι_-｀)","ಠ౪ಠ","ಥ_ಥ","(/≥▽≤/)","ヾ(o◕∀◕)ﾉ ヾ(o◕∀◕)ﾉ ヾ(o◕∀◕)ﾉ","*★,°*:.☆\\(￣▽￣)/$:*.°★*",
			"ヾ (o ° ω ° O ) ノ゙","╰(*°▽°*)╯ "," (｡◕ˇ∀ˇ◕）","o(*≧▽≦)ツ","≖‿≖✧",">ㅂ<","ˋ▽ˊ","\\(•ㅂ•)/♥",
			"✪ε✪","✪υ✪","✪ω✪","눈_눈",",,Ծ‸Ծ,,","π__π","（/TДT)/","ʅ（´◔౪◔）ʃ","(｡☉౪ ⊙｡)","o(*≧▽≦)ツ┏━┓拍桌狂笑",
			" (●'◡'●)ﾉ♥","<(▰˘◡˘▰)>","｡◕‿◕｡","(｡・`ω´･)","(♥◠‿◠)ﾉ","ʅ(‾◡◝) "," (≖ ‿ ≖)✧","（´∀｀*)",
			"（＾∀＾）","(o^∇^o)ﾉ","ヾ(=^▽^=)ノ","(*￣∇￣*)"," (*´∇｀*)","(*ﾟ▽ﾟ*)","(｡･ω･)ﾉﾞ","(≡ω≡．)",
			"(｀･ω･´)","(´･ω･｀)","(●´ω｀●)φ)"}
		rand.Seed(time.Now().Unix())
		return happy[rand.Intn(len(happy))]
}

func (this *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  msg := r.PostFormValue("msg")

	uid := r.PostFormValue("uid")

	robotName := r.PostFormValue("robot")

	if uid == ""{
		uid = "myself"
	}

	if robotName == ""{
		fmt.Println("robot名字未定义")
		return
	}

	switch this.ApiName {
	case "message":
		var reply string
		fmt.Println(msg)
		if strings.Contains(msg, "颜文字"){
			e := randomEmoticon()
			fmt.Fprint(w, e)
			return
		}
		turingConfig, _ := getConfig("turing")
		robotConfig,err := getConfig(robotName)
		if err != nil{
			fmt.Println("robot配置未定义")
			return
		}else if _,exist :=robotConfig["key"]; exist == false{
			fmt.Println("robot配置未定义")
			return
		}
		resp, err := http.PostForm(turingConfig["base_url"], url.Values{"uid":{uid}, "key": {robotConfig["key"]}, "info": {msg}})
		if err != nil{
			fmt.Println("request fail")
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			fmt.Println("request fail")
			return
		}
		fmt.Println(string(body))
		replyMap := JsonDecode(string(body))
		if replyMap["code"] == 100000{
			reply = replyMap["text"].(string)
		}else if replyMap["code"] == 200000{
			reply = "给你个链接接着  " + replyMap["url"].(string)
		}else if replyMap["code"] == 40002{
			reply = "不知道你在说虾米～"
		}else if replyMap["code"] == 40004{
			reply = "今天太累了，明天再聊吧"
		}else{
			reply = "哦"
		}
		fmt.Fprint(w, reply)
	default:
		fmt.Fprint(w, "Invalid api")
	}
}

func main(){
  var err error

  http.Handle("/message", &ApiServer{ApiName: "message"})

	serverConfig, err := getConfig("server")
	if err != nil {
		log.Fatal("server config error:", err)
	}

	fmt.Println("listen on port " + serverConfig["port"])

  if err = http.ListenAndServe(":"+serverConfig["port"], nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
