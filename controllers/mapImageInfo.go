package controllers

import (
	// "encoding/json"
	"errors"
	// "fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	// "os"
	"time"
)

var G_MapImageRefreshInterval = 30 * time.Second

type CarMapImageInfo struct {
	CarID     string
	ImageName string
	UID       int64
}

func NewCarMapImageInfo(carID, imageName string, uid int64) *CarMapImageInfo {
	return &CarMapImageInfo{
		CarID:     carID,
		ImageName: imageName,
		UID:       uid,
	}
}

type CarMapImageInfoList []*CarMapImageInfo

func (this CarMapImageInfoList) Find(carID string) *CarMapImageInfo {
	for _, mii := range this {
		if mii.CarID == carID {
			return mii
		}
	}
	return nil
}

func (this CarMapImageInfoList) HasImageTemp(carID string) (string, error) {
	mii := this.Find(carID)
	if mii != nil {
		if (time.Now().UnixNano() - mii.UID) < int64(G_MapImageRefreshInterval) {
			return mii.ImageName, nil
		}
	}
	return "", errors.New("没有找到")
}
func (this CarMapImageInfoList) Add(_mii *CarMapImageInfo) CarMapImageInfoList {
	mii := this.Find(_mii.CarID)
	if mii != nil {
		mii.ImageName = _mii.ImageName
		mii.UID = _mii.UID
		return this
	} else {
		return append(this, _mii)
	}
}
func (this CarMapImageInfoList) UsingImage(imageName string) bool {
	for _, mii := range this {
		if mii.ImageName == imageName {
			return true
		}
	}
	return false
}

// type mapImageInfo struct {
// 	BagageID  string
// 	ImageName string
// 	UID       int64
// }

// func NewMapImageInfo(bagageID, imageName string, uid int64) *mapImageInfo {
// 	return &mapImageInfo{
// 		BagageID:  bagageID,
// 		ImageName: imageName,
// 		UID:       uid,
// 	}
// }

// type mapImageInfoList []*mapImageInfo

// func (this mapImageInfoList) Find(bagageID string) *mapImageInfo {
// 	for _, mii := range this {
// 		if mii.BagageID == bagageID {
// 			return mii
// 		}
// 	}
// 	return nil
// }

// func (this mapImageInfoList) HasImageTemp(bagageID string) (string, error) {
// 	mii := this.Find(bagageID)
// 	if mii != nil {
// 		if (time.Now().UnixNano() - mii.UID) < int64(10*time.Second) {
// 			return mii.ImageName, nil
// 		}
// 	}
// 	return "", errors.New("没有找到")
// }
// func (this mapImageInfoList) Add(_mii *mapImageInfo) mapImageInfoList {
// 	mii := this.Find(_mii.BagageID)
// 	if mii != nil {
// 		mii.ImageName = _mii.ImageName
// 		mii.UID = _mii.UID
// 		return this
// 	} else {
// 		return append(this, _mii)
// 	}
// }
