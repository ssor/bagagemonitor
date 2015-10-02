package controllers

import (
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	// "os"
	"time"
)

func NewBagagePosInfo(bagageID string, lon, lat float64, flag string) *BagagePosInfo {
	addedTime := time.Now().Format("2006-01-02 15:04:05")
	return &BagagePosInfo{
		BagageID:  bagageID,
		TimeStamp: addedTime,
		Flag:      flag,
		Longitude: lon,
		Latitude:  lat,
	}
}
func (this *BagagePosInfo) String() string {
	return fmt.Sprintf("Flag: %10s  BagageID:%10s   time: %10s     Position:(%6s, %6s)  Image: %s",
		this.Flag, this.BagageID, this.TimeStamp, this.Longitude, this.Latitude, this.ImageName)
}
func (b *BagagePosInfo) equals(bi *BagagePosInfo) bool {
	return b.Latitude == bi.Latitude && b.Longitude == bi.Longitude
}

func (this *BagagePosInfo) update(timeStamp string, lon, lat float64, imageName, flag string) {
	this.TimeStamp = timeStamp
	this.Longitude = lon
	this.Latitude = lat
	this.ImageName = imageName
}

func (this BagagePosInfoList) ListName() string {
	return "订单位置信息"
}
func (this BagagePosInfoList) InfoList() (list []string) {
	for _, bi := range this {
		list = append(list, bi.String())
	}
	return
}
func (bl BagagePosInfoList) forEach(f func(*BagagePosInfo)) {
	if len(bl) <= 0 {
		return
	}
	f(bl[0])
	bl[1:].forEach(f)
}

type BagagePosInfoPredictor func(*BagagePosInfo) bool

func (bl BagagePosInfoList) findOne(p BagagePosInfoPredictor) *BagagePosInfo {
	if len(bl) <= 0 {
		return nil
	}
	if p(bl[0]) {
		return bl[0]
	}
	return bl[1:].findOne(p)
}
func (bl BagagePosInfoList) find(p BagagePosInfoPredictor) BagagePosInfoList {
	return bl.findRecursive(p, BagagePosInfoList{})
}
func (bl BagagePosInfoList) findRecursive(p BagagePosInfoPredictor, l BagagePosInfoList) BagagePosInfoList {
	if len(bl) <= 0 {
		return l
	}
	if p(bl[0]) {
		l = append(l, bl[0])
	}
	return bl[1:].findRecursive(p, l)
}
func (this BagagePosInfoList) UsingImage(imageName string) bool {
	for _, bpi := range this {
		if bpi.ImageName == imageName {
			return true
		}
	}
	return false
}

// func (this BagagePosInfoList) Find(bagageID string) *BagagePosInfo {
// 	for _, bi := range this {
// 		if bi.BagageID == bagageID {
// 			return bi
// 		}
// 	}
// 	return nil
// }

// //是否存在单号和坐标完全一致的订单信息
// func (this BagagePosInfoList) BagageInfoRepeat(bpi *BagagePosInfo) bool {
// 	bi := this.Find(bpi.BagageID)
// 	if bi == nil {
// 		return false
// 	} else if bi.Latitude == bpi.Latitude && bi.Longitude == bpi.Longitude {
// 		return true
// 	}
// 	return false
// }
