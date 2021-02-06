package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var SignalInterruptHandler *InterruptHandler

func InitInterruptHandler() {
	SignalInterruptHandler = NewInterruptHandle()
}

func InterruptHandleAddFunc(f func()) {
	if SignalInterruptHandler != nil {
		SignalInterruptHandler.registerFunc(f)
	}
}

// ---------------------------------------------------
// SIGINT信号监听，func列表执行后从容退出
func NewInterruptHandle() *InterruptHandler {
	h := &InterruptHandler{
		make([]func(), 0),
	}
	go h.handle()
	return h
}

type InterruptHandler struct {
	FuncList []func()
}

func (h *InterruptHandler) registerFunc(f func()) {
	h.FuncList = append(h.FuncList, f)
}

func (h *InterruptHandler) handle() {
	ch := make(chan os.Signal, 0)
	// SIGINT 监听ctrl+C，SIGTERM监听结束(含普通kill pid)
	//signal.Notify(ch, syscall.SIGINT)
	//signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	//注意：signal.Notify同一个信号多次，可以捕获多次。

	if s, ok := <-ch; ok {
		fmt.Println("捕捉到" + s.String() + "信号")
		for _, f := range h.FuncList {
			f()
		}
	}

	//取消监听
	signal.Stop(ch)
	close(ch)
	// 使用了信号监听，可交给注入的函数去平滑结束主协程
	// os.Exit(0)
}
