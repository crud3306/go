并发编程(goroutine)


7.1 sync.WaitGroup
--------------
Go 语言提供了 sync.WaitGroup 和 channel 两种方式支持协程(goroutine)的并发。

例如我们希望并发下载 N 个资源，多个并发协程之间不需要通信，那么就可以使用 sync.WaitGroup，等待所有并发协程执行结束。
```golang
import (
	"fmt"
	"sync"
	"time"
)

var wg *sync.WaitGroup

func download(url string) {
	fmt.Println("start to download", url)
	time.Sleep(time.Second) // 模拟耗时操作
	wg.Done()
}

func main() {
	wg = sync.WaitGroup{}
	
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go download("a.com/" + string(i+'0'))
	}
	
	wg.Wait()
	fmt.Println("Done!")
}
```
wg.Add(1)：为 wg 添加一个计数，wg.Done()，减去一个计数。
go download()：启动新的协程并发执行 download 函数。
wg.Wait()：等待所有的协程执行结束。

```sh
$ go run .
start to download a.com/2
start to download a.com/0
start to download a.com/1
Done!

real    0m1.563s
```
可以看到串行需要 3s 的下载操作，并发后，只需要 1s。



7.2 channel
--------------
```golang
var ch = make(chan string, 10) // 创建大小为 10 的缓冲信道

func download(url string) {
	fmt.Println("start to download", url)
	time.Sleep(time.Second)
	ch <- url // 将 url 发送给信道
}

func main() {
	for i := 0; i < 3; i++ {
		go download("a.com/" + string(i+'0'))
	}

	for i := 0; i < 3; i++ {
		msg := <-ch // 等待信道返回消息。
		fmt.Println("finish", msg)
	}
	
	fmt.Println("Done!")
}
```
使用 channel 信道，可以在协程之间传递消息。阻塞等待并发协程返回消息。
```sh
$ go run .
start to download a.com/2
start to download a.com/0
start to download a.com/1
finish a.com/2
finish a.com/1
finish a.com/0
Done!

real    0m1.528s
```