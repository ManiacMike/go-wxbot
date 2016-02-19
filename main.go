package main

import (
	"fmt"
	"net/http"
  "log"
	"net/url"
	"io/ioutil"
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


func (this *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  msg := r.PostFormValue("msg")

	uid := r.PostFormValue("uid")

	if uid == ""{
		uid = "myself"
	}

	switch this.ApiName {
	case "message":
		var reply string
		fmt.Println(msg)
		turingConfig, _ := getConfig("turing")
		resp, err := http.PostForm(turingConfig["base_url"], url.Values{"uid":{uid}, "key": {turingConfig["key"]}, "info": {msg}})
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
