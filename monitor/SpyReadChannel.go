package monitor


import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

//当pageChan有值时才进行读操作（即，进行第一次爬虫获取到这个值的时候；启动多个goroutine继续爬取之后的页面数据
func GetContentsByPageChan(pageChan chan int,contentChannel chan MonitorContent){
	maxI,ok := <- pageChan
	if ok {
		for i:=2; i<=maxI ;i++  {
			go GetContent("https://store.steampowered.com/search/?specials=1&page="+strconv.Itoa(i),contentChannel,nil,nil)
		}
		close(pageChan)
	}
}

//当maxContentSize有值时才进行读操作（即，进行第一次爬虫获取到这个值的时候)
func SaveContents(contentChannel chan MonitorContent,maxContentSize chan int){
	i := 0
	var size int
	var c redis.Conn

	if v,f :=<-maxContentSize;f{
		size = v
		//fmt.Println("contentSize",size)
		close(maxContentSize)
		var err error
		c,err = redis.Dial("tcp","127.0.0.1:6379",redis.DialPassword("123456"))
		if err!=nil {
			panic(err)
		}
	}
	for v := range contentChannel{
		i++
		_,err := c.Do("ZADD","gameContent",v.Id,v)	 //将内容放到redis里
		if err!=nil {
			fmt.Println("redis error:",err.Error())
		}
		/*else{
			fmt.Println("reply:",reply)
		}*/
		//fmt.Println(v.Thumbnail,i)
		if size == i {
			close(contentChannel)
		}
	}
}
