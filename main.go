package main

import (
	"steam-discount/monitor"
	"time"
	"steam-discount/service"
	"fmt"
)

func main(){
	initData()

	contents := service.GetGameContents()	//获取存在redis的全部gameContent
	for _,data := range contents{
		fmt.Println(data)
	}

	content :=service.GetGameContentById("1000")		//通过id取值
	fmt.Println("contentById:",content)
}

func initData(){
	var contentChannel = make(chan monitor.MonitorContent,50)
	var pageSizeChannel = make(chan int )
	var maxContentSize = make(chan int)

	go monitor.GetContentsByPageChan(pageSizeChannel,contentChannel)

	go monitor.GetContent("https://store.steampowered.com/search/?specials=1",contentChannel,pageSizeChannel,maxContentSize)

	monitor.SaveContents(contentChannel,maxContentSize,time.Second*30)
}
