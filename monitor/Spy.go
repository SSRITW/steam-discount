package monitor

import (
	"net/http"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"steam-discount/common"
	"strings"
	"time"
	"io/ioutil"
)


const (
	WINDOWS = iota+1
	MAC
	LINUX
)

type Game struct {
	Id int
	Name string
	PayUrl string
	Thumbnail string
	Price float64
	IssueDate string
	SupportPlatforms []int
}

type MonitorContent struct {
	Game
	Off string
	AtferOffPrice float64
}

//获取网页内容
func GetContent(url string) ( contents []MonitorContent){
	resp,err := http.Get(url)
	if err!=nil {
		fmt.Println("httpGetError:",err.Error())
	}
	if resp.StatusCode==200 {
		body := resp.Body
		contents = finTargetedContent(body)
		defer body.Close()
	}
	return contents
}

//查找目标内容并放入MonitorContent内
func finTargetedContent(content io.Reader) ( contents []MonitorContent){
	reader ,err := goquery.NewDocumentFromReader(content)
	if err!=nil {
		fmt.Println(err.Error())
		return nil
	}

	reader.Find("#search_result_container div").Eq(1).Find("a").Each(func(i int , tag *goquery.Selection){
		content := MonitorContent{}
		id,_ := tag.Attr("data-ds-appid")
		name := tag.Find(".search_name .title").Text()
		payUrl,_ := tag.Attr("href")
		thumbnail,_ := tag.Find(".search_capsule img").Attr("src")
		issueDate := tag.Find(".search_released").Text()
		off := tag.Find(".search_discount").Text()
		prices := tag.Find(".search_price").Text()
		var supportPlatforms []int;
		tag.Find(".search_name p span").Each(func(i int, selection *goquery.Selection) {//通过样式（class）找出游戏支撑的平台
			hadClass, _:= selection.Attr("class")
			if strings.Contains(hadClass,"win") {
				supportPlatforms = append(supportPlatforms,WINDOWS)
			}else if(strings.Contains(hadClass,"mac")){
				supportPlatforms = append(supportPlatforms,MAC)
			}else if(strings.Contains(hadClass,"linux")){
				supportPlatforms = append(supportPlatforms,LINUX)
			}
		})

		content.Id = common.StrToInt(id)
		content.Name = name
		content.PayUrl = payUrl
		content.Thumbnail = downPicture(thumbnail,id)
		content.IssueDate = formatDate(issueDate)
		content.Off = strings.TrimSpace(off)	//去除分行符和空格
		content.Price,content.AtferOffPrice = separatePrices(prices)
		content.SupportPlatforms = supportPlatforms

		contents = append(contents,content)
	})
	return
}

//通过分隔包含打折前后价格的字符串，得到打折前后价格
func separatePrices(prices string)(price float64 , offPrice float64){
	prices = strings.TrimSpace(prices)
	pricesSlice := strings.Split(prices,"¥")
	size := len(pricesSlice)
	if size==3 {
		price = common.StrToFloat64(strings.TrimSpace(pricesSlice[1]))
		offPrice = common.StrToFloat64(strings.TrimSpace(pricesSlice[2]))
	}
	return
}

//通过格式化日期字符串为time，再由time格式化为想要的格式的string
func formatDate(dateStr string) (formatStr string){
	v,e := time.Parse("2 Jan, 2006",dateStr)
	if e!=nil {
		fmt.Println(e.Error())
	}else{
		formatStr = v.Format("2006-01-02")
	}
	return
}

//将游戏的图标保存到本地
func downPicture(pictureUrl string,gameIdStr string)(localUrl string){
	resp,err := http.Get(pictureUrl)
	if err!=nil {
		fmt.Println("downPicture:",err.Error())
	}
	if resp.StatusCode==200 {
		body := resp.Body
		contents,_ := ioutil.ReadAll(body)
		localUrl = "/resource/images/"+gameIdStr+".jpg"
		err := ioutil.WriteFile(".."+localUrl, contents, 0666) //写入文件(字节数组)
		if err!=nil {
			fmt.Println("downPicture:",err.Error())
		}
		defer body.Close()
	}
	return
}

