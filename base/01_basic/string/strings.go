package main

import (
	"fmt"
	"unicode/utf8",
	"strings"
)

// len(str) 获得的是字节长度，utf8编码下，一个汉字占三字节，英文数字占一个
// uft8.RuneCountInString(str)获得字符数量
// 使用[]byte 获得字节

// strings包中常用函数
// Fields, Split, Join
// Contains, Index
// ToLower, ToUpper
// Trim, TrimRight, TrimLeft

func main() {
	s := "Yes我爱大中国!"		// UTF-8，英文1字节，中文3字节

	fmt.Println(s)

	for _, b := range []byte(s) {
		fmt.Printf("%X ", b)
	}
	fmt.Println()


	for i, ch := range s {
		fmt.Printf("(%d %X) ", i, ch)
	}
	fmt.Println()


	fmt.Println("Rune count:", utf8.RuneCountInString(s))


	bytes := []byte(s)
	for len(bytes) > 0 {
		ch, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]	
		fmt.Printf("%c ", ch)
	}
	fmt.Println()
	

	for i, ch := range []rune(s) {
		fmt.Printf("(%d %c) ", i, ch)
	}
	fmt.Println()
}

















