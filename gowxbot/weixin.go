// web weixin client
package gowxbot

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func debugPrint(content interface{}) {
	fmt.Println(content)
}

type Wxweb struct {
	uuid        string
	baseURI     string
	redirectURI string
	uin         string
	sid         string
	skey        string
	passTicket  string
	deviceID    string
	SyncKey     map[string]interface{}
	synckey     string
	User        map[string]interface{}
	BaseRequest map[string]interface{}
	syncHost    string
	httpClient  *http.Client
}

func (w *Wxweb) getUUID(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/jslogin"
	urlstr += "?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + w._unixStr()
	data, _ := w._get(urlstr, false)
	re := regexp.MustCompile(`"([\S]+)"`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		w.uuid = find[1]
		return true
	}
	return false

}

func (w *Wxweb) _run(desc string, f func(...interface{}) bool, args ...interface{}) {
	start := time.Now().UnixNano()
	var result bool
	if len(args) > 1 {
		result = f(args)
	} else if len(args) == 1 {
		result = f(args[0])
	} else {
		result = f()
	}
	useTime := fmt.Sprintf("%.5f", float64(time.Now().UnixNano()-start)/1000000000)
	if result {
		fmt.Println("成功,用时" + useTime + "秒")
	} else {
		fmt.Println("失败\n[*] 退出程序")
		os.Exit(1)
	}
}

func (w *Wxweb) _post(urlstr string, params map[string]interface{}, jsonFmt bool) ([]byte, error) {
	var err error
	var resp *http.Response
	if jsonFmt == true {
		jsonPost, err := json.Marshal(params)
		if err != nil {
			return []byte(""), err
		}

		debugPrint(jsonPost)
		requestBody := bytes.NewBuffer([]byte(jsonPost))
		request, err := http.NewRequest("POST", urlstr, requestBody)
		if err != nil {
			return []byte(""), err
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
		request.Header.Add("Referer", "https://wx.qq.com/")
		request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
		resp, err = w.httpClient.Do(request)
	} else {
		v := url.Values{}
		for key, value := range params {
			v.Add(key, value.(string))
		}
		resp, err = w.httpClient.PostForm(urlstr, v)
	}

	if err != nil || resp == nil {
		fmt.Println(err)
		return []byte(""), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return []byte(""), err
	}
	return body, nil
}

func (w *Wxweb) _get(urlstr string, jsonFmt bool) (string, error) {
	var err error
	res := ""
	request, _ := http.NewRequest("GET", urlstr, nil)
	request.Header.Add("Referer", "https://wx.qq.com/")
	request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	resp, err := w.httpClient.Do(request)
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

func (w *Wxweb) _unixStr() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (w *Wxweb) genQRcode(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/qrcode/" + w.uuid
	urlstr += "?t=webwx"
	urlstr += "&_=" + w._unixStr()
	path := "qrcode.jpg"
	out, err := os.Create(path)
	resp, err := w._get(urlstr, false)
	_, err = io.Copy(out, bytes.NewReader([]byte(resp)))
	if err != nil {
		return false
	}
	if runtime.GOOS == "darwin" {
		exec.Command("open", path).Run()
	} else {
		go func() {
			http.HandleFunc("/qrcode", func(w http.ResponseWriter, req *http.Request) {
				http.ServeFile(w, req, "qrcode.jpg")
				return
			})
			http.ListenAndServe(":8889", nil)
		}()
	}
	return true

}

func (w *Wxweb) waitForLogin(tip int) bool {
	time.Sleep(time.Duration(tip) * time.Second)
	url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
	url += "?tip=" + strconv.Itoa(tip) + "&uuid=" + w.uuid + "&_=" + w._unixStr()
	data, _ := w._get(url, false)
	re := regexp.MustCompile(`window.code=(\d+);`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		code := find[1]
		if code == "201" {
			return true
		} else if code == "200" {
			re := regexp.MustCompile(`window.redirectURI="(\S+?)";`)
			find := re.FindStringSubmatch(data)
			if len(find) > 1 {
				rURI := find[1] + "&fun=new"
				w.redirectURI = rURI
				re = regexp.MustCompile(`/`)
				finded := re.FindAllStringIndex(rURI, -1)
				w.baseURI = rURI[:finded[len(finded)-1][0]]
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

func (w *Wxweb) login(args ...interface{}) bool {
	data, _ := w._get(w.redirectURI, false)
	type Result struct {
		Skey       string `xml:"skey"`
		Wxsid      string `xml:"wxsid"`
		Wxuin      string `xml:"wxuin"`
		PassTicket string `xml:"passTicket"`
	}
	v := Result{}
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}
	w.skey = v.Skey
	w.sid = v.Wxsid
	w.uin = v.Wxuin
	w.passTicket = v.PassTicket
	w.BaseRequest = make(map[string]interface{})
	w.BaseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	w.BaseRequest["Sid"] = v.Wxsid
	w.BaseRequest["Skey"] = v.Skey
	w.BaseRequest["deviceID"] = w.deviceID
	return true
}

func (w *Wxweb) webwxinit(args ...interface{}) bool {
	url := fmt.Sprintf("%s/webwxinit?passTicket=%s&skey=%s&r=%s", w.baseURI, w.passTicket, w.skey, w._unixStr())
	params := make(map[string]interface{})
	params["BaseRequest"] = w.BaseRequest
	res, err := w._post(url, params, true)
	if err != nil {
		return false
	}
	ioutil.WriteFile("tmp.txt", []byte(res), 777)
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return false
	}
	w.User = data["User"].(map[string]interface{})
	w.SyncKey = data["SyncKey"].(map[string]interface{})
	w._setsynckey()

	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(float64)
	return retCode == 0
}

func (w *Wxweb) _setsynckey() {
	keys := []string{}
	for _, keyVal := range w.SyncKey["List"].([]interface{}) {
		key := strconv.Itoa(int(keyVal.(map[string]interface{})["Key"].(float64)))
		value := strconv.Itoa(int(keyVal.(map[string]interface{})["Val"].(float64)))
		keys = append(keys, key+"_"+value)
	}
	w.synckey = strings.Join(keys, "|")
	debugPrint(w.synckey)
}

func (w *Wxweb) synccheck() (string, string) {
	urlstr := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", w.syncHost)
	v := url.Values{}
	v.Add("r", w._unixStr())
	v.Add("sid", w.sid)
	v.Add("uin", w.uin)
	v.Add("skey", w.skey)
	v.Add("deviceID", w.deviceID)
	v.Add("synckey", w.synckey)
	v.Add("_", w._unixStr())
	urlstr = urlstr + "?" + v.Encode()
	data, _ := w._get(urlstr, false)
	re := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	find := re.FindStringSubmatch(data)
	if len(find) > 2 {
		retcode := find[1]
		selector := find[2]
		debugPrint(fmt.Sprintf("retcode:%s,selector,selector%s", find[1], find[2]))
		return retcode, selector
	}
	return "9999", "0"

}

func (w *Wxweb) testsynccheck(args ...interface{}) bool {
	SyncHost := []string{
		"webpush.wx.qq.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
		//"webpush.wechatapp.com"
	}
	for _, host := range SyncHost {
		w.syncHost = host
		retcode, _ := w.synccheck()
		if retcode == "0" {
			return true
		}
	}
	return false
}

func (w *Wxweb) webwxstatusnotify(args ...interface{}) bool {
	urlstr := fmt.Sprintf("%s/webwxstatusnotify?lang=zh_CN&passTicket=%s", w.baseURI, w.passTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = w.BaseRequest
	params["Code"] = 3
	params["FromUserName"] = w.User["UserName"]
	params["ToUserName"] = w.User["UserName"]
	params["ClientMsgId"] = int(time.Now().Unix())
	res, err := w._post(urlstr, params, true)
	if err != nil {
		return false
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return false
	}
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(float64)
	return retCode == 0
}

func (w *Wxweb) webgetchatroommember(chatroomID string) (map[string]string, error) {
	urlstr := fmt.Sprintf("%s/webwxbatchgetcontact?type=ex&r=%s&passTicket=%s", w.baseURI, w._unixStr(), w.passTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = w.BaseRequest
	params["Count"] = 1
	params["List"] = []map[string]string{}
	l := []map[string]string{}
	params["List"] = append(l, map[string]string{
		"UserName":   chatroomID,
		"chatroomID": "",
	})
	members := []string{}
	stats := make(map[string]string)
	res, err := w._post(urlstr, params, true)
	fmt.Println(urlstr)
	debugPrint(params)
	if err != nil {
		return stats, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return stats, err
	}
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
	urlstr = fmt.Sprintf("%s/webwxbatchgetcontact?type=ex&r=%s&passTicket=%s", w.baseURI, w._unixStr(), w.passTicket)
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
		params["BaseRequest"] = w.BaseRequest
		params["Count"] = len(blockmembers)
		blockmemberslist := []map[string]string{}
		for _, g := range blockmembers {
			blockmemberslist = append(blockmemberslist, map[string]string{
				"UserName":        g,
				"EncrychatroomID": chatroomID,
			})
		}
		params["List"] = blockmemberslist
		debugPrint(urlstr)
		debugPrint(params)
		dic, err := w._post(urlstr, params, true)
		if err == nil {
			userlistTmp := make(map[string]interface{})
			err = json.Unmarshal(dic, &userlistTmp)
			if err == nil {
				userlist := userlistTmp["ContactList"]
				for _, u := range userlist.([]interface{}) {
					if u.(map[string]interface{})["Sex"].(float64) == 1 {
						man++
					} else if u.(map[string]interface{})["Sex"].(float64) == 2 {
						woman++
					}
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

func (w *Wxweb) webwxsync() interface{} {
	urlstr := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s&passTicket=%s", w.baseURI, w.sid, w.skey, w.passTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = w.BaseRequest
	params["SyncKey"] = w.SyncKey
	params["rr"] = ^int(time.Now().Unix())
	res, err := w._post(urlstr, params, true)
	if err != nil {
		return false
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return false
	}
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(float64)
	if retCode == 0 {
		w.SyncKey = data["SyncKey"].(map[string]interface{})
		w._setsynckey()
	}
	return data
}

func (w *Wxweb) handleMsg(r interface{}) {
	myNickName := w.User["NickName"].(string)
	for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
		// fmt.Printf("[*] message: %v \n", msg)
		// msg = msg.(map[string]interface{})
		msgType := msg.(map[string]interface{})["MsgType"].(float64)
		fromUserName := msg.(map[string]interface{})["FromUserName"].(string)
		// name = Wxweb.getUserRemarkName(msg['FromUserName'])
		content := msg.(map[string]interface{})["Content"].(string)
		content = strings.Replace(content, "&lt;", "<", -1)
		content = strings.Replace(content, "&gt;", ">", -1)
		content = strings.Replace(content, " ", " ", 1)
		// msgid := msg.(map[string]interface{})["MsgId"].(float64)
		if msgType == 1 {
			var ans string
			var err error
			if fromUserName[:2] == "@@" {
				debugPrint(content + "|0045")
				contentSlice := strings.Split(content, ":<br/>")
				// people := contentSlice[0]
				content = contentSlice[1]
				if strings.Contains(content, "@"+myNickName) {
					realcontent := strings.TrimSpace(strings.Replace(content, "@"+myNickName, "", 1))
					debugPrint(realcontent + "|0046")
					if realcontent == "统计人数" {
						stat, err := w.webgetchatroommember(fromUserName)
						if err == nil {
							ans = "据统计群里男生" + stat["man"] + "人，女生" + stat["woman"] + "人 (ó㉨ò)"
						}
					} else {
						ans, err = w.getReplyByAPI(realcontent, "", fromUserName, "")
					}
				} else if strings.Contains(content, "撩@") {
					name := strings.Replace(content, "撩@", "", 1)
					name = strings.Replace(name, "\u003cbr/\u003e", "", 1)
					ans, err = w.getReplyByAPI(LoveWordsQuest, "", fromUserName, "")
					if err == nil {
						ans = "@" + name + " " + ans
					}
				} else if content == "撩我" {
					ans, err = w.getReplyByAPI(LoveWordsQuest, "", fromUserName, "")
				}
			} else {
				ans, err = w.getReplyByAPI(content, fromUserName, "", "")
			}
			debugPrint(ans)
			debugPrint(content)
			if err != nil {
				debugPrint(err)
			} else if ans != "" {
				go w.webwxsendmsg(ans, fromUserName)
			}
		} else if msgType == 51 {
			fmt.Println("[*] 成功截获微信初始化消息")
		}
	}
}

func (w *Wxweb) getReplyByAPI(realcontent, fromUserName, groupID, userIDName string) (string, error) {
	username := fromUserName[1:33]
	return getAnswer(realcontent, username, groupID, userIDName, w.User["NickName"].(string))
}

func (w *Wxweb) webwxsendmsg(message string, toUseNname string) bool {
	urlstr := fmt.Sprintf("%s/webwxsendmsg?passTicket=%s", w.baseURI, w.passTicket)
	clientMsgID := w._unixStr() + "0" + strconv.Itoa(rand.Int())[3:6]
	params := make(map[string]interface{})
	params["BaseRequest"] = w.BaseRequest
	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = message
	msg["FromUserName"] = w.User["UserName"]
	msg["ToUserName"] = toUseNname
	msg["LocalID"] = clientMsgID
	msg["ClientMsgId"] = clientMsgID
	params["Msg"] = msg
	data, err := w._post(urlstr, params, true)
	if err != nil {
		debugPrint(err)
		return false
	}
	debugPrint(data)
	return true

}

func (w *Wxweb) _init() {
	gCookieJar, _ := cookiejar.New(nil)
	httpclient := http.Client{
		CheckRedirect: nil,
		Jar:           gCookieJar,
	}
	w.httpClient = &httpclient
	rand.Seed(time.Now().Unix())
	str := strconv.Itoa(rand.Int())
	w.deviceID = "e" + str[2:17]
}

func (w *Wxweb) test() {

}

func (w *Wxweb) Start() {
	fmt.Println("[*] 微信网页版 ... 开动")
	w._init()
	w._run("[*] 正在获取 uuid ... ", w.getUUID)
	w._run("[*] 正在获取 二维码 ... ", w.genQRcode)
	if runtime.GOOS == "darwin" {
		fmt.Println("[*] 请使用微信扫描二维码以登录 ... ")
	} else {
		fmt.Println("[*] 打开链接扫码登录 http://127.0.0.1:8889/qrcode")
	}
	for {
		if w.waitForLogin(1) == false {
			continue
		}
		fmt.Println("[*] 请在手机上点击确认以登录 ... ")
		if w.waitForLogin(0) == false {
			continue
		}
		break
	}
	w._run("[*] 正在登录 ... ", w.login)
	w._run("[*] 微信初始化 ... ", w.webwxinit)
	w._run("[*] 开启状态通知 ... ", w.webwxstatusnotify)
	w._run("[*] 进行同步线路测试 ... ", w.testsynccheck)
	for {
		retcode, selector := w.synccheck()
		if retcode == "1100" {
			fmt.Println("[*] 你在手机上登出了微信，债见")
			break
		} else if retcode == "1101" {
			fmt.Println("[*] 你在其他地方登录了 WEB 版微信，债见")
			break
		} else if retcode == "0" {
			if selector == "2" {
				r := w.webwxsync()
				debugPrint(r)
				switch r.(type) {
				case bool:
				default:
					w.handleMsg(r)
				}
			} else if selector == "0" {
				time.Sleep(1)
			} else if selector == "6" || selector == "4" {
				w.webwxsync()
				time.Sleep(1)
			}
		}
	}

}

func forgeheadget(urlstr string) ([]byte, error) {

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
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
	body, _ := ioutil.ReadAll(response.Body)
	return body, nil
}
