

怎么给WaitGroup加超时时间呢？


示例1
---------------
```golang
func main() {
	var w = &sync.WaitGroup{}
	var ch = make(chan bool)

	w.Add(2)
	go func() {
	    time.Sleep(time.Second * 2)
	    fmt.Println("我2秒")
	    w.Done()
	}()
	go func() {
	    time.Sleep(time.Second * 6)
	    fmt.Println("我6秒")
	    w.Done()
	}()

	go func() {
	    w.Wait()
	    ch <- false
	}()
	
	select {
	case <-time.After(time.Second * 5):
	    fmt.Println("我超时了")
	case <-ch:
	    fmt.Println("我结束了")
	}   
}
```



示例2
---------------
go用chan实现WaitGroup并支持超时

```golang
package main
 
import "fmt"
import "time"
import "sync"
 
type group struct {
	gc chan bool
	tk *time.Ticker
	cap int
	mutex sync.Mutex
}
 
func WaitGroup(timeOuteRec int) *group{
	timeout     := time.Millisecond * time.Duration(timeOuteRec)

	wg := group{
		gc   : make(chan bool),
		cap  :  0,
		tk   : time.NewTicker(timeout),
	}
 
	return &wg
}

func(w *group)Add(index int){
	w.mutex.Lock()
	w.cap = w.cap+index
	w.mutex.Unlock()
 
	go func(w *group, index int) {
		for i := 0; i<index; i++{
			fmt.Println("exec goruntine product")
			w.gc<- true
		}
	}(w,index)
}

func(w *group)Done(){
	<-w.gc
	fmt.Println("exec goruntine consumer")

	w.mutex.Lock()
	w.cap--
	w.mutex.Unlock()
}

func(w *group)Wait(){
	defer w.tk.Stop()
	
	for  {
		select {
		case <-w.tk.C:
			fmt.Println("time out exec over")
			return;
		default:
			w.mutex.Lock()
			if w.cap == 0 {
				fmt.Println("all goruntine exec over")
				return;
			}
			w.mutex.Unlock()
		}
	}
}

func main() {
	fmt.Println("start...")

	wg := WaitGroup(10)

	wg.Add(1)
	xxx := func(wg *group) {
		fmt.Println("doning...")
		wg.Done()
	}
	go xxx(wg)
 
	wg.Wait()
	
	fmt.Println("exec return")
}
```


