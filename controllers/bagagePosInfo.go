package controllers

import (
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	// "os"
	// "time"
)

//订单的位置信息
type BagagePosInfo struct {
	BagageID       string
	TimeStamp      string
	SogouLongitude string
	SogouLatitude  string
	Flag           string //放置在位置图上的文字标识
	ImageName      string //下载的位置图片
}

func NewBagagePosInfo(bagageID, timeStamp, lon, lat, flag string) *BagagePosInfo {
	return &BagagePosInfo{
		BagageID:       bagageID,
		TimeStamp:      timeStamp,
		Flag:           flag,
		SogouLongitude: lon,
		SogouLatitude:  lat,
	}
}
func (this *BagagePosInfo) String() string {
	return fmt.Sprintf("Flag: %10s  BagageID:%10s   time: %10s     Position:(%6s, %6s)  Image: %s",
		this.Flag, this.BagageID, this.TimeStamp, this.SogouLongitude, this.SogouLatitude, this.ImageName)
}
func (this *BagagePosInfo) update(timeStamp, lon, lat, imageName, flag string) {
	this.TimeStamp = timeStamp
	this.SogouLongitude = lon
	this.SogouLatitude = lat
	this.ImageName = imageName
}

type BagagePosInfoList []*BagagePosInfo

func (this BagagePosInfoList) ListName() string {
	return "订单位置信息"
}
func (this BagagePosInfoList) InfoList() (list []string) {
	for _, bi := range this {
		list = append(list, bi.String())
	}
	return
}
func (this BagagePosInfoList) UsingImage(imageName string) bool {
	for _, bpi := range this {
		if bpi.ImageName == imageName {
			return true
		}
	}
	return false
}
func (this BagagePosInfoList) Find(bagageID string) *BagagePosInfo {
	for _, bi := range this {
		if bi.BagageID == bagageID {
			return bi
		}
	}
	return nil
}

//是否存在单号和坐标完全一致的订单信息
func (this BagagePosInfoList) BagageInfoRepeat(bpi *BagagePosInfo) bool {
	bi := this.Find(bpi.BagageID)
	if bi == nil {
		return false
	} else if bi.SogouLatitude == bpi.SogouLatitude && bi.SogouLongitude == bpi.SogouLongitude {
		return true
	}
	return false
}
