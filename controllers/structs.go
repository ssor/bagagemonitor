package controllers

import (
	"encoding/xml"
	// "fmt"
)

//传入位置信息时的数据格式
type Pakage struct {
	ID         string
	Longitude  float64
	Latitude   float64
	BagageList []string
}

//订单的位置信息
type BagagePosInfo struct {
	BagageID  string
	TimeStamp string
	Longitude float64
	Latitude  float64
	Flag      string //放置在位置图上的文字标识
	ImageName string //下载的位置图片
}
type BagagePosInfoList []*BagagePosInfo

// type ImageInfo struct {
// 	Name    string
// 	Deleted bool
// }
// type ImageInfoList []*ImageInfo

//weixin input parsed result
type weixinInputTextMsg struct {
	ToUserName, FromUserName string
	CreateTime               string
	MsgType                  string
	Content                  string
	MsgId                    string
}

//weixin response type
type weixinResponseNews struct {
	XMLName                  xml.Name `xml:"xml"`
	ToUserName, FromUserName string
	CreateTime               int64
	MsgType                  string
	ArticleCount             int
	Articles                 ArticleList
}

//weixin url
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
