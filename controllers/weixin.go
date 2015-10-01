package controllers

import (
	// "encoding/xml"
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	"crypto/sha1"
	"encoding/xml"
	"io"
	// "os"
	"sort"
	"strings"
	"time"
)

//weixin request, 接收订单查询请求，返回地图信息
func (this *MainController) ReceiveMsg() {
	body := this.Ctx.Input.CopyBody()
	// this.Ctx.Request.
	DebugTraceF("输入：%s", string(body))
	response := ""
	defer func() {
		DebugTraceF("输出：%s", response)
		this.Ctx.WriteString(response)
	}()
	if msg, err := parseComingInMessage(body); err != nil {
		DebugSysF("解析接收到的微信消息时发生错误：%s", err.Error())
		return
	} else {
		DebugTraceF(msg.String())
		/*
			查找快递相应的位置图片，没有则根据情况使用默认图片
			如果图片下载失败，使用正在查找的图片代替
		*/
		bagageID := msg.Content
		var weixinRes *weixinResponseNews
		defer func() {
			if bytes, err := xml.Marshal(weixinRes); err != nil {
				DebugSysF("序列化返回信息出错：%s", err.Error())
				return
			} else {
				response = string(bytes)
			}
		}()
		//如果没有该单号
		bagagePosInfo := g_bagages.findOne(func(b *BagagePosInfo) bool { return b.BagageID == bagageID })
		if bagagePosInfo == nil {
			//构造一个没有单号的返回信息
			imageUrl := ImageDirDefaultPath + "error.png"
			articleList := ArticleList{&Article{NewArticleItem("订单状态查询", "没有该订单信息", imageUrl, imageUrl)}}
			weixinRes = NewWeixinResponseNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articleList)
			return
		}
		if bagagePosInfo.ImageName == "" {
			// 使用没有地图的默认图片构造返回信息
			imageUrl := ImageDirDefaultPath + "rt.png"
			articleList := ArticleList{&Article{NewArticleItem("订单状态查询", "暂未找到订单的位置", imageUrl, imageUrl)}}
			weixinRes = NewWeixinResponseNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articleList)
			return
		}
		//有单号有图片
		imageUrl := imageDirPath + bagagePosInfo.ImageName
		articleList := ArticleList{&Article{NewArticleItem("订单状态查询", fmt.Sprintf("单号 %s 最新位置 %s", bagageID, bagagePosInfo.TimeStamp), imageUrl, imageUrl)}}
		weixinRes = NewWeixinResponseNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articleList)
	}
}

//解析微信的xml消息
/*
<xml>
	 <ToUserName><![CDATA[toUser]]></ToUserName>
	 <FromUserName><![CDATA[fromUser]]></FromUserName>
	 <CreateTime>1348831860</CreateTime>
	 <MsgType><![CDATA[text]]></MsgType>
	 <Content><![CDATA[this is a test]]></Content>
	 <MsgId>1234567890123456</MsgId>
 </xml>
*/
func parseComingInMessage(body []byte) (*weixinInputTextMsg, error) {
	var msg weixinInputTextMsg
	if err := xml.Unmarshal(body, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

//weixin server test request
func (this *MainController) Index() {
	signature := this.GetString("signature")
	timestamp := this.GetString("timestamp")
	nonce := this.GetString("nonce")
	echostr := this.GetString("echostr")
	if isLegel(signature, timestamp, nonce, token) == true {
		// res.send(echostr)
		this.Ctx.WriteString(echostr)
	} else {
		// res.send('')
		// this.
		this.Ctx.WriteString("")
	}
}

func isLegel(signature, timestamp, nonce, token string) bool {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil)) == signature
}

func NewWeixinResponseNews(to, from string, time int64, articles ArticleList) *weixinResponseNews {
	return &weixinResponseNews{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time,
		MsgType:      "news",
		ArticleCount: len(articles),
		Articles:     articles,
	}
}
func NewArticleItem(title, des, picurl, url string) *ArticleItem {
	return &ArticleItem{
		Title:       title,
		Description: des,
		PicUrl:      picurl,
		Url:         url,
	}
}

func (this *weixinInputTextMsg) String() string {
	return fmt.Sprintf("From %10s To %10s   Content:%s", this.FromUserName, this.ToUserName, this.Content)
}
