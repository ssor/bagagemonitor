package controllers

import (
	"encoding/xml"
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	// "os"
	// "time"
)

type weixinResponseNews struct {
	XMLName                  xml.Name `xml:"xml"`
	ToUserName, FromUserName string
	CreateTime               int64
	MsgType                  string
	ArticleCount             int
	Articles                 ArticleList
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

type ArticleList []*Article

type Article struct {
	Item *ArticleItem `xml:"item"`
}
type ArticleItem struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}

func NewArticleItem(title, des, picurl, url string) *ArticleItem {
	return &ArticleItem{
		Title:       title,
		Description: des,
		PicUrl:      picurl,
		Url:         url,
	}
}

type weixinInputTextMsg struct {
	ToUserName, FromUserName string
	CreateTime               string
	MsgType                  string
	Content                  string
	MsgId                    string
}

func (this *weixinInputTextMsg) String() string {
	return fmt.Sprintf("From %10s To %10s   Content:%s", this.FromUserName, this.ToUserName, this.Content)
}
