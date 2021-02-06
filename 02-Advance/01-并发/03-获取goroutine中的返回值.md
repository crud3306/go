
获取goroutine中的返回值
=============

执行go协程时, 是没有返回值的, 这时候需要用到go语言中特色的channel来获取到返回值. 通过channel拿到返回值有两种处理形式, 一种形式是具有go风格特色的, 即发送给一个for channel或select channel的独立goroutine中, 由该独立的goroutine来处理函数的返回值. 还有一种传统的做法, 就是将所有goroutine的返回值都集中到当前函数, 然后统一返回给调用函数.

- 发送给独立的goroutine处理程序
- 在当前函数中聚合返回


发送给独立的goroutine处理程序
-------------
```golang
package main

import (
	"fmt"
	"sync"
	"time"
)

var responseChannel = make(chan string, 15)

func httpGet(url int, limiter chan bool, wg *sync.WaitGroup) {
	// 函数执行完毕时 计数器-1
	defer wg.Done()
	fmt.Println("http get:", url)
	responseChannel <- fmt.Sprintf("Hello Go %d", url)
	// 释放一个坑位
	<- limiter
}

func ResponseController() {
	for rc := range responseChannel {
		fmt.Println("response: ", rc)
	}
}

func main() {

	// 启动接收response的控制器
	go ResponseController()

	wg := &sync.WaitGroup{}
	// 控制并发数为10
	limiter := make(chan bool, 20)

	for i := 0; i < 99; i++ {
		// 计数器+1
		wg.Add(1)
		limiter <- true
		go httpGet(i, limiter, wg)
	}
	// 等待所以协程执行完毕
	wg.Wait() // 当计数器为0时, 不再阻塞
	fmt.Println("所有协程已执行完毕")
}
```
这种具有Go语言特色的处理方式的关键在于, 你需要预先创建一个用于处理返回值的公共管道. 然后定义一个一直在读取该管道的函数, 该函数需要预先以单独的goroutine形式启动.

最后当执行到并发任务时, 每个并发任务得到结果后, 都会将结果通过管道传递到之前预先启动的goroutine中.






在当前函数中聚合返回
------------------
```golang
package main

import (
	"fmt"
	"sync"
)

func httpGet(url int,response chan string, limiter chan bool, wg *sync.WaitGroup) {
	// 函数执行完毕时 计数器-1
	defer wg.Done()
	// 将拿到的结果, 发送到参数中传递过来的channel中
	response <- fmt.Sprintf("http get: %d", url)
	// 释放一个坑位
	<- limiter
}

// 将所有的返回结果, 以 []string 的形式返回
func collect(urls []int) []string {
	var result []string

	wg := &sync.WaitGroup{}
	// 控制并发数为10
	limiter := make(chan bool, 5)
	defer close(limiter)

	// 函数内的局部变量channel, 专门用来接收函数内所有goroutine的结果
	responseChannel := make(chan string, 20)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}
	// 启动读取结果的控制器
	go func() {
		// wgResponse计数器+1
		wgResponse.Add(1)
		// 读取结果
		for response :=  range responseChannel {
			// 处理结果
			result = append(result, response)
		}
		// 当 responseChannel被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
		wgResponse.Done()
	}()

	for _, url := range urls {
		// 计数器+1
		wg.Add(1)
		limiter <- true
		// 这里在启动goroutine时, 将用来收集结果的局部变量channel也传递进去
		go httpGet(url,responseChannel, limiter, wg)
	}

	// 等待所以协程执行完毕
	wg.Wait() // 当计数器为0时, 不再阻塞
	fmt.Println("所有协程已执行完毕")

	// 关闭接收结果channel
	close(responseChannel)

	// 等待wgResponse的计数器归零
	wgResponse.Wait()
	
	// 返回聚合后结果
	return result
}

func main() {
	urls := []int{1,2,3,4,5,6,7,8,9,10}

	result := collect(urls)
	fmt.Println(result)
}
```