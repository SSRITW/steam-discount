package monitor

import (
	"testing"
	"fmt"
)

func Test_GetContent(t *testing.T){
	result := GetContent("https://store.steampowered.com/search/?specials=1");
	for _,v := range result{
		fmt.Println(v.Thumbnail)
	}
}
