package main

import (
	"fmt"
	"time"
)

func worker(c chan int, i int) {
	defer func(){
		<- c
	}()

	time.Sleep(1*time.Millisecond)

	fmt.Println(i)
}

func main(){
	channel := make(chan int, 10)
	count := 21
	for i := 0; i < count; i++ {
		// fmt.Printf("11： %d \n", i)

		channel <- i

		// fmt.Printf("22： %d \n", i)

		go worker(channel, i)

		// fmt.Printf("33： %d \n", i)
	}

	for j := 100; j < 110; j++ {
		channel <- j

		// fmt.Printf("44： %d \n", j)
	}
	// close(channel)
	// time.Sleep(5000*time.Millisecond)
}