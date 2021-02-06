package main

import (
	"demo01/service"
	"fmt"
	"time"
)

func main() {

	task := &Task{
		enable: true,
	}

	// 初始化
	service.InitInterruptHandler()
	// 添加回调方法
	service.InterruptHandleAddFunc(task.SignalCallback)

	task.Do()
}


// Task 任务
type Task struct {
	enable bool
}

func (t *Task) Do() {
	i := 0
	for t.enable {
		i++
		fmt.Println("task doing1", i)
		time.Sleep(5 * time.Second)
		fmt.Println("task doing2", i)
	}

	fmt.Println("task end", i)
}

func (t *Task) SignalCallback() {
	fmt.Println("执行捕获到信号后的收尾工作")
	// 这里的处理逻辑，视业务而定
	// 如果业务是循环处理数据，则此处更改循环的标志，让退出循环即可。 如：wileFlag=false
	t.enable = false
}
