package main

import (
	"fmt"
)

func main(){
	a := []int{1, 2}

	fmt.Println(a)

	a = a[1:]
	fmt.Println(a)

	// a = a[1:]
	// a = a[1:]
	// a = a[1:]
	// fmt.Println(a)

	a = a[1:]
	fmt.Println(a)
}