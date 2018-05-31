package monitor

import (
	"testing"
	"strconv"
	"fmt"
)

//
func Test_GetContent(t *testing.T){
	var contentChannel = make(chan MonitorContent,50)
	var pageSizeChannel = make(chan int )
	var maxContentSize = make(chan int)

	go readPageSizeChannel(pageSizeChannel,contentChannel)

	go GetContent("https://store.steampowered.com/search/?specials=1",contentChannel,pageSizeChannel,maxContentSize)

	readContentChannel(contentChannel,maxContentSize)

}

//当pageChan有值时才进行读操作（即，进行第一次爬虫获取到这个值的时候；启动多个goroutine继续爬取之后的页面数据
func readPageSizeChannel(pageChan chan int,contentChannel chan MonitorContent){
	maxI,ok := <- pageChan
	if ok {
		for i:=2; i<=maxI ;i++  {
			go GetContent("https://store.steampowered.com/search/?specials=1&page="+strconv.Itoa(i),contentChannel,nil,nil)
		}
		close(pageChan)
	}
}

//当maxContentSize有值时才进行读操作（即，进行第一次爬虫获取到这个值的时候
func readContentChannel(contentChannel chan MonitorContent,maxContentSize chan int){
	i := 0
	var size int
	if v,f :=<-maxContentSize;f{
		size = v
		fmt.Println("contentSize",size)
		close(maxContentSize)
	}
	for v := range contentChannel{
		i++
		fmt.Println(v.Thumbnail,i)
		if size == i {
			close(contentChannel)
		}
	}
}
