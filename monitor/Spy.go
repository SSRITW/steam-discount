package monitor

import (
	"net/http"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
	"time"
	"io/ioutil"
	"math/rand"
	"strconv"
)


const (
	WINDOWS = iota+1
	MAC
	LINUX
)

//id可能会存在的标签
var idAttrNames = [...] string{"data-ds-packageid","data-ds-appid","data-ds-bundleid"}

//代理
var userAgent = [...]string{"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"}

type Game struct {
	Id int `redis:"id"`
	Name string	`redis:"name"`
	PayUrl string `redis:"payUrl"`
	Thumbnail string `redis:"thumbnail"`
	Price float64 `redis:"price"`
	IssueDate string `redis:"issueDate"`
	SupportPlatforms []int `redis:"dupportPlatforms"`
}

type MonitorContent struct {
	Game
	Off string `redis:"off"`
	AtferOffPrice float64 `redis:"atferOffPrice"`
}

//获取网页内容
func GetContent(url string , contents chan MonitorContent,pageSize chan int,maxContentSize chan int){
	defer func() {	//用于recover可能会出现的panic
		if r := recover(); r != nil {
			fmt.Println("getContentPanic", r)
		}
	}()
	//fmt.Println("now is ",url)
	request,_ := http.NewRequest("GET",url,nil)
	request.Header.Set("User-Agent", getRandomUserAgent())
	client := http.DefaultClient
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("Get请求返回错误:", url, err.Error())
		return
	}
	if res.StatusCode==200 {
		body := res.Body
		finTargetedContent(body,contents,pageSize,maxContentSize)
		defer body.Close()
	}
}

//查找目标内容并放入MonitorContent内
func finTargetedContent(content io.Reader,contents chan MonitorContent,pageSize chan int,maxContentSize chan int){
	reader ,err := goquery.NewDocumentFromReader(content)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}


	if pageSize!=nil {
		pSize,contentSize := getPageSizeAndMaxContentSize(reader)
		pageSize <- pSize
		maxContentSize <- contentSize
	}

	reader.Find("#search_result_container div").Eq(1).Find("a").Each(func(i int , tag *goquery.Selection){
		content := MonitorContent{}

		var id string
		var exists bool
		for i:=0 ; i<len(idAttrNames) && !exists ; i++ {
			id,exists = tag.Attr(idAttrNames[i])
		}

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

		content.Id,_ = strconv.Atoi(id)
		content.Name = name
		content.PayUrl = payUrl
		content.Thumbnail = downPicture(thumbnail,id)
		content.IssueDate = formatDate(issueDate)
		//去除分行符和空格、‘-’和‘%’符号
		content.Off = strings.TrimSpace(off)
		if content.Off!="" {
			content.Off = content.Off[1:len(content.Off)-1]
		}

		content.Price,content.AtferOffPrice = separatePrices(prices)
		content.SupportPlatforms = supportPlatforms

		contents <- content
	})
	return
}

//通过分隔包含打折前后价格的字符串，得到打折前后价格
func separatePrices(prices string)(price float64 , offPrice float64){
	prices = strings.TrimSpace(prices)
	pricesSlice := strings.Split(prices,"¥")
	size := len(pricesSlice)
	if size==3 {
		price,_ = strconv.ParseFloat(strings.TrimSpace(pricesSlice[1]),32)
		offPrice,_ = strconv.ParseFloat(strings.TrimSpace(pricesSlice[2]),32)
	}
	return
}

//通过格式化日期字符串为time，再由time格式化为想要的格式的string
func formatDate(dateStr string) (formatStr string){
	if dateStr=="" {
		return
	}
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
	request,_ := http.NewRequest("GET",pictureUrl,nil)
	request.Header.Set("User-Agent", getRandomUserAgent())
	client := http.DefaultClient
	res, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Get请求返回错误:", pictureUrl, err.Error())
		return
	}
	if res.StatusCode==200 {
		body := res.Body
		contents,_ := ioutil.ReadAll(body)
		localUrl = "/resource/images/"+gameIdStr+".jpg"
		err := ioutil.WriteFile("."+localUrl, contents, 0666) //写入文件(字节数组)
		if err!=nil {
			fmt.Println("downPicture:",err.Error())
		}
		defer body.Close()
	}
	return
}

//利用随机数，随机获取一个代理
func getRandomUserAgent() string {
	return userAgent[rand.Intn(len(userAgent))]
}

//获取最大页码和内容总条数
func getPageSizeAndMaxContentSize(reader *goquery.Document)(pageSize int ,maxContentSize int){
	contentSizeStr := reader.Find("#search_result_container .search_pagination .search_pagination_left").Text()

	pageButtonSelection := reader.Find("#search_result_container .search_pagination .search_pagination_right a")

	//倒数第二个a标签是最后一页的页码
	pageStr := pageButtonSelection.Eq(pageButtonSelection.Length()-2).Text()

	pageSize,_ = strconv.Atoi(pageStr)

	contentSizeStr = strings.TrimSpace(contentSizeStr)
	//去除空格之后`showing 1 - 25 of 1434`，以`of`为分隔符拆分字符串，取到` 1434 `,最后去除空格
	contentSizeStr = strings.TrimSpace(strings.Split(contentSizeStr,"of")[1])

	maxContentSize ,_ = strconv.Atoi(contentSizeStr)
	return
}

