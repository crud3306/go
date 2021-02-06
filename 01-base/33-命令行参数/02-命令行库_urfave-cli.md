
github.com/urfave/cli 使用


目录示例
-------------
```sh
core
├── comm
│   └── flag.go
├── console
│   └── console_main.go
├── http
│   └── http_main.go
├── service
│   ├── config.go
│   └── signal.go
├── main.go
```



code示例
--------------

main.go
```golang
package main

import (
	"xxx/common"
	"xxx/console"
	"xxx/http"
	"xxx/service"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

/*
 *############################### rpc 服务 ###############################
 * go run main.go rpc -c ./config/qa
 *
 */
func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		fmt.Println(c.String("name"), "===")
		return nil
	}

	// 通用参数
	app.Flags = []cli.Flag{
		// 环境变量
		cli.StringFlag{
			Name:        "env, e",
			Value:       "dev",
			Usage:       "environment type",
			Destination: &common.FlagEnv,
		},
		// 指定配置文件路径
		cli.StringFlag{
			Name:        "config, c",
			Value:       "",
			Usage:       "config file path",
			Destination: &common.FlagConfigPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "http",           // 命令全称, 命令简写
			Aliases: []string{"http"}, // 命令简写
			Usage:   "http service",   // 命令详细描述
			Action: func(c *cli.Context) { // 命令处理函数
				// 监听配置, 后续从watcher中获取配置内容, config服务基于他获取配置
				service.InitConfigWatcher(common.FlagEnv, common.FlagConfigPath)

				// 初始化全局		信号监听服务
				service.InitInterruptHandler()
				// 注册全局的interrupt信号监听调用函数
				service.InterruptHandleAddFunc(xxxStopFunc)

				// 启动服务
				http.Main()
			},
		},
		{
			Name:    "task",           // 命令全称, 命令简写
			Aliases: []string{"task"}, // 命令简写
			Usage:   "task service",   // 命令详细描述
			Flags: []cli.Flag{ // 各command的参数
				// 任务资源
				cli.StringFlag{
					Name:        "resource, r",
					Value:       "runall",
					Usage:       "run resource",
					Destination: &common.FlagResource,
				},
			},
			Action: func(c *cli.Context) { // 命令处理函数

				// 启动服务
				console.Main(common.FlagResource)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

```


common/flag.go
```golang
package common

// FlagEnv 环境
var FlagEnv string

// FlagConfigPath 配置文件
var FlagConfigPath string

// FlagResource 资源
var FlagResource string
```


service/signal.go
```golang
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
```


http/http_main.go
```golang
package http
import "fmt"

func Main(){
	fmt.Println("http")
	//这里启动http服务
}
```


console/console_main.go
```golang
package console
import "fmt"

func Main(flagResource string){
	fmt.Println("console", flagResource)

	// 这里通过flagResource, 执行不同的业务逻辑
}
```


