package controllers

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/config"
	// "errors"
	"fmt"
	"github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	"io"
	"net/http"
	"os"
	// "path/filepath"
	// "strconv"
	"github.com/ungerik/go-dry"
	"time"
)

/*
	编译命令：
	GOOS=windows GOARCH=386 go build main.go
*/
/*
	运作原理：
	* 本系统开放接口，外部系统可以通过接口添加需要监控的订单及其坐标位置（目前只支持搜狗坐标）
	* 只在更新坐标位置时更新位置图片

	需要考虑的大量数据的问题：
	1. 海量位置信息及其它信息的管理
	2. 海量位置图片的产生问题
	3. 图片文件的管理问题
	4. 图片的静态服务
	5. 大量订单的状态刷新问题，中间件
	6. 位置刷新的时间间隔的设定，应该根据接近目的地而间隔变小
*/
const (
	DEFAULT_REFRESH_INTERVAL = 30 * time.Second
	ImageDirPath             = "static/img/"
	ImageDirDefaultPath      = "static/img/defaultImages/"
)

var (
	token                               = "nodewebgis" //微信接口
	localhost                           = "http://localhost/"
	imageDirPath                        = localhost + "images/"
	G_iniconf    config.ConfigContainer = nil
	g_bagages                           = BagagePosInfoList{}
	// g_imageNamesExpired                        = ImageInfoList{} //注册为过期的图片，可以删除
)

func init() {
	initConfig()
	go startIntervalCheck(5 * time.Hour)
}
func initConfig() {
	var err error
	G_iniconf, err = config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		beego.Error(err.Error())
	} else {
		localip := G_iniconf.String("localip")
		if len(localip) <= 0 {
			localip = "localhost"
		}
		localhost = fmt.Sprintf("http://%s/", localip)
		DebugInfoF("本地网络：%s", localhost)
		imageDirPath = localhost + "images/"
	}
}
func startIntervalCheck(interval time.Duration) {
	ticker := time.Tick(interval)
	for {
		select {
		case <-ticker:
			removeExpiredImage()
		}
	}
}

type MainController struct {
	beego.Controller
}

//add bagage pos info
func (m *MainController) AddBagage() {
	responseHandler(m, func(m *MainController) (interface{}, error) {
		body := m.Ctx.Input.CopyBody()
		list := BagagePosInfoList{}
		err := json.Unmarshal(body, &list)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("解析传入数据有误：%s", err))
			// DebugMustF("解析传入数据有误：%s", err)
			// m.CustomAbort(http.StatusBadRequest, "数据格式有误")
		} else {
			if len(list) > 0 {
				DebugInfoF("新添加了 %d 个订单", len(list))
				DebugPrintList_Trace(list)
				// downloadBagageImage(list)
				list.forEach(func(b *BagagePosInfo) {
					if g_bagages.findOne(func(bi *BagagePosInfo) bool { return bi.equals(b) }) == nil {
						go addBagagePosInfo(b, 15*time.Second) //下载失败15秒后重试
					}
				})
			}
			return nil, nil
			// m.ServeJson()
		}
	})
}

//get all bagage
func (this *MainController) BagageList() {
	this.Data["json"] = g_bagages
	this.ServeJson()
}

// //添加新的订单位置信息后，直接开始下载对应的图片
// //下载成功，加入到总的订单列表
// //下载失败，等待N秒后，重新开始
// func downloadBagageImage(bagages BagagePosInfoList) {
// 	bagages.forEach(func(b *BagagePosInfo) {
// 		if g_bagages.findOne(func(bi *BagagePosInfo) bool { return bi.equals(b) }) == nil {
// 			go addBagagePosInfo(b, 15*time.Second) //下载失败15秒后重试
// 		}
// 	})
// }

func addBagagePosInfo(bpi *BagagePosInfo, interval time.Duration) {
	for {
		DebugTraceF("准备下载地图 %s", bpi.BagageID)
		if imageName, result := downloadMap(bpi); result == true {
			bpi.ImageName = imageName
			bi := g_bagages.findOne(func(b *BagagePosInfo) bool { return b.BagageID == bpi.BagageID })
			if bi == nil {
				g_bagages = append(g_bagages, bpi)
			} else {
				// g_imageNamesExpired = g_imageNamesExpired.RegisterImage(bi.ImageName) //将之前使用的图片注册到可删除列表
				bi.update(bpi.TimeStamp, bpi.Longitude, bpi.Latitude, bpi.ImageName, bpi.Flag)
			}
			return
		} else {
			DebugInfoF("下载订单 %s 位置地图出错，%d 秒后重试", bpi.BagageID, interval)
			time.Sleep(interval)
		}
	}
}

//移除过期不用的图片
func removeExpiredImage() {
	if images, err := dry.ListDirFiles(ImageDirPath); err != nil {
		DebugMustF("removeExpiredImage error: %s", err)
	} else {
		// images = filterString(images, func(s string) { return inStringList([]string{"error.png", "rt.png"}, s) })
		for _, image := range images {
			if g_bagages.findOne(func(b *BagagePosInfo) bool { return b.ImageName == image }) == nil {
				if e := os.Remove(ImageDirPath + image); e != nil {
					DebugSysF("remove image error: %s", e)
				}
			}
		}
	}

	// ImageDirPath := "static/img/"
	// for _, ii := range g_imageNamesExpired {
	// 	if err := os.Remove(ImageDirPath + ii.Name); err != nil {
	// 		DebugSysF("删除过期图片时出错：%s", err.Error())
	// 	} else {
	// 		DebugTraceF("删除过期图片：%s", ii.Name)
	// 		ii.SetDeleted()
	// 	}
	// }
	// g_imageNamesExpired = g_imageNamesExpired.Clear()
	// return
}

// func inStringList(src []string, t string) bool {
// 	for _, s := range src {
// 		if s == t {
// 			return true
// 		}
// 	}
// 	return false
// }
// func filterString(src []string, p func(string) bool) []string {
// 	l := []string{}
// 	for _, s := range src {
// 		if p(s) {
// 			l = append(l, s)
// 		}
// 	}
// 	return l
// }

//下载地图
// 百度地图API：http://api.map.baidu.com/staticimage?center=116.403874,39.914888&width=300&height=200&zoom=11
// 搜狗地图API：http://api.go2map.com/engine/api/static/image+{'points':'116.36620044708252,39.96220463653672',height:'450','width':550,'zoom':9,'center':'116.36620044708252,39.96220463653672',labels:'搜狐网络大厦',pss:'S1756',city:'北京'}.png
func downloadMap(bpi *BagagePosInfo) (downloadImageName string, result bool) {
	uid := time.Now().UnixNano()
	imageName := fmt.Sprintf("%s_%d.png", bpi.Flag, uid) //使用车辆编号，当不同的快递在同一辆车上时可以复用
	DebugTraceF("保存的图片名称：%s", imageName)
	url := fmt.Sprintf("http://api.map.baidu.com/staticimage?width=512&height=341&center=%s,%s&zoom=13&markers=%s,%s&markerStyles=l,%s",
		bpi.Longitude, bpi.Latitude, bpi.Longitude, bpi.Latitude, "")
	// url := fmt.Sprintf("http://api.go2map.com/engine/api/static/image+{'points':'%s,%s',height:'341','width':512,'zoom':11,'center':'%s,%s',labels:'%s',pss:'S1756'}.png",
	// 	bpi.Longitude, bpi.Latitude, bpi.Longitude, bpi.Latitude, "")
	DebugTraceF("获取快递最新位置地图链接：%s", url)
	if err := downloadFromUrl(url, "static/img/"+imageName); err != nil {
		DebugMustF("下载 %s 的位置地图出错：%s", bpi.BagageID, err)
		return "", false
	} else {
		return imageName, true
	}
}
func downloadFromUrl(url, fileName string) error {
	rawURL := url
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()

	//fmt.Printf("Downloading file %s...", fileName)
	//fmt.Println()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(rawURL) // add a filter to check redirect

	if err != nil {
		// fmt.Println(err)
		// panic(err)
		return err
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		return err
	}
	DebugTraceF("下载完成 (%d bytes) %s", size, fileName)
	return nil
	// fmt.Printf("%s,%s,%v", resp.Status, rawURL, size)
	// fmt.Println()
}

type logicHandler func(m *MainController) (interface{}, error)

func responseHandler(m *MainController, handler logicHandler) {
	response := NewResponseMsg(0)
	defer func() {
		m.Data["json"] = response
		m.ServeJson()
	}()
	if value, err := handler(m); err != nil {
		DebugMustF("error: %s", err.Error())
		response = NewResponseMsg(1, err.Error())
	} else {
		response.Data = value
	}
}

type ResponseMsg struct {
	Code    int
	Message string
	Data    interface{}
}

func NewResponseMsg(code int, msg ...string) *ResponseMsg {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return &ResponseMsg{
		Code:    code,
		Message: message,
	}
}
