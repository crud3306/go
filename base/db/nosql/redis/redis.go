package main

// goredis源码 
// https://github.com/astaxie/goredis/blob/master/redis.go
// 用法非常简单

import (
	"fmt"
	"github.com/astaxie/goredis"
)

func main() {
	// client := goredis.Client //不能这种方式来写，会报错
	var client goredis.Client
	client.Addr = "192.168.2.221:6379"


	// test string
	// ==========
	err := client.Set("qtest", []byte("hello 123"))
	if err != nil {
		panic(err)
	}

	s, err1 := client.Get("qtest")
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(string(s))
	// fmt.Printf("%s", s)



	// test delete key
	// ===========
	// client.Del("qtest")



	// test hash
	// ==========
	f := make(map[string]interface{})
	f["name"] = "zhangsan"
	f["age"] = 12
	f["sex"] = "male"
	err = client.Hmset("qtest1", f)
	if err != nil {
		panic(err)
	}



	// test zset
	// ==========
	_, err = client.Zadd("qteest2", []byte("haha"), 100)
	if err != nil {
		panic(err)
	}


	// 更多操作见源码
	// https://github.com/astaxie/goredis/blob/master/redis.go
	// https://github.com/astaxie/goredis
}




























