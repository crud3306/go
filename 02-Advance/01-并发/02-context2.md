
简介
在 Go http包的Server中，每一个请求都有一个对应的 goroutine 去处理。请求处理函数通常会启动额外的 goroutine 用来访问后端服务，比如数据库和RPC服务。用来处理一个请求的 goroutine 通常需要访问一些与请求特定的数据，比如终端用户的身份认证信息、验证相关的token、请求的截止时间。 当一个请求被取消或超时时，所有用来处理该请求的 goroutine 都应该迅速退出，然后系统才能释放这些 goroutine 占用的资源。

在Google 内部，我们开发了 Context 包，专门用来简化对于处理单个请求的多个 goroutine 之间与请求域的数据、取消信号、截止时间等相关操作，这些操作可能涉及多个 API 调用。


源码剖析
=============

context.Context 接口
-------------
context 包的核心
```golang
//  context 包里的方法是线程安全的，可以被多个 goroutine 使用    
type Context interface {               
    // 当Context 被 canceled 或是 times out 的时候，Done 返回一个被 closed 的channel      
    Done() <-chan struct{}        

    // 在 Done 的 channel被closed 后， Err 代表被关闭的原因   
    Err() error 

    // 如果存在，Deadline 返回Context将要关闭的时间  
    Deadline() (deadline time.Time, ok bool)

    // 如果存在，Value 返回与 key 相关了的值，不存在返回 nil  
    Value(key interface{}) interface{}
}
```
我们不需要手动实现这个接口，context 包已经给我们提供了两个，一个是 Background()，一个是 TODO()，这两个函数都会返回一个 Context 的实例。只是返回的这两个实例都是空 Context。


主要结构
-------------
cancelCtx 结构体继承了 Context ，实现了 canceler 方法：
```golang
//*cancelCtx 和 *timerCtx 都实现了canceler接口，实现该接口的类型都可以被直接canceled
type canceler interface {
    cancel(removeFromParent bool, err error)
    Done() <-chan struct{}
}        

type cancelCtx struct {
    Context
    done chan struct{} // closed by the first cancel call.
    mu       sync.Mutex
    children map[canceler]bool // set to nil by the first cancel call
    err      error             // 当其被cancel时将会把err设置为非nil
}

func (c *cancelCtx) Done() <-chan struct{} {
    return c.done
}

func (c *cancelCtx) Err() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.err
}

func (c *cancelCtx) String() string {
    return fmt.Sprintf("%v.WithCancel", c.Context)
}

//核心是关闭c.done
//同时会设置c.err = err, c.children = nil
//依次遍历c.children，每个child分别cancel
//如果设置了removeFromParent，则将c从其parent的children中删除
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    c.mu.Lock()
    if c.err != nil {
        c.mu.Unlock()
        return // already canceled
    }
    c.err = err
    close(c.done)
    for child := range c.children {
        // NOTE: acquiring the child's lock while holding parent's lock.
        child.cancel(false, err)
    }
    c.children = nil
    c.mu.Unlock()

    if removeFromParent {
        removeChild(c.Context, c) // 从此处可以看到 cancelCtx的Context项是一个类似于parent的概念
    }
}
```

timerCtx 结构继承 cancelCtx
-------------
```golang
type timerCtx struct {
    cancelCtx //此处的封装为了继承来自于cancelCtx的方法，cancelCtx.Context才是父亲节点的指针
    timer *time.Timer // Under cancelCtx.mu. 是一个计时器
    deadline time.Time
}
```


valueCtx 结构继承 cancelCtx
-------------
```golang
type valueCtx struct {
    Context
    key, val interface{}
}
```


主要方法
```golang
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key interface{}, val interface{}) Context
```

- WithCancel 对应的是 cancelCtx ,其中，返回一个 cancelCtx ，同时返回一个 CancelFunc，CancelFunc 是 context 包中定义的一个函数类型：type CancelFunc func()。调用这个 CancelFunc 时，关闭对应的c.done，也就是让他的后代goroutine退出。

- WithDeadline 和 WithTimeout 对应的是 timerCtx ，WithDeadline 和 WithTimeout 是相似的，WithDeadline 是设置具体的 deadline 时间，到达 deadline 的时候，后代 goroutine 退出，而 WithTimeout 简单粗暴，直接 return WithDeadline(parent, time.Now().Add(timeout))。

- WithValue 对应 valueCtx ，WithValue 是在 Context 中设置一个 map，拿到这个 Context 以及它的后代的 goroutine 都可以拿到 map 里的值。



详细 context 包源码解读:go源码解读

使用原则
-------------
不要把 Context 存在一个结构体当中，显式地传入函数。Context 变量需要作为第一个参数使用，一般命名为ctx

即使方法允许，也不要传入一个 nil 的 Context ，如果你不确定你要用什么 Context 的时候传一个 context.TODO

使用 context 的 Value 相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数

同样的 Context 可以用来传递到不同的 goroutine 中，Context 在多个goroutine 中是安全的

使用示例
```golang
package main

import (
    "fmt"
    "time"
    "golang.org/x/net/context"
)

// 模拟一个最小执行时间的阻塞函数
func inc(a int) int {
    res := a + 1                // 虽然我只做了一次简单的 +1 的运算,
    time.Sleep(1 * time.Second) // 但是由于我的机器指令集中没有这条指令,
    // 所以在我执行了 1000000000 条机器指令, 续了 1s 之后, 我才终于得到结果。B)
    return res
}

// 向外部提供的阻塞接口
// 计算 a + b, 注意 a, b 均不能为负
// 如果计算被中断, 则返回 -1
func Add(ctx context.Context, a, b int) int {
    res := 0
    for i := 0; i < a; i++ {
        res = inc(res)
        select {
        case <-ctx.Done():
            return -1
        default:
        }
    }
    for i := 0; i < b; i++ {
        res = inc(res)
        select {
        case <-ctx.Done():
            return -1
        default:
        }
    }
    return res
}

func main() {
    {
        // 使用开放的 API 计算 a+b
        a := 1
        b := 2
        timeout := 2 * time.Second
        ctx, _ := context.WithTimeout(context.Background(), timeout)
        res := Add(ctx, 1, 2)
        fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
    }
    {
        // 手动取消
        a := 1
        b := 2
        ctx, cancel := context.WithCancel(context.Background())
        go func() {
            time.Sleep(2 * time.Second)
            cancel() // 在调用处主动取消
        }()
        res := Add(ctx, 1, 2)
        fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
    }
}
```


官方示例
```golang
// func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

package main

import (
	"context"
	"fmt"
)

func main() {
	// gen generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the context once
	// they are done consuming generated integers not to leak
	// the internal goroutine started by gen.
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return // returning not to leak the goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
```


```golang
//func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// Even though ctx will be expired, it is good practice to call its
	// cancelation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

}
```

```golang
//func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Pass a context with a timeout to tell a blocking function that it
	// should abandon its work after the timeout elapses.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}

}
```


```golang
//func WithValue(parent Context, key, val interface{}) Context

package main

import (
	"context"
	"fmt"
)

func main() {
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, favContextKey("color"))

}
```


