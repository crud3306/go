package main

import (
	"fmt"

	"github.com/astaxie/beego/httplib"
)

func main() {
	req := httplib.Get("https://www.douban.com/")
	str, err := req.String()
	if err != nil {
		panic(err)
	}

	fmt.Println(str)
}