package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gogf/gf/encoding/gurl"
)

// ServerChanAPI 配置
type ServerChanAPI struct {
	SCKEY string // server酱的key http://sc.ftqq.com
	Text  string // 标题内容 text：消息标题，最长为256，必填。
	Desp  string // 内容文本 desp：消息内容，最长64Kb，可空，支持MarkDown。
}

// ServerChanPost 转发至server酱公众号的信息！
func (s ServerChanAPI) ServerChanPost() {
	sckey := gurl.Encode(s.SCKEY) // 对key 进行url编码
	text := gurl.Encode(s.Text)   // 对 标题进行 url 编码
	desp := gurl.Encode(s.Desp)   // 对内容进行url编码
	// 将组合好的方式用post发送出去
	resp, err := http.Post("https://sc.ftqq.com/"+sckey+".send?text="+text+"&desp="+desp, "text/html;charset=utf-8", nil)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close() // 关闭连接
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	if strings.Contains(string(body), "success") == true { // 对返回的body 进行判断 是否包含success
		log.Println(s.Text + "  的消息已经推送到微信")
	} else {
		log.Println("微信消息推送失败请重试，please try again")

	}

}
