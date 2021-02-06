package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func main() {
	
	get()

	// post()

}

func get() {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	s, err := ioutil.ReadAll(resp.Body)
	// s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", s)
}

func post() {
	resp, err := http.Post("http://www.baidu.com", "application/x-www-form-urlencoded", strings.NewReader("id=1"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	s, err := ioutil.ReadAll(resp.Body)
	// s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", s)
}
