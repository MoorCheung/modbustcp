package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/chanxuehong/wechat/mp/core"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func Get(str string){
	url := "http://hq.sinajs.cn/list=%s"
	url = fmt.Sprintf(url,str)
	resp,err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}else{
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}else{
			src:=string(bytes)
			srcDecoder := mahonia.NewDecoder("gbk")
			desDecoder := mahonia.NewDecoder("utf-8")
			resStr:= srcDecoder.ConvertString(src)
			_, resBytes, _ := desDecoder .Translate([]byte(resStr), true)
			src = string(resBytes)
			src = src[:len(src)-6]
			fmt.Println(src)
			strs:= strings.Split(src,`var hq_str_sh600101="`)
			arr := strings.Split(strs[1],",")
			fmt.Println(arr)
			// 查找行首以 H 开头，以空格结尾的字符串 $guize = '/^var hq_str_(.*)="(.*)"/';
			re := regexp.MustCompile(`^var hq_str_(\S)*=\"(.*)`)
			fmt.Println(re.FindAllStringSubmatch(src, -1))
		}

	}

}
const (
	wxAppId     = "appid"
	wxAppSecret = "appsecret"

	wxOriId         = "oriid"
	wxToken         = "token"
	wxEncodedAESKey = "aeskey"
)

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server
)

func init() {
	mux := core.NewServeMux()

	msgHandler = mux
	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}
func main(){

	log.Println(http.ListenAndServe(":80", nil))
	//Get("sh600101")
}