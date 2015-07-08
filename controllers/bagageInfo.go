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

var (
// url_bagageInfoList = "http://111.67.197.251:9002/bagageList"
)

//获取所有订单车辆绑定信息
// func requestBagageInfoList() {
// 	resp, body, errs := gorequest.New().Post(url_bagageInfoList).Send("").End()
// 	if errs != nil {
// 		for _, _err := range errs {
// 			DebugSysF("请求快递列表信息时出错：%s", _err.Error())
// 		}
// 		return
// 	}
// 	DebugTraceF("请求快递列表 信息结果状态：%s", resp.Status)
// 	DebugTraceF(body)

// 	var list bagageInfoList
// 	if err := json.Unmarshal([]byte(body), &list); err != nil {
// 		DebugSysF(err.Error())
// 		return
// 	} else {
// 		DebugPrintList_Trace(list)
// 		G_bagageInfos = list
// 	}
// }

//订单与车辆的关系，订单的位置需要根据与其绑定的车辆的位置确定
type bagageInfo struct {
	BagageID  string
	TimeStamp string
	CarID     string
}

func (this *bagageInfo) String() string {
	return fmt.Sprintf("单号:%10s    车辆编号：%10s   绑定时间: %s", this.BagageID, this.CarID, this.TimeStamp)
}

type bagageInfoList []*bagageInfo

func (this bagageInfoList) FindBagage(bagageID string) *bagageInfo {
	for _, bi := range this {
		if bi.BagageID == bagageID {
			return bi
		}
	}
	return nil
}

func (this bagageInfoList) ListName() string {
	return "订单信息列表："
}
func (this bagageInfoList) InfoList() []string {
	list := []string{}
	for _, bi := range this {
		list = append(list, bi.String())
	}
	return list
}
