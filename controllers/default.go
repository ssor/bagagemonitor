package controllers

import (
	"encoding/json"
	"encoding/xml"
	"github.com/astaxie/beego/config"
	// "errors"
	"crypto/sha1"
	"fmt"
	"github.com/astaxie/beego"
	// "github.com/parnurzeal/gorequest"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	// "strconv"
	"strings"
	"time"
)

/*
	编译命令：
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build main.go
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
)

var (
	token                               = "nodewebgis" //微信接口
	localhost                           = "http://localhost/"
	imageDirPath                        = localhost + "images/"
	G_iniconf    config.ConfigContainer = nil
	g_bagages                           = BagagePosInfoList{}
	// bagageStatusUrl                                  = "http://111.67.197.251:9002/getBagageStatus" //post，获取单号信息，获取的是与该单号绑定的车的位置信息
	// G_bagageInfos                                    = bagageInfoList{}
	// G_CarMapImageInfoList                            = CarMapImageInfoList{}
	// G_MapImageRefreshInterval                        = DEFAULT_REFRESH_INTERVAL //刷新位置
)

func init() {
	// DebugPrintList_Trace(g_bagages)
	initConfig()
	go startIntervalCheck(5 * time.Second)
	// g_bagages = append(g_bagages, NewBagagePosInfo("b001", "12:12:12", "", "111", "222"))
	// requestBagageInfoList()
	// refreshBagagePosImage()

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
			// requestBagageInfoList()
			// refreshBagagePosImage()
			removeExpiredImage()
		}
	}
}

//移除过期不用的图片
func removeExpiredImage() {
	ImageDirPath := "static/img/"
	fileNameList := []string{}
	walkFn := func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") == true {
			return nil
		}
		if info.Name() == "error.png" || info.Name() == "rt.png" {
			return nil
		}
		if info.IsDir() == false {
			fileNameList = append(fileNameList, info.Name())
		}
		return nil
	}
	if err := filepath.Walk(ImageDirPath, walkFn); err != nil {
		DebugSysF(err.Error())
		return
	}
	// sort.Strings(fileNameList)
	// fmt.Println(fileNameList)
	for _, name := range fileNameList {
		if g_bagages.UsingImage(name) == false {
			if err := os.Remove(ImageDirPath + name); err != nil {
				DebugSysF("删除过期图片时出错：%s", err.Error())
			} else {
				DebugTraceF("删除过期图片：%s", name)
			}
		}
	}
}

// //根据车辆绑定信息，获取订单的位置信息
// func refreshBagagePosImage() {
// 	for _, bi := range G_bagageInfos {
// 		time.Sleep(1 * time.Second)
// 		go bagageInfoRequest(bi.BagageID)
// 	}
// }

// //获取指定订单的位置状态信息，并转化成地图
// func bagageInfoRequest(bagageID string) {
// 	content := fmt.Sprintf(`{"bagageID":"%s"}`, bagageID)
// 	resp, body, errs := gorequest.New().Post(bagageStatusUrl).Send(content).End()
// 	if errs != nil {
// 		for _, _err := range errs {
// 			DebugSysF("请求单号信息时出错：%s", _err.Error())
// 		}
// 		return
// 	}
// 	DebugTraceF("查询单号 %s 状态 结果状态：%s", bagageID, resp.Status)
// 	DebugTraceF(body)
// 	if len(strings.TrimSpace(string(body))) <= 0 {
// 		return
// 	}
// 	var bpi BagagePosInfo
// 	if err := json.Unmarshal([]byte(body), &bpi); err != nil {
// 		DebugSysF(err.Error())
// 		return
// 	} else {
// 		(&bpi).BagageID = bagageID
// 		DebugTrace(bpi.String())
// 		if imageName, err := G_CarMapImageInfoList.HasImageTemp(bagageID, bpi.SogouLongitude, bpi.SogouLatitude); err != nil {
// 			downloadMap(&bpi)
// 		} else {
// 			DebugTraceF("快递 %s 地图位置有缓存 %s", bagageID, imageName)
// 		}
// 	}
// }

//下载地图
// 搜狗地图API：http://api.go2map.com/engine/api/static/image+{'points':'116.36620044708252,39.96220463653672',height:'450','width':550,'zoom':9,'center':'116.36620044708252,39.96220463653672',labels:'搜狐网络大厦',pss:'S1756',city:'北京'}.png
func downloadMap(bpi *BagagePosInfo) (downloadImageName string, result bool) {
	uid := time.Now().UnixNano()
	imageName := fmt.Sprintf("%s_%d.png", bpi.Flag, uid) //使用车辆编号，当不同的快递在同一辆车上时可以复用
	DebugTraceF("保存的图片名称：%s", imageName)
	url := fmt.Sprintf("http://api.go2map.com/engine/api/static/image+{'points':'%s,%s',height:'341','width':512,'zoom':11,'center':'%s,%s',labels:'%s',pss:'S1756'}.png",
		bpi.SogouLongitude, bpi.SogouLatitude, bpi.SogouLongitude, bpi.SogouLatitude, "")
	DebugTraceF("获取快递最新位置地图链接：%s", url)
	if err := DownloadFromUrl(url, "static/img/"+imageName); err != nil {
		DebugMustF("下载 %s 的位置地图出错：%s", bpi.BagageID, err)
		return "", false
	} else {
		return imageName, true
	}
}
func DownloadFromUrl(url, fileName string) error {
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

// func DownloadFromUrl(filePath, url string, chanDownloadCount chan bool) {
// 	// fileTempPath := "./tmp/"
// 	//如果路径中包含文件夹，需要首先建立该文件夹
// 	// fileName := path.Base(url)
// 	// file, err := os.OpenFile(fileTempPath+fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
// 	if err != nil {
// 		DebugSysF(err.Error())
// 		chanDownloadCount <- false
// 		return
// 	}
// 	fileDl, err := downloader.NewFileDl(url, file, -1)
// 	if err != nil {
// 		DebugSysF("下载 [%s] 出错：%s  url: %s", filePath, err.Error(), url)
// 		os.Remove(filePath)
// 		chanDownloadCount <- false
// 		return
// 	}
// 	var chanExit = make(chan bool)
// 	var chanProgress = make(chan bool)
// 	var chanAbort = make(chan bool)
// 	fileDl.OnStart(func() {
// 		// fmt.Println("开始下载")
// 		for {
// 			select {
// 			case <-chanExit:
// 				status := fileDl.GetStatus()
// 				// fmt.Println(fmt.Sprintf(format, status.Downloaded, fileDl.Size, h, 0, "[FINISH]"))
// 				DebugTraceF("[%s] 下载完成，共 %d 字节", filePath, status.Downloaded)
// 				// DebugTrace("关闭文件"+GetFileLocation())
// 				file.Close()
// 				chanDownloadCount <- true
// 				return
// 			case <-chanAbort:
// 				i := 0
// 				for {
// 					if err := file.Close(); err == nil {
// 						DebugInfo("下载取消成功，关闭了文件 [" + filePath + "]" + GetFileLocation())
// 						break
// 					}
// 					time.Sleep(time.Second * 1)
// 					i++
// 					if i > 3 {
// 						DebugMust("下载取消失败，无法关闭文件 [" + filePath + "]" + GetFileLocation())
// 						break
// 					}
// 				}
// 				chanDownloadCount <- false
// 				return
// 			case <-chanProgress:
// 				// format := "\033[2K\r%v/%v [%s] (当前速度： %v byte/s) %v"
// 				// status := fileDl.GetStatus()
// 				// var i = float64(status.Downloaded) / float64(fileDl.Size) * 50
// 				// if i < 0 {
// 				// 	i = 0
// 				// }
// 				// h := strings.Repeat("=", int(i)) + strings.Repeat(" ", 50-int(i))

// 				// fmt.Println(fmt.Sprintf(format, status.Downloaded, fileDl.Size, h, status.Speeds, "[DOWNLOADING]"))
// 			}
// 		}
// 	})
// 	fileDl.OnAbort(func() {
// 		chanAbort <- true
// 	})
// 	fileDl.OnProgress(func() {
// 		chanProgress <- true
// 	})
// 	fileDl.OnFinish(func() {
// 		chanExit <- true
// 	})

// 	fileDl.OnError(func(errCode int, err error) {
// 		fmt.Println(errCode, err)
// 	})

// 	DebugTraceF("开始下载 url: %s", url)
// 	fileDl.Start()
// 	// return fileDl, nil
// }

type MainController struct {
	beego.Controller
}

//接收订单查询请求，返回地图信息
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
		bagagePosInfo := g_bagages.Find(bagageID)
		if bagagePosInfo == nil {
			//构造一个没有单号的返回信息
			imageUrl := imageDirPath + "error.png"
			articleList := ArticleList{&Article{NewArticleItem("订单状态查询", "没有该订单信息", imageUrl, imageUrl)}}
			weixinRes = NewWeixinResponseNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articleList)
			return
		}
		if bagagePosInfo.ImageName == "" {
			// 使用没有地图的默认图片构造返回信息
			imageUrl := imageDirPath + "rt.png"
			articleList := ArticleList{&Article{NewArticleItem("订单状态查询", "暂未找到订单的位置", imageUrl, imageUrl)}}
			weixinRes = NewWeixinResponseNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articleList)
			return
		}
		// mapImageInfo := G_CarMapImageInfoList.Find(bagageInfo.CarID)
		// if mapImageInfo == nil || mapImageInfo.ImageName == "" { //有单号没图片
		// 	// 使用没有地图的默认图片构造返回信息
		// 	imageUrl := imageDirPath + "rt.png"
		// 	articleList := ArticleList{&Article{NewArticleItem("订单状态查询", "暂未找到订单的位置", imageUrl, imageUrl)}}
		// 	weixinRes = NewWeixinResponseNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articleList)
		// 	return
		// }
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

//添加新的订单位置信息后，直接开始下载对应的图片
//下载成功，加入到总的订单列表
//下载失败，等待N秒后，重新开始
func downloadBagageImage(bagages BagagePosInfoList) {
	f := func(bpi *BagagePosInfo, interval time.Duration) {
		for {
			DebugTraceF("准备下载地图 %s", bpi.BagageID)
			if imageName, result := downloadMap(bpi); result == true {
				bpi.ImageName = imageName
				bi := g_bagages.Find(bpi.BagageID)
				if bi == nil {
					g_bagages = append(g_bagages, bpi)
				} else {
					bi.update(bpi.TimeStamp, bpi.SogouLongitude, bpi.SogouLatitude, bpi.ImageName, bpi.Flag)
				}
				return
			} else {
				DebugInfoF("下载订单 %s 位置地图出错，%d 秒后重试", bpi.BagageID, interval)
				time.Sleep(interval)
			}
		}
	}
	for _, b := range bagages {
		if g_bagages.BagageInfoRepeat(b) == false { //重复不需要下载
			go f(b, 15*time.Second) //下载失败15秒后重试
		}
	}
}
func (this *MainController) AddBagage() {
	body := this.Ctx.Input.CopyBody()
	list := BagagePosInfoList{}
	err := json.Unmarshal(body, &list)
	if err != nil {
		DebugMustF("解析传入数据有误：%s", err)
		this.CustomAbort(http.StatusBadRequest, "数据格式有误")
	} else {
		if len(list) > 0 {
			DebugInfoF("新添加了 %d 个订单", len(list))
			DebugPrintList_Trace(list)
			downloadBagageImage(list)
		}
		// this.CustomAbort(http.StatusOK, "")
		this.ServeJson()
	}
}
func (this *MainController) BagageList() {
	this.Data["json"] = g_bagages
	this.ServeJson()
}
func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplNames = "index.tpl"
}
