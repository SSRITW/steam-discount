package monitor

import (
	"net/http"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
)


type Game struct {
	Id int
	Name string
	PayUrl string
	Price float64
	issueDate string
	SupportPlatforms []string
}

type MonitorCentent struct {
	Game
	Off string
	AtferOffPrice float64
}

//获取网页内容
func GetContent(url string) (io.Reader,bool){
	resp,err := http.Get(url)
	if err!=nil {
		fmt.Println("httpGetError:",err.Error())
	}
	if resp.StatusCode==200 {
		body := resp.Body
		defer body.Close()
		return body,true
	}
	return nil,false
}


func FinTargetedContent(content io.Reader,outsideSelector string,contentIndex int,) ( []string, error){
	reader ,err := goquery.NewDocumentFromReader(content)
	if err!=nil {
		return nil,err
	}
	reader.Find(outsideSelector,)
	return targetContents
}


