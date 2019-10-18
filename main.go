package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gogf/gf/os/gcfg"
)

// Desp 是消息内容 server酱用的
var Desp string

// Text 是消息标题 server酱用的
var Text string

func main() {

	gcfg.Instance().SetFileName("google_to_wechat_config.toml") // 指定配置文件名字
	email := gcfg.Instance().GetString("google.email")          // 读取配置文件的邮箱
	passwrod := gcfg.Instance().GetString("google.password")    // 读取配置文件的密码
	sckey := gcfg.Instance().GetString("ServerChan.sckey")      // 读取配置文件的server酱KEY

	ctx := context.Background()
	// 下面这段代码是打开chrome界面进行调试的。 正式版本应该注销
	// options := []chromedp.ExecAllocatorOption{
	// 	chromedp.Flag("headless", false),
	// 	chromedp.Flag("hide-scrollbars", false),
	// 	chromedp.Flag("mute-audio", false),
	// }
	options := []chromedp.ExecAllocatorOption{ // 这个是正式环境使用的
		chromedp.Flag("headless", false), // 因为谷歌限制策略只能以带窗口模式运行，也方便我们观察
		chromedp.DisableGPU,              // 禁用chrome 的GPU 因为服务器不需要
		chromedp.NoDefaultBrowserCheck,   // 禁用默认浏览器检查
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...) // 切片

	c, cc := chromedp.NewExecAllocator(ctx, options...)
	defer cc()
	// create context
	ctx, cancel := chromedp.NewContext(c)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.google.com/accounts/Login?hl=zh-CN`),   //浏览网址
		chromedp.WaitVisible(`input[type="email"]`, chromedp.NodeVisible),     // 等待元素加载完毕
		chromedp.Sleep(2*time.Second),                                         // 等待 2秒
		chromedp.SendKeys(`input[type="email"]`, email, chromedp.NodeVisible), // 给指定元素框 输入email数据
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`#identifierNext`, chromedp.ByID), // 点击下一步
		chromedp.WaitVisible(`input[type="password"]`, chromedp.NodeVisible),
		chromedp.Sleep(3*time.Second),
		chromedp.SendKeys(`input[type="password"]`, passwrod, chromedp.ByQuery),
		chromedp.Click(`#passwordNext`, chromedp.ByID), // 如果启用两步验证 建议打开 Google prompt 这样可以直接点确认 或者直接关闭两步验证
		chromedp.WaitVisible(`input[aria-label="Search Google Account"]`, chromedp.NodeVisible),
		chromedp.Navigate(`https://voice.google.com/messages`),

		// -------------------------------
		// 发送消息
		// chromedp.Click(`div[aria-label="Send new message"]`, chromedp.NodeVisible),
		// chromedp.SendKeys(`label[ng-click="delegateClick()"]`, "8336721001‬", chromedp.NodeVisible), // 回车没搞定
		// chromedp.Sleep(1*time.Second),
		// chromedp.Click(`textarea[aria-label="Type a message"]`, chromedp.NodeVisible),
		// chromedp.SendKeys(`textarea[aria-label="Type a message"]`, "shaoxia.xyz", chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatalln("登陆有误", err)
	}
	// 这里开始无限循环获取新的消息。
	for {
		err := chromedp.Run(ctx, WaitNewMessages())
		if err != nil {
			log.Fatalln("获取新消息有误", err)
		}
		log.Println("收到" + strings.TrimSpace(Text) + "发来的消息:" + strings.TrimSpace(Desp))
		// 当获取到新消息之后，就推送到 server酱
		s := ServerChanAPI{
			SCKEY: sckey,
			Text:  Text,
			Desp:  Desp,
		}
		s.ServerChanPost()
		time.Sleep(2 * time.Second)
	}

}

// WaitNewMessages 等待新的消息出来
func WaitNewMessages() chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.WaitVisible(`a[aria-label="Messages: 1 unread"]`, chromedp.NodeVisible), // 这个表示 消息未读 （点亮的 数字 1）
		chromedp.Click(`#messaging-view > div > md-content > div > gv-conversation-list > md-virtual-repeat-container > div > div.md-virtual-repeat-offsetter > div:nth-child(1) > div > gv-text-thread-item > gv-thread-item > div`, chromedp.NodeVisible),                                                                                        // 选择未读消息 (主要是将消息更改成已读)
		chromedp.Text(`#messaging-view > div > md-content > div > gv-conversation-list > md-virtual-repeat-container > div > div.md-virtual-repeat-offsetter > div:nth-child(1) > div > gv-text-thread-item > gv-thread-item > div > div.rkljfb-MZArnb.flex > div > gv-annotation`, &Text, chromedp.NodeVisible),                                   // 选择要发送的消息 手机号码
		chromedp.Text(`#messaging-view > div > md-content > div > gv-conversation-list > md-virtual-repeat-container > div > div.md-virtual-repeat-offsetter > div:nth-child(1) > div > gv-text-thread-item > gv-thread-item > div > div.rkljfb-MZArnb.flex > ng-transclude > gv-thread-item-detail > gv-annotation`, &Desp, chromedp.NodeVisible), // 选择要发送的消息 消息内容
	}
}

//*-----------------------------------

// 思路 1 打开gui 浏览器，等待用户输入账号密码并登陆成功，
// 检测GUI浏览器关闭
// 启动隐藏浏览器继续爬虫 gv信息

// https://zhangguanzhang.github.io/2019/07/14/chromedp/

// https://blog.csdn.net/yang731227/article/details/89202458
