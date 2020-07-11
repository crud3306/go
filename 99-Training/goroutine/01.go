package main

import (
    "fmt"
    "time"
)

func subtask(in chan int, out chan int) {
    i := 0
    for {
        fmt.Println("--- come from parent ", <-in)
        out <- i
        i++
    }
}

func test(i int) {
    // time.Sleep(4*time.Second)
    fmt.Printf("test %d \n", i)
}

func main() {
    // var in = make(chan int, 1)
    var in = make(chan int)
    var out = make(chan int)
    // cc := <-out
    // fmt.Println(cc)

    for i := 0; i < 10; i++ {
        go subtask(in, out)

        fmt.Println("come 001")
        in <- i
        fmt.Println("come 002")

        

        go test(1)
        go test(2)

        

        fmt.Println("come 003")

        tmp := <-out
        fmt.Println("come from subtask ", tmp)
    }

    time.Sleep(5*time.Second)
}