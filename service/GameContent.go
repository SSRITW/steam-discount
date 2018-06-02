package service

import (
	"steam-discount/monitor"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"encoding/json"
)

//获取所有游戏信息
func GetGameContents()(data []monitor.MonitorContent){
	c,err := redis.Dial("tcp","127.0.0.1:6379",redis.DialPassword("123456"))
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	var tempData map[string]string
	tempData,_ = redis.StringMap(c.Do("HGETALL","gameContent"))

	for _,v := range tempData {
		var temp monitor.MonitorContent
		if tempErr :=json.Unmarshal([]byte(v),&temp); tempErr!=nil{
			panic(tempErr)
		}
		data = append(data,temp)
	}
	return
}

//通过游戏id获取到信息
func GetGameContentById(id string)(data monitor.MonitorContent){
	c,err := redis.Dial("tcp","127.0.0.1:6379",redis.DialPassword("123456"))
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	var tempData interface{}
	tempData,err = c.Do("HGET","gameContent",id)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	if tempData!=nil {
		if tempErr :=json.Unmarshal(tempData.([]byte),&data); tempErr!=nil{
			panic(tempErr)
		}
	}
	return
}



