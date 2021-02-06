package main

import (
	"fmt"
	"io/ioutil"
)

// switch基本用法
// 注意：在go中，switch的每个case会自动break，所以不用写。如果不想break，需要加fallthrough
func eval(a, b int, op string) int {
	var result int
	switch op {
	case "+":
		result = a + b
		// 试一下fallthrough语句
		//fallthrough
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	default:
		panic("unsupported operator:" + op)
	}

	return result
}

// switch语名，注意：switch后面可以不带表达式，这时需在case中带
func grade(score int) string {
	g := ""
	switch {
	case score < 0 || score > 100:
		panic(fmt.Sprintf(
			"Wrong score: %d", score))
	case score < 60:
		g = "F"
	case score < 80:
		g = "C"
	case score < 90:
		g = "B"
	case score <= 100:
		g = "A"
	}
	return g
}

func main() {
	// 基本的if写法
	// const filename = "abc.txt"
	// contents, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s\n", contents)
	// }

	// If "abc.txt" is not found,
	// please check what current directory is,
	// and change filename accordingly.
	const filename = "abc.txt"
	// if 的简写方式
	if contents, err := ioutil.ReadFile(filename); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", contents)
	}
	// 注意：通这种方式写，赋值的contents,err 变量在if结构外面是读不到的，只能在if的结构里读到

	fmt.Println(
		grade(0),
		grade(59),
		grade(60),
		grade(82),
		grade(99),
		grade(100),
		// Uncomment to see it panics.
		// grade(-3),
	)
}
