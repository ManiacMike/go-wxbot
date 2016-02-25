//web weixin client
package main

import (
  "fmt"
  "net/http"
  // "log"
	"net/url"
  "time"
  "regexp"
  "os"
  "os/exec"
  "io/ioutil"
  "io"
  "bytes"
  "strconv"
  "strings"
  // "strings"
)

type wxweb struct {
  uuid string
}

type jsonData map[string]interface{}

func (self *wxweb) getUuid(args ...interface{}) bool{
  urlstr := "https://login.weixin.qq.com/jslogin"
  params := make(map[string]interface{})
  params["appid"] = "wx782c26e4c19acffb"
	params["fun"] = "new"
	params["lang"] = "zh_CN"
  params["_"] = self.unixStr()
  data,_ := self._post(urlstr, params, false)
  re := regexp.MustCompile("\"(\\S+?)\"")
  find := re.FindAllString(data, -1)
  if len(find) > 0{
    self.uuid = strings.Replace(find[0],"\"","",-1)
    return true
  }else{
    return false
  }
}

func (self *wxweb) _run(desc string, f func(...interface{}) bool, args... interface{}) {
  fmt.Print(desc)
  var result bool
  if len(args) > 1 {
    result = f(args)
  } else if len(args) == 1 {
    result = f(args[0])
  } else {
    result = f()
  }
  if result{
    fmt.Println("成功")
  }else{
    fmt.Println("失败\n[*] 退出程序");
    os.Exit(1)
  }
}

func (self *wxweb) _post(urlstr string, params jsonData, jsonFmt bool) (string, error){
  var err error
  res := ""
  postParams := url.Values{}
  for key,value := range params{
    postParams.Add(key, value.(string))
  }
  resp, err := http.PostForm(urlstr, postParams)
  if err != nil{
    return res,err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil{
    return res,err
  }
  return string(body),nil
}

func (self *wxweb) unixStr() string{
  return strconv.Itoa(int(time.Now().Unix()))
}

//TODO support linux
func (self *wxweb) genQRcode(args ...interface{}) bool{
  urlstr := "https://login.weixin.qq.com/qrcode/" + self.uuid
  urlstr += "?t=webwx"
  urlstr += "&_=" + self.unixStr()
  path := "qrcode.jpg"
  out, err := os.Create(path)
  defer out.Close()
  resp, err := http.Get(urlstr)
  defer resp.Body.Close()
  pix, err := ioutil.ReadAll(resp.Body)
  _, err = io.Copy(out, bytes.NewReader(pix))
  if err != nil{
    return false
  }else{
    exec.Command("open",path).Run()
    return true
  }
}

func (self *wxweb) start(){
  fmt.Println("[*] 微信网页版 ... 开动")
  self._run("[*] 正在获取 uuid ... ", self.getUuid)
  self._run("[*] 正在获取 二维码 ... ", self.genQRcode)
}
