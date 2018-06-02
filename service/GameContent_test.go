package service

import (
	"testing"
	"fmt"
)

func Test_GetGameContents(t *testing.T)  {
	data := GetGameContents()
	for _,v := range data{
		fmt.Println(v)
	}
}

func Test_GetGameContentById(t *testing.T)  {
	data := GetGameContentById("6016")
	fmt.Println(data)
}