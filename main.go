package main

import "steam-discount/monitor"

func main(){
	initData()
}

func initData(){
	var contentChannel = make(chan monitor.MonitorContent,50)
	var pageSizeChannel = make(chan int )
	var maxContentSize = make(chan int)

	go monitor.GetContentsByPageChan(pageSizeChannel,contentChannel)

	go monitor.GetContent("https://store.steampowered.com/search/?specials=1",contentChannel,pageSizeChannel,maxContentSize)

	monitor.SaveContents(contentChannel,maxContentSize)
}
