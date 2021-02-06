package main
// 写法

import (
	"fmt"
	"testing"	// 每一个test文件须import一个testing
)

// 执行test case，使用如下命令：
// go test
// go test -v

// 使用TestMain作为初始化test，并且使用m.Run()来调用其他tests可以完成一些需要初始化操作的testing，
// 比如数据库连接，文件打开，REST服务登录等
func TestMain(m *testing.M) {
	fmt.Println("test main first")

	// 如果没有在TestMain中调用m.Run()则只执行TestMain，其它test case不会被执行
	m.Run()
}

// 每个test case均必须以Test开头，并且符合TestXxx形式，否则go test会直接跳过测试不执行
// 必须带参数 t *testing.T
// 如果要测试性能，必须带参数b *testing.B
func TestPrint(t *testing.T) {
	t.SkipNow() // t.SkipNow()为跳过当前test，直接给出PASS结果，然后继续处理其它test case

	fmt.Println("hey")
	// 调用测试代码

	// 比较测试结果与预期结果
	// 如果结果不符合预期
	t.Errorf("return value not valid") // t.Errorf打印错误信息，并且当前test case会被跳过
}

// 每个test case执行并不一定是从上到下的顺序，如果要保证顺序，可以用t.Run来执行subtests
func TestSub(t *testing.T) {
	// 使用t.Run来执行subtests可以做到控制test输出以及test的顺序
	t.Run("a1", func(t *testing.T) { fmt.Println("a1") })
	t.Run("a2", func(t *testing.T) { fmt.Println("a2") })
	t.Run("a3", func(t *testing.T) { fmt.Println("a3") })

	// 调用其它test
	t.Run("testHaha", testHaha)
}

// 此test不能被直接执行，因它的命令不合TestXxxx的规范，但是可以被其它Test case调用
func testHaha(t *testing.T) {
	fmt.Println("haha")
}


// 执行bench，使用如下命令：
// go test -bench=.

// benchmark写法 BenchmarkXxxx
// 须传参数 b *testing.B
// 函数中的循环条件用到b.N
func BenchmarkAll(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// fmt.Printf("bench %d", n)
		aaa(n)
	}
}

func aaa(n int) int {
	// for n > 0 {
	// 	n--
	// }
	return n
}














































