package main

import (
	"steam-discount/monitor"
	"time"
	"fmt"
	"github.com/robfig/cron"
)

func main(){
	initData()

	c := cron.New()
	c.AddFunc("* */10 * * * *", initData) //十分钟爬一次数据
	c.Start()

	for i:=0;i>=0;{		//使主线程不会停止
		i = 1
	}

	/*
	contents := service.GetGameContents()	//获取存在redis的全部gameContent
	for _,data := range contents{
		fmt.Println(data)
	}
	content := service.GetGameContentById("448500")	//获取id数据
	*/

}

func initData(){
	var contentChannel = make(chan monitor.MonitorContent,50)
	var pageSizeChannel = make(chan int )
	var maxContentSize = make(chan int)

	go monitor.GetContentsByPageChan(pageSizeChannel,contentChannel)

	go monitor.GetContent("https://store.steampowered.com/search/?specials=1",contentChannel,pageSizeChannel,maxContentSize)

	monitor.SaveContents(contentChannel,maxContentSize,time.Second*30)
	fmt.Println("init data done")
}
