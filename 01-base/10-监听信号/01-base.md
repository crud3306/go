
信号(Signal)是Linux, 类Unix和其它POSIX兼容的操作系统中用来进程间通讯的一种方式。对于Linux系统来说，信号就是软中断，用来通知进程发生了异步事件。一个信号就是一个异步的通知，发送给某个进程，或者同进程的某个线程，告诉它们某个事件发生了。

当信号发送到某个进程中时，操作系统会中断该进程的正常流程，并进入相应的信号处理函数执行操作，完成后再回到中断的地方继续执行。

如果目标进程先前注册了某个信号的处理程序(signal handler),则此处理程序会被调用，否则缺省的处理程序被调用。



Go程序中处理Signal信号，比如收到SIGTERM信号后优雅的关闭程序，以及 goroutine结束通知等。

Go 语言提供了对信号处理的包（os/signal）。

Go 中对信号的处理主要使用os/signal包中的两个方法：一个是notify方法用来监听收到的信号；一个是 stop方法用来取消监听。

Go信号通知机制可以通过往一个channel中发送os.Signal实现。



发送信号
============
- kill 系统调用(system call)可以用来发送一个特定的信号给进程。
- kill 命令允许用户发送一个特定的信号给进程。
- raise 库函数可以发送特定的信号给当前进程。

在Linux下运行man kill可以查看此命令的介绍和用法。

The command kill sends the specified signal to the specified process or process group. If no signal is specified, the TERM signal is sent. The TERM signal will kill processes which do not catch this signal. For other processes, it may be necessary to use the KILL (9) signal, since this signal cannot be caught.
Most modern shells have a builtin kill function, with a usage rather similar to that of the command described here. The ‘-a’ and ‘-p’ options, and the possibility to specify pids by command name is a local extension.
If sig is 0, then no signal is sent, but error checking is still performed.

一些异常比如除以0或者segmentation violation相应的会产生SIGFPE和SIGSEGV信号，缺省情况下导致core dump和程序退出。
内核在某些情况下发送信号，比如在进程往一个已经关闭的管道写数据时会产生SIGPIPE信号.
在进程的终端敲入特定的组合键也会导致系统发送某个特定的信号给此进程：

- Ctrl-C 发送 INT signal (SIGINT)，通常导致进程结束
- Ctrl-Z 发送 TSTP signal (SIGTSTP); 通常导致进程挂起(suspend)
- Ctrl-\ 发送 QUIT signal (SIGQUIT); 通常导致进程结束 和 dump core.
- Ctrl-T (不是所有的UNIX都支持) 发送INFO signal (SIGINFO); 导致操作系统显示此运行命令的信息
- kill pid 会发送 SIGTERM信号给进程。实际上默认执行的：kill 15 pid
- kill -9 pid 会发送 SIGKILL信号给进程。



处理信号
============
Signal handler可以通过signal()系统调用进行设置。如果没有设置，缺省的handler会被调用，当然进程也可以设置忽略此信号。
有两种信号不能被拦截和处理: SIGKILL和SIGSTOP。

当接收到信号时，进程会根据信号的响应动作执行相应的操作，信号的响应动作有以下几种：

1.中止进程(Term)  
2.忽略信号(Ign)  
3.中止进程并保存内存信息(Core)  
4.停止进程(Stop)   
5.继续运行进程(Cont)  
6.用户可以通过signal或sigaction函数修改信号的响应动作（也就是常说的”注册信号”）。另外，在多线程中，各线程的信号响应动作都是相同的，不能对某个线程设置独立的响应动作。  



信号类型
============
不同平台的信号定义或许有些不同。下面列出了POSIX中定义的信号。 
Linux 使用34-64信号用作实时系统中。
命令man 7 signal提供了官方的信号介绍。

查看全部信息
```sh
kill -l
#输出
 1) SIGHUP      2) SIGINT       3) SIGQUIT      4) SIGILL       5) SIGTRAP
 6) SIGABRT     7) SIGBUS       8) SIGFPE       9) SIGKILL      10) SIGUSR1
11) SIGSEGV     12) SIGUSR2     13) SIGPIPE     14) SIGALRM     15) SIGTERM
16) SIGSTKFLT   17) SIGCHLD     18) SIGCONT     19) SIGSTOP     20) SIGTSTP
21) SIGTTIN     22) SIGTTOU     23) SIGURG      24) SIGXCPU     25) SIGXFSZ
26) SIGVTALRM   27) SIGPROF     28) SIGWINCH    29) SIGIO       30) SIGPWR
31) SIGSYS      34) SIGRTMIN    35) SIGRTMIN+1  36) SIGRTMIN+2  37) SIGRTMIN+3
38) SIGRTMIN+4  39) SIGRTMIN+5  40) SIGRTMIN+6  41) SIGRTMIN+7  42) SIGRTMIN+8
43) SIGRTMIN+9  44) SIGRTMIN+10 45) SIGRTMIN+11 46) SIGRTMIN+12 47) SIGRTMIN+13
48) SIGRTMIN+14 49) SIGRTMIN+15 50) SIGRTMAX-14 51) SIGRTMAX-13 52) SIGRTMAX-12
53) SIGRTMAX-11 54) SIGRTMAX-10 55) SIGRTMAX-9  56) SIGRTMAX-8  57) SIGRTMAX-7
58) SIGRTMAX-6  59) SIGRTMAX-5  60) SIGRTMAX-4  61) SIGRTMAX-3  62) SIGRTMAX-2
63) SIGRTMAX-1  64) SIGRTMAX    
``` 


在POSIX.1-1990标准中定义的信号列表
---------------
```sh
信号	值	动作	说明
SIGHUP	1	Term	终端控制进程结束(终端连接断开)
SIGINT	2	Term	用户发送INTR字符 (Ctrl+C)触发  (注意：常用、重要)
SIGQUIT	3	Core	用户发送QUIT字符 (Ctrl+反斜杠)触发
SIGILL	4	Core	非法指令(程序错误、试图执行数据段、栈溢出等)
SIGABRT	6	Core	调用abort函数触发
SIGFPE	8	Core	算术运行错误(浮点运算错误、除数为零等)
SIGKILL	9	Term	无条件结束程序(不能被捕获、阻塞或忽略)
SIGSEGV	11	Core	无效内存引用(试图访问不属于自己的内存空间、对只读内存空间进行写操作)
SIGPIPE	13	Term	消息管道损坏(FIFO/Socket通信时，管道未打开而进行写操作)
SIGALRM	14	Term	时钟定时信号
SIGTERM	15	Term	结束程序(可以被捕获、阻塞或忽略) （普通kill触发，或 kill 15 pid）  (注意：常用、重要)
SIGUSR1	30,10,16	Term	用户保留
SIGUSR2	31,12,17	Term	用户保留
SIGCHLD	20,17,18	Ign	子进程结束(由父进程接收)
SIGCONT	19,18,25	Cont	继续执行已经停止的进程(不能被阻塞)
SIGSTOP	17,19,23	Stop	停止进程(不能被捕获、阻塞或忽略)
SIGTSTP	18,20,24	Stop	停止进程(可以被捕获、阻塞或忽略)
SIGTTIN	21,21,26	Stop	后台程序从终端中读取数据时触发
SIGTTOU	22,22,27	Stop	后台程序向终端中写数据时触发
```

在SUSv2和POSIX.1-2001标准中的信号列表:
----------------
```sh
信号	值	动作	说明
SIGTRAP	5	Core	Trap指令触发(如断点，在调试器中使用)
SIGBUS	0,7,10	Core	非法地址(内存地址对齐错误)
SIGPOLL		Term	Pollable event (Sys V). Synonym for SIGIO
SIGPROF	27,27,29	Term	性能时钟信号(包含系统调用时间和进程占用CPU的时间)
SIGSYS	12,31,12	Core	无效的系统调用(SVr4)
SIGURG	16,23,21	Ign	有紧急数据到达Socket(4.2BSD)
SIGVTALRM	26,26,28	Term	虚拟时钟信号(进程占用CPU的时间)(4.2BSD)
SIGXCPU	24,24,30	Core	超过CPU时间资源限制(4.2BSD)
SIGXFSZ	25,25,31	Core	超过文件大小资源限制(4.2BSD)
```

Windows中没有SIGUSR1,可以用SIGBREAK或者SIGINT代替。



Go中的Signal发送和处理
==============
有时候我们想在Go程序中处理Signal信号，比如收到SIGTERM信号后优雅的关闭程序(参看下一节的应用)。
Go信号通知机制可以通过往一个channel中发送os.Signal实现。
首先我们创建一个os.Signal channel，然后使用signal.Notify注册要接收的信号。

```golang
package main
import "fmt"
import "os"
import "os/signal"
import "syscall"
func main() {
    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigs
        fmt.Println(sig)
        done <- true
    }()
    
    fmt.Println("awaiting signal")
    <-done
    fmt.Println("exiting")
}
```
go run main.go执行这个程序，敲入ctrl-C会发送SIGINT信号。 此程序接收到这个信号后会打印退出。



一个例子
---------------
假如这样一个场景:  

一个写日志的服务，每次要等到满了200条数据才会触发写入阿里日志服务操作， 就会遇到一个问题，当程序被人为中断退出时，但是又没有达到200数据，如果不做其它处理，则这些数据就会丢失了， 为防止存在上述情况的丢失，可以使用Signal来做，

```golang
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {

    done := make(chan bool, 1)

    log.Println("hello world")

    watchSignal()
    //试下执行多次watchSignal()会怎么样，发现每个都能捕获到信号
    //watchSignal()

    log.Println("main wait")
    <-done
    log.Println("main end")
}

func watchSignal() {
    ch := make(chan os.Signal, 0)
    // signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

    go func() {
        select {
        case _signal, _ := <-ch:
            log.Println("sig", _signal)

            doSomeThingsFunc()

            //取消监听
            signal.Stop(ch)
            close(ch)
            return
        }
    }()
}
```
上述代码的好处时当收到Signal信号之后，就可以选择调用doSomeThingsFunc()方法来作出相应的处理






Go网络服务器如果无缝重启
==============
Go很适合编写服务器端的网络程序。DevOps经常会遇到的一个情况是升级系统或者重新加载配置文件，在这种情况下我们需要重启此网络程序，如果网络程序暂停的时间较长，则给客户的感觉很不好。

如何实现优雅地重启一个Go网络程序呢。

主要要解决两个问题:

- 进程重启不需要关闭监听的端口
- 既有请求应当完全处理或者超时

文章Graceful Restart in Golang中提供了一种方式，而Florian von Bock根据此思路实现了一个框架endless。

此框架使用起来超级简单:
```sh
err := endless.ListenAndServe("localhost:4242", mux)
```
只需替换http.ListenAndServe和http.ListenAndServeTLS。

它会监听这些信号： syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM和syscall.SIGTSTP。

此文章提到的思路是：

通过exec.Command fork一个新的进程，同时继承当前进程的打开的文件(输入输出，socket等)
```golang
file := netListener.File() // this returns a Dup()
path := "/path/to/executable"
args := []string{
    "-graceful"}
cmd := exec.Command(path, args...)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.ExtraFiles = []*os.File{file}
err := cmd.Start()
if err != nil {
    log.Fatalf("gracefulRestart: Failed to launch, error: %v", err)
}
````

子进程初始化


网络程序的启动代码
```golang
server := &http.Server{Addr: "0.0.0.0:8888"}
 var gracefulChild bool
 var l net.Listever
 var err error
 flag.BoolVar(&gracefulChild, "graceful", false, "listen on fd open 3 (internal use only)")
 if gracefulChild {
     log.Print("main: Listening to existing file descriptor 3.")
     f := os.NewFile(3, "")
     l, err = net.FileListener(f)
 } else {
     log.Print("main: Listening on a new file descriptor.")
     l, err = net.Listen("tcp", server.Addr)
 }
```

父进程停止
```golang
if gracefulChild {
    parent := syscall.Getppid()
    log.Printf("main: Killing parent pid: %v", parent)
    syscall.Kill(parent, syscall.SIGTERM)
}
server.Serve(l)
````

因此，处理特定的信号可以实现程序无缝的重启。



其它  

graceful shutdown实现非常的简单，通过简单的信号处理就可以实现。本文介绍的是graceful restart,要求无缝重启，所以所用的技术相当的hack。

Facebook的工程师也提供了http和net的实现： facebookgo。



