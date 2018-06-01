package monitor

import (
	"testing"
)


func Test_GetContent(t *testing.T){
	var contentChannel = make(chan MonitorContent,50)
	var pageSizeChannel = make(chan int )
	var maxContentSize = make(chan int)

	go GetContentsByPageChan(pageSizeChannel,contentChannel)

	go GetContent("https://store.steampowered.com/search/?specials=1",contentChannel,pageSizeChannel,maxContentSize)

	SaveContents(contentChannel,maxContentSize)
}
