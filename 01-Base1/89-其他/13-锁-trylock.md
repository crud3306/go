

实现tryLock
===============

Go标准库的sync/Mutex、RWMutex实现了sync/Locker接口， 提供了Lock()和UnLock()方法，可以获取锁和释放锁，我们可以方便的使用它来控制我们对共享资源的并发控制上。

但是标准库中的Mutex.Lock的锁被获取后，如果在未释放之前再调用Lock则会被阻塞住，这种设计在有些情况下可能不能满足我的需求。有时候我们想尝试获取锁，如果获取到了，没问题继续执行，如果获取不到，我们不想阻塞住，而是去调用其它的逻辑，这个时候我们就想要TryLock方法了。

虽然很早(13年)就有人给Go开发组提需求了，但是这个请求并没有纳入官方库中，最终在官方库的清理中被关闭了，也就是官方库目前不会添加这个方法。

顺便说一句， sync/Mutex的源代码实现可以访问这里，它应该是实现了一种自旋(spin)加休眠的方式实现， 有兴趣的读者可以阅读源码，或者阅读相关的文章，比如 Go Mutex 源码剖析。这不是本文要介绍的内容，读者可以找一些资料来阅读。

好了，转入正题，看看几种实现TryLock的方式吧。



使用 unsafe 操作指针
---------------
如果你查看sync/Mutex的代码，会发现Mutext的数据结构如下所示：
```golang
type Mutex struct {
	state int32
	sema  uint32
}
```
它使用state这个32位的整数来标记锁的占用，所以我们可以使用CAS来尝试获取锁。

实现的逻辑相对简单，就是使用golang atomic标准库做compareAndSet原子更新, 如果更新成功为拿到锁，否则，反之。 cas 底层是依赖cpu的指令集的，cas的操作包含三个操作数 —— 内存位置（V）、预期原值（A）和新值(B)。 如果内存位置的值与预期原值相匹配，那么处理器会自动将该位置值更新为新值 。否则，处理器不做任何操作。

代码实现如下：
```golang
package trylock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}
```
使用起来和标准库的Mutex用法一样。
```golang
func main() {
	var m trylock.Mutex
	m.Lock()
	go func() {
		m.Lock()
	}()
	
	time.Sleep(time.Second)
	fmt.Printf("TryLock: %t\n", m.TryLock()) //false
	fmt.Printf("TryLock: %t\n", m.TryLock()) // false
	m.Unlock()
	fmt.Printf("TryLock: %t\n", m.TryLock()) //true
	fmt.Printf("TryLock: %t\n", m.TryLock()) //false
	m.Unlock()
	fmt.Printf("TryLock: %t\n", m.TryLock()) //true
	m.Unlock()
}
```
注意TryLock不是检查锁的状态，而是尝试获取锁，所以TryLock返回true的时候事实上这个锁已经被获取了。



或者封装成如下形式
```golang
package trylock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	LockedFlag   int32 = 1
	UnlockedFlag int32 = 0
)

type Mutex struct {
	in     sync.Mutex
	status *int32
}

func NewMutex() *Mutex {
	status := UnlockedFlag
	return &Mutex{
		status: &status,
	}
}

func (m *Mutex) Lock() {
	m.in.Lock()
}

func (m *Mutex) Unlock() {
	m.in.Unlock()
	atomic.AddInt32(m.status, UnlockedFlag)
}

func (m *Mutex) TryLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), UnlockedFlag, LockedFlag) {
		atomic.AddInt32(m.status, LockedFlag)
		return true
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	if atomic.LoadInt32(m.status) == LockedFlag {
		return true
	}
	return false
}
```



实现自旋锁
---------------
上面一节给了我们启发，利用 uint32和CAS操作我们可以一个自定义的锁:
```golang
type SpinLock struct {
	f uint32
}
func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.f, 0)
}
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.f, 0, 1)
}
```
整体来看，它好像是标准库的一个精简版，没有休眠和唤醒的功能。

当然这个自旋锁可以在大并发的情况下CPU的占用率可能比较高，这是因为它的Lock方法使用了自旋的方式，如果别人没有释放锁，这个循环会一直执行，速度可能更快但CPU占用率高。

当然这个版本还可以进一步的优化，尤其是在复制的时候。下面是一个优化的版本:
```golang
type spinLock uint32
func (sl *spinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched() //without this it locks up on GOMAXPROCS > 1
	}
}
func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}
func (sl *spinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1)
}
func SpinLock() sync.Locker {
	var lock spinLock
	return &lock
}
```



使用 channel 实现
---------------
另一种方式是使用channel:
```golang
type ChanMutex chan struct{}

func (m *ChanMutex) Lock() {
	ch := (chan struct{})(*m)
	ch <- struct{}{}
}

func (m *ChanMutex) Unlock() {
	ch := (chan struct{})(*m)
	select {
	case <-ch:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *ChanMutex) TryLock() bool {
	ch := (chan struct{})(*m)
	select {
	case ch <- struct{}{}:
		return true
	default:
	}
	return false
}
```
有兴趣的同学可以关注我的同事写的库 lrita/gosync。


channel lock 示例2
```golang
package main

import (
    "sync"
)

val (
	qLock *Lock
	_lock sync.Mutex
)

// Lock try lock
type QLock struct {
    c chan struct{}
}

// NewLock generate a try lock
func NewLock() Lock {
	_lock.Lock()
	defer _lock.Unlock()
	
	if qLock == nil {
		qLock = &QLock{}
	    l.c = make(chan struct{}, 1)
	    
	    l.c <- struct{}{}	
	}
    
    return l
}

// Lock try lock, return lock result
func (l *QLock) Lock() bool {
    lockResult := false
    select {
    case <-l.c:
        lockResult = true
    default:
    }
    return lockResult
}

// Unlock , Unlock the try lock
func (l *QLock) Unlock() {
    l.c <- struct{}{}
}

var counter int

func main() {
    var l = NewLock()
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if !l.Lock() {
                // log error
                println("lock failed")
                return
            }

            counter++
            println("current counter", counter)
            l.Unlock()
        }()
    }
    wg.Wait()
}
```


性能比较
---------------
首先看看上面三种方式和标准库中的Mutex、RWMutex的Lock和Unlock的性能比较：
```sh
BenchmarkMutex_LockUnlock-4         	100000000	        16.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkRWMutex_LockUnlock-4       	50000000	        36.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnsafeMutex_LockUnlock-4   	100000000	        16.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannMutex_LockUnlock-4    	20000000	        65.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkSpinLock_LockUnlock-4      	100000000	        18.6 ns/op	       0 B/op	       0 allocs/op
```
可以看到单线程(goroutine)的情况下｀spinlock｀并没有比标准库好多少，反而差一点,并发测试的情况比较好，如下表中显示，这是符合预期的。

unsafe方式和标准库差不多。

channel方式的性能就比较差了。

```sh
BenchmarkMutex_LockUnlock_C-4         	20000000	        75.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkRWMutex_LockUnlock_C-4       	20000000	       100 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnsafeMutex_LockUnlock_C-4   	20000000	        75.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannMutex_LockUnlock_C-4    	10000000	       231 ns/op	       0 B/op	       0 allocs/op
BenchmarkSpinLock_LockUnlock_C-4      	50000000	        32.3 ns/op	       0 B/op	       0 allocs/op
```


再看看三种实现TryLock方法的锁的性能：
```sh
BenchmarkUnsafeMutex_Trylock-4        	50000000	        34.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannMutex_Trylock-4         	20000000	        83.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkSpinLock_Trylock-4           	50000000	        30.9 ns/op	       0 B/op	       0 allocs/op
```