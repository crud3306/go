package main 

import (
	"crypto/md5"
	"fmt"
)

func main() {
	md5Inst := md5.New()
	md5Inst.Write([]byte("zhangesan"))
	result := md5Inst.Sum([]byte(""))
	// fmt.Printf("%v", string(result))
	fmt.Printf("%x\n", result)
}