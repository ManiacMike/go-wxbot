//web weixin client
package main

import (
  "fmt"
  "net/http"
  // "log"
	"net/url"
  "io/ioutil"
  "time"
  "regexp"
  // "strings"
)

type wxweb struct {
  uuid string
}

type jsonData map[string]interface{}

func (self *wxweb) getUuid() bool{
  urlstr := "https://login.weixin.qq.com/jslogin"
  params := make(map[string]interface{})
  params["appid"] = "wx782c26e4c19acffb"
	params["fun"] = "new"
	params["lang"] = "zh_CN"
  params["_"] = string(time.Now().Unix())
  data,_ := self._post(urlstr, params, false)
  re := regexp.MustCompile("\"(\\S+?)\"")
  find := re.FindAllString(data, -1)
  if len(find) > 0{
    fmt.Println(find)
    self.uuid = find[0]
    return true
  }else{
    return false
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
