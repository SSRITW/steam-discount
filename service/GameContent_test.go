package service

import (
	"testing"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type test struct {
	Id int
	Name string
}

func Test_GetGameContents(t *testing.T)  {
	c,err := redis.Dial("tcp","127.0.0.1:6379",redis.DialPassword("123456"))
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	var tempSlice test
	for i:=0;i<5 ;i++  {
		temp := test{i,"test"}
		_,err := c.Do("ZADD","testcontent",temp.Id,temp)	 //将内容放到redis里
		if err!=nil {
			fmt.Println("redis error:",err.Error())
		}
	}
	tempData,err2 := redis.Values(c.Do("ZRANGE","testcontent",0,-1))
	if err2!=nil {
		fmt.Println(err2.Error())
		return
	}
	e := redis.ScanStruct(tempData, &tempSlice)
	if e!=nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(tempSlice)
}
