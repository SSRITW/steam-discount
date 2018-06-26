# steam-discount
抓取steam打折游戏信息保存到redis <br>

### 用到的类库
go get github.com/PuerkitoBio/goquery <br>
go get github.com/garyburd/redigo/redis <br>

### 简介
利用http.NewRequest创建一个使用随机代理的请求 <br>
使用goquery解析到请求到的内容 <br>
将解析到的内容存储在redis内 <br>

### 弊端
目前仍没有把redis连接信息和url写入配置文件内，修改比较麻烦 <br>
解析数据是依据页面html标签的id或者样式获取的，当steam的页面样式改变时，就需要更改解析代码 <br>
