package monitor

import (
	"testing"
	"fmt"
)

func Test_GetContent(t *testing.T){
	res,state := GetContent("https://store.steampowered.com/search/?specials=1");
	if state {
		//fmt.Println(res)
		testStrs := FinTargetedContent(res)
		for i,v:=range testStrs{
			fmt.Println("index ",i,v)
		}
	}
}
