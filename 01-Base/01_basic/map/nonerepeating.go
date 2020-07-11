package main

import (
	"fmt"
)



func lengthOfNonRepeatingSubStr(s string) int {
	// lastOccurred := make(map[byte]int)
	lastOccurred := make(map[rune]int)
	start := 0
	maxLength := 0

	// for i, ch := range []byte(s) { // 只能处理英文字符串
	for i, ch := range []rune(s) {	// 换成rune，则可以处理中文文混排的字符串

		fmt.Println(i, ch, lastOccurred, start, maxLength, lastOccurred[ch])

		lastI, ok := lastOccurred[ch]
		if ok && lastI >= start {
			start = lastI + 1
		}

		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}

		lastOccurred[ch] = i
	}

	return maxLength
}

func main() {
	fmt.Println(lengthOfNonRepeatingSubStr("abcdacbccc"))

	// fmt.Println(lengthOfNonRepeatingSubStr("abcdefg"))

	fmt.Println(lengthOfNonRepeatingSubStr("ic我是中国人中国人123"))

	fmt.Println(lengthOfNonRepeatingSubStr("这里是中国人"))

	fmt.Println(lengthOfNonRepeatingSubStr("一二三三二"))
}























