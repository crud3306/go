package main 

import (
	"fmt"
	_ "time"
)

var ch chan int = make(chan int)

func foo() {
	for j:=0; j<10; j++ {
		fmt.Println("iam foo")
		cc := <- ch // 从ch取数据，如果ch中还没放数据，那就挂起main线，直到foo函数中放数据为止
		fmt.Println(cc)
	}
}

func main() {

	fmt.Println(111)    

	fmt.Println(222)    

	ch <- 0 // 向ch中加数据，如果没有其他goroutine来取走这个数据，那么挂起foo, 直到main函数把0这个数据拿走
	go foo()

	fmt.Println(333)    
}