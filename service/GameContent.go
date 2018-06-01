package service

import (
	"steam-discount/monitor"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

func GetGameContents()(data []monitor.MonitorContent){
	c,err := redis.Dial("tcp","127.0.0.1:6379",redis.DialPassword("123456"))
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	var tempData []interface{}
	tempData,err = redis.Values(c.Do("ZRANGE","gameContent",0,-1))
	if err!=nil {
		fmt.Println(err.Error())
		return
	}

	if err := redis.ScanSlice(tempData, &data); err != nil {
		panic(err)
	}
	return
}



