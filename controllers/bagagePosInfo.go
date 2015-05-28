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
type bagagePosInfo struct {
	CarID          string
	TimeStamp      string
	SogouLongitude string
	SogouLatitude  string
	BagageID       string
}

func (this *bagagePosInfo) String() string {
	return fmt.Sprintf("CarID: %10s  BagageID:%10s   time: %10s     Position:(%10s, %10s)", this.CarID, this.BagageID, this.TimeStamp, this.SogouLongitude, this.SogouLatitude)
}
