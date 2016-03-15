//web weixin client
package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	// "log"
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	// "strings"
)

const debug = false

func debugPrint(content interface{}) {
	if debug == true {
		fmt.Println(content)
	}
}

type wxweb struct {
	uuid         string
	base_uri     string
	redirect_uri string
	uin          string
	sid          string
	skey         string
	pass_ticket  string
	deviceId     string
	SyncKey      map[string]interface{}
	synckey      string
	User         map[string]interface{}
	BaseRequest  map[string]interface{}
	syncHost     string
	http_client  *http.Client
}

func (self *wxweb) getUuid(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/jslogin"
	urlstr += "?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + self._unixStr()
	data, _ := self._get(urlstr, false)
	re := regexp.MustCompile(`"([\S]+)"`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		self.uuid = find[1]
		return true
	} else {
		return false
	}
}

func (self *wxweb) _run(desc string, f func(...interface{}) bool, args ...interface{}) {
	start := time.Now().UnixNano()
	fmt.Print(desc)
	var result bool
	if len(args) > 1 {
		result = f(args)
	} else if len(args) == 1 {
		result = f(args[0])
	} else {
		result = f()
	}
	useTime := fmt.Sprintf("%.5f", (float64(time.Now().UnixNano()-start) / 1000000000))
	if result {
		fmt.Println("成功,用时" + useTime + "秒")
	} else {
		fmt.Println("失败\n[*] 退出程序")
		os.Exit(1)
	}
}

func (self *wxweb) _post(urlstr string, params map[string]interface{}, jsonFmt bool) (string, error) {
	var err error
	var resp *http.Response
	if jsonFmt == true {
		jsonPost := JsonEncode(params)
		debugPrint(jsonPost)
		requestBody := bytes.NewBuffer([]byte(jsonPost))
		request, err := http.NewRequest("POST", urlstr, requestBody)
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
		request.Header.Add("Referer", "https://wx.qq.com/")
		request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
		resp, err = self.http_client.Do(request)
		// resp, err = self.http_client.Post(urlstr, "application/json;charset=utf-8", requestBody)
	} else {
		v := url.Values{}
		for key, value := range params {
			v.Add(key, value.(string))
		}
		resp, err = self.http_client.PostForm(urlstr, v)
	}

	if err != nil || resp == nil {
		fmt.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	} else {
		defer resp.Body.Close()
	}
	return string(body), nil
}

func (self *wxweb) _get(urlstr string, jsonFmt bool) (string, error) {
	var err error
	res := ""
	// v := url.Values{}
	// for key, value := range params {
	// 	v.Add(key, value.(string))
	// }
	// urlstr = urlstr + "?" + v.Encode()
	request, _ := http.NewRequest("GET", urlstr, nil)
	request.Header.Add("Referer", "https://wx.qq.com/")
	request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	resp, err := self.http_client.Do(request)
	// resp, err := self.http_client.Get(urlstr)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return string(body), nil
}

func (self *wxweb) _unixStr() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (self *wxweb) genQRcode(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/qrcode/" + self.uuid
	urlstr += "?t=webwx"
	urlstr += "&_=" + self._unixStr()
	path := "qrcode.jpg"
	out, err := os.Create(path)
	resp, err := self._get(urlstr, false)
	_, err = io.Copy(out, bytes.NewReader([]byte(resp)))
	if err != nil {
		return false
	} else {
		if runtime.GOOS == "darwin" {
			exec.Command("open", path).Run()
		} else {
			go func() {
				fmt.Println("please open on web broswer ip:8889/qrcode")
				http.HandleFunc("/qrcode", func(w http.ResponseWriter, req *http.Request) {
					http.ServeFile(w, req, "qrcode.jpg")
					return
				})
				http.ListenAndServe(":8889", nil)
			}()
		}
		return true
	}
}

func (self *wxweb) waitForLogin(tip int) bool {
	time.Sleep(time.Duration(tip) * time.Second)
	url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
	url += "?tip=" + strconv.Itoa(tip) + "&uuid=" + self.uuid + "&_=" + self._unixStr()
	data, _ := self._get(url, false)
	re := regexp.MustCompile(`window.code=(\d+);`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		code := find[1]
		if code == "201" {
			return true
		} else if code == "200" {
			re := regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
			find := re.FindStringSubmatch(data)
			if len(find) > 1 {
				r_uri := find[1] + "&fun=new"
				self.redirect_uri = r_uri
				re = regexp.MustCompile(`/`)
				finded := re.FindAllStringIndex(r_uri, -1)
				self.base_uri = r_uri[:finded[len(finded)-1][0]]
				return true
			}
			return false
		} else if code == "408" {
			fmt.Println("[登陆超时]")
		} else {
			fmt.Println("[登陆异常]")
		}
	}
	return false
}

func (self *wxweb) login(args ...interface{}) bool {
	data, _ := self._get(self.redirect_uri, false)
	type Result struct {
		Skey        string `xml:"skey"`
		Wxsid       string `xml:"wxsid"`
		Wxuin       string `xml:"wxuin"`
		Pass_ticket string `xml:"pass_ticket"`
	}
	v := Result{}
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}
	self.skey = v.Skey
	self.sid = v.Wxsid
	self.uin = v.Wxuin
	self.pass_ticket = v.Pass_ticket
	self.BaseRequest = make(map[string]interface{})
	self.BaseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	self.BaseRequest["Sid"] = v.Wxsid
	self.BaseRequest["Skey"] = v.Skey
	self.BaseRequest["DeviceID"] = self.deviceId
	return true
}

func (self *wxweb) webwxinit(args ...interface{}) bool {
	url := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%s", self.base_uri, self.pass_ticket, self.skey, self._unixStr())
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	res, err := self._post(url, params, true)
	if err != nil {
		return false
	}
	//log
	ioutil.WriteFile("tmp.txt", []byte(res), 777)
	//log

	data := JsonDecode(res).(map[string]interface{})
	self.User = data["User"].(map[string]interface{})
	self.SyncKey = data["SyncKey"].(map[string]interface{})
	self._setsynckey()

	//interface int和int型不能使用==
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	return retCode == 0
}

func (self *wxweb) _setsynckey() {
	keys := []string{}
	for _, keyVal := range self.SyncKey["List"].([]interface{}) {
		key := strconv.Itoa(int(keyVal.(map[string]interface{})["Key"].(int)))
		value := strconv.Itoa(int(keyVal.(map[string]interface{})["Val"].(int)))
		keys = append(keys, key+"_"+value)
	}
	self.synckey = strings.Join(keys, "|")
	debugPrint(self.synckey)
}

func (self *wxweb) synccheck() (string, string) {
	urlstr := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", self.syncHost)
	v := url.Values{}
	v.Add("r", self._unixStr())
	v.Add("sid", self.sid)
	v.Add("uin", self.uin)
	v.Add("skey", self.skey)
	v.Add("deviceid", self.deviceId)
	v.Add("synckey", self.synckey)
	v.Add("_", self._unixStr())
	urlstr = urlstr + "?" + v.Encode()
	data, _ := self._get(urlstr, false)
	re := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	find := re.FindStringSubmatch(data)
	if len(find) > 2 {
		retcode := find[1]
		selector := find[2]
		debugPrint(fmt.Sprintf("retcode:%s,selector,selector%s", find[1], find[2]))
		return retcode, selector
	} else {
		return "9999", "0"
	}
}

func (self *wxweb) testsynccheck(args ...interface{}) bool {
	SyncHost := []string{
		"webpush.weixin.qq.com",
		"webpush2.weixin.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
		//"webpush.wechatapp.com"
	}
	for _, host := range SyncHost {
		self.syncHost = host
		retcode, _ := self.synccheck()
		if retcode == "0" {
			return true
		}
	}
	return false
}

func (self *wxweb) webwxstatusnotify(args ...interface{}) bool {
	urlstr := fmt.Sprintf("%s/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", self.base_uri, self.pass_ticket)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["Code"] = 3
	params["FromUserName"] = self.User["UserName"]
	params["ToUserName"] = self.User["UserName"]
	params["ClientMsgId"] = int(time.Now().Unix())
	res, err := self._post(urlstr, params, true)
	if err != nil {
		return false
	}
	data := JsonDecode(res).(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	return retCode == 0
}

func (self *wxweb) webgetchatroommember(chatroomId string) (map[string]string, error) {
	urlstr := fmt.Sprintf("%s/webwxbatchgetcontact?type=ex&r=%s&pass_ticket=%s", self.base_uri, self._unixStr(), self.pass_ticket)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["Count"] = 1
	params["List"] = []map[string]string{}
	l := []map[string]string{}
	params["List"] = append(l, map[string]string{
		"UserName":   chatroomId,
		"ChatRoomId": "",
	})
	members := []string{}
	stats := make(map[string]string)
	res, err := self._post(urlstr, params, true)
	fmt.Println(urlstr)
	debugPrint(params)
	if err != nil {
		return stats, err
	}
	data := JsonDecode(res).(map[string]interface{})
	RoomContactList := data["ContactList"].([]interface{})[0].(map[string]interface{})["MemberList"]
	man := 0
	woman := 0
	for _, v := range RoomContactList.([]interface{}) {
		if m, ok := v.([]interface{}); ok {
			for _, s := range m {
				members = append(members, s.(map[string]interface{})["UserName"].(string))
			}
		} else {
			members = append(members, v.(map[string]interface{})["UserName"].(string))
		}
	}
	urlstr = fmt.Sprintf("%s/webwxbatchgetcontact?type=ex&r=%s&pass_ticket=%s", self.base_uri, self._unixStr(), self.pass_ticket)
	length := 50
	debugPrint(members)
	mnum := len(members)
	block := int(math.Ceil(float64(mnum) / float64(length)))
	k := 0
	for k < block {
		offset := k * length
		var l int
		if offset+length > mnum {
			l = mnum
		} else {
			l = offset + length
		}
		blockmembers := members[offset:l]
		params := make(map[string]interface{})
		params["BaseRequest"] = self.BaseRequest
		params["Count"] = len(blockmembers)
		blockmemberslist := []map[string]string{}
		for _, g := range blockmembers {
			blockmemberslist = append(blockmemberslist, map[string]string{
				"UserName":        g,
				"EncryChatRoomId": chatroomId,
			})
		}
		params["List"] = blockmemberslist
		debugPrint(urlstr)
		debugPrint(params)
		dic, err := self._post(urlstr, params, true)
		if err == nil {
			debugPrint("flag")
			userlist := JsonDecode(dic).(map[string]interface{})["ContactList"]
			for _, u := range userlist.([]interface{}) {
				if u.(map[string]interface{})["Sex"].(int) == 1 {
					man++
				} else if u.(map[string]interface{})["Sex"].(int) == 2 {
					woman++
				}
			}
		}
		k++
	}
	stats = map[string]string{
		"woman": strconv.Itoa(woman),
		"man":   strconv.Itoa(man),
	}
	return stats, nil
}

func (self *wxweb) webwxsync() interface{} {
	urlstr := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s&pass_ticket=%s", self.base_uri, self.sid, self.skey, self.pass_ticket)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["SyncKey"] = self.SyncKey
	params["rr"] = ^int(time.Now().Unix())
	res, err := self._post(urlstr, params, true)
	if err != nil{
		return false
	}
	data := JsonDecode(res).(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	if retCode == 0 {
		self.SyncKey = data["SyncKey"].(map[string]interface{})
		self._setsynckey()
	}
	return data
}

func (self *wxweb) handleMsg(r interface{}) {
	myNickName := self.User["NickName"].(string)
	for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
		// fmt.Println("[*] 你有新的消息，请注意查收")
		// msg = msg.(map[string]interface{})
		msgType := msg.(map[string]interface{})["MsgType"].(int)
		fromUserName := msg.(map[string]interface{})["FromUserName"].(string)
		// name = self.getUserRemarkName(msg['FromUserName'])
		content := msg.(map[string]interface{})["Content"].(string)
		content = strings.Replace(content, "&lt;", "<", -1)
		content = strings.Replace(content, "&gt;", ">", -1)
		content = strings.Replace(content, " ", " ", 1)
		// msgid := msg.(map[string]interface{})["MsgId"].(int)
		if msgType == 1 {
			var ans string
			var err error
			if fromUserName[:2] == "@@" {
				contentSlice := strings.Split(content, ":<br/>")
				// people := contentSlice[0]
				content = contentSlice[1]
				if strings.Contains(content, "@"+myNickName) {
					realcontent := strings.TrimSpace(strings.Replace(content, "@"+myNickName, "", 1))
					debugPrint(realcontent)
					if realcontent == "统计人数" {
						stat, err := self.webgetchatroommember(fromUserName)
						if err == nil {
							ans = "据统计群里男生" + stat["man"] + "人，女生" + stat["woman"] + "人 (ó㉨ò)"
						}
					} else {
						ans, err = self.getReplyByApi(realcontent, fromUserName)
					}
				} else if strings.Contains(content, "撩@") {
					name := strings.Replace(content, "撩@", "", 1)
					name = strings.Replace(name, "\u003cbr/\u003e", "", 1)
					ans, err = self.getReplyByApi(LOVEWORDS_QUEST, fromUserName)
					if err == nil {
						ans = "@" + name + " " + ans
					}
				} else if content == "撩我" {
					ans, err = self.getReplyByApi(LOVEWORDS_QUEST, fromUserName)
				}
			} else {
				ans, err = self.getReplyByApi(content, fromUserName)
			}
			debugPrint(ans)
			debugPrint(content)
			if err != nil {
				debugPrint(err)
			} else if ans != "" {
				go self.webwxsendmsg(ans, fromUserName)
			}
		} else if msgType == 51 {
			fmt.Println("[*] 成功截获微信初始化消息")
		}
	}
}

func (self *wxweb) getReplyByApi(realcontent string, fromUserName string) (string, error) {
	return getAnswer(realcontent, fromUserName, self.User["NickName"].(string))
}

func (self *wxweb) webwxsendmsg(message string, toUseNname string) bool {
	urlstr := fmt.Sprintf("%s/webwxsendmsg?pass_ticket=%s", self.base_uri, self.pass_ticket)
	clientMsgId := self._unixStr() + "0" + strconv.Itoa(rand.Int())[3:6]
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = message
	msg["FromUserName"] = self.User["UserName"]
	msg["ToUserName"] = toUseNname
	msg["LocalID"] = clientMsgId
	msg["ClientMsgId"] = clientMsgId
	params["Msg"] = msg
	data, err := self._post(urlstr, params, true)
	if err != nil {
		debugPrint(err)
		return false
	} else {
		debugPrint(data)
		return true
	}
}

func (self *wxweb) _init() {
	gCookieJar, _ := cookiejar.New(nil)
	httpclient := http.Client{
		CheckRedirect: nil,
		Jar:           gCookieJar,
	}
	self.http_client = &httpclient
	rand.Seed(time.Now().Unix())
	str := strconv.Itoa(rand.Int())
	self.deviceId = "e" + str[2:17]
}

func (self *wxweb) test() {

}

func (self *wxweb) start() {
	fmt.Println("[*] 微信网页版 ... 开动")
	self._init()
	self._run("[*] 正在获取 uuid ... ", self.getUuid)
	self._run("[*] 正在获取 二维码 ... ", self.genQRcode)
	fmt.Println("[*] 请使用微信扫描二维码以登录 ... ")
	for {
		if self.waitForLogin(1) == false {
			continue
		}
		fmt.Println("[*] 请在手机上点击确认以登录 ... ")
		if self.waitForLogin(0) == false {
			continue
		}
		break
	}
	self._run("[*] 正在登录 ... ", self.login)
	self._run("[*] 微信初始化 ... ", self.webwxinit)
	self._run("[*] 开启状态通知 ... ", self.webwxstatusnotify)
	self._run("[*] 进行同步线路测试 ... ", self.testsynccheck)
	for {
		retcode, selector := self.synccheck()
		if retcode == "1100" {
			fmt.Println("[*] 你在手机上登出了微信，债见")
			break
		} else if retcode == "1101" {
			fmt.Println("[*] 你在其他地方登录了 WEB 版微信，债见")
			break
		} else if retcode == "0" {
			if selector == "2" {
				r := self.webwxsync()
				debugPrint(r)
				switch r.(type) {
				case bool:
				default:
					self.handleMsg(r)
				}
			} else if selector == "0" {
				time.Sleep(1)
			} else if selector == "6" || selector == "4" {
				self.webwxsync()
				time.Sleep(1)
			}
		}
	}

}

func forgeheadget(urlstr string) string {

	client := &http.Client{}

	reqest, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate, sdch")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Host", "login.weixin.qq.com")
	reqest.Header.Add("Referer", "https://wx.qq.com/")
	reqest.Header.Add("Upgrade-Insecure-Requests", "1")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}
