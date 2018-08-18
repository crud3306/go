package main

import (
	"encoding/json"
	"fmt"
)

// 用到 "encoding/json"

// 把相应数据转换成json
// s, err := json.Marshal(interface{})

// 把相应的json字符串解码
// json.Unmarshal(s string, *interface{})

type Student struct {
	Name string `json:"s_name"`
	Age int
}

func main() {
	// 数组转json
	x := [5]int{1, 2, 3, 4, 5}
	s, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(ss))
	

	// map转json
	x1 := make(map[string]float64)
	x1["zhangsan"] = 100.4
	x1["lishi"] = 101
	s1, err1 := json.Marshal(x1)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(string(s1))


	// struct转json
	student := Student{"zhan", 22}
	s2, err2 := json.Marshal(student)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(string(s2))


	// 对s3对应的json字符串 解码
	var s4 interface{}
	json.Unmarshal([]byte(s2), &s4)
	fmt.Printf("%v", s4)
}