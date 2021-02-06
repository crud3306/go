
前言

本篇文章介绍如何分析golang程序的内存使用情况。包含以下几种方法的介绍：

- 执行前添加系统环境变量GODEBUG='gctrace=1'来跟踪打印垃圾回收器信息
- 在代码中使用runtime.ReadMemStats来获取程序当前内存的使用情况
- 使用pprof工具

注意，本篇文章前后有关联，需要顺序阅读。


从十来行的demo开始
=====================
```golang
package main

import (
    "log"
    "runtime"
    "time"
)

func f() {
    container := make([]int, 8)
    log.Println("> loop.")
    // slice会动态扩容，用它来做堆内存的申请
    for i := 0; i < 32*1000*1000; i++ {
        container = append(container, i)
    }
    log.Println("< loop.")
    // container在f函数执行完毕后不再使用
}

func main() {
    log.Println("start.")
    f()

    log.Println("force gc.")
    runtime.GC() // 调用强制gc函数

    log.Println("done.")
    time.Sleep(1 * time.Hour) // 保持程序不退出
}
```

编译并运行
```sh
go build -o snippet_mem && ./snippet_mem

#打印如下信息：
2019/04/06 14:23:16 start.
2019/04/06 14:23:16 > loop.
2019/04/06 14:23:17 < loop.
2019/04/06 14:23:17 force gc.
2019/04/06 14:23:18 done.
```

使用top命令查看snippet_mem进程的内存RSS占用为470M。

> top -p $(pidof snippet_mem)

分析：

直观上来说，这个程序在f()函数执行完后，切片的内存应该被释放，不应该占用470M那么大。



下面让我们使用一些手段来分析程序的内存使用情况。



GODEBUG中的gctrace
=====================
我们在执行demo程序之前添加环境变量GODEBUG='gctrace=1'来跟踪打印垃圾回收器信息

> go build -o snippet_mem && GODEBUG='gctrace=1' ./snippet_mem

在分析demo程序的输出信息之前，先把gctrace输出信息的格式以及字段的含义放前面，一会我们的分析要基于这部分内容。

gctrace输出信息的格式以及字段的含义对应的官方文档：https://godoc.org/runtime

我对它做的翻译如下：
```sh
gctrace: 设置gctrace=1会使得垃圾回收器在每次回收时汇总所回收内存的大小以及耗时，
并将这些内容汇总成单行内容打印到标准错误输出中。
这个单行内容的格式以后可能会发生变化。
目前它的格式：
    gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
各字段的含义：
    gc #        GC次数的编号，每次GC时递增
    @#s         距离程序开始执行时的时间
    #%          GC占用的执行时间百分比
    #+...+#     GC使用的时间
    #->#-># MB  GC开始，结束，以及当前活跃堆内存的大小，单位M
    # MB goal   全局堆内存大小
    # P         使用processor的数量
如果信息以"(forced)"结尾，那么这次GC是被runtime.GC()调用所触发。

如果gctrace设置了任何大于0的值，还会在垃圾回收器将内存归还给系统时打印一条汇总信息。
这个将内存归还给系统的操作叫做scavenging。
这个汇总信息的格式以后可能会发生变化。
目前它的格式：
    scvg#: # MB released  printed only if non-zero
    scvg#: inuse: # idle: # sys: # released: # consumed: # (MB)
各字段的含义:
    scvg#        scavenge次数的变化，每次scavenge时递增
    inuse: #     MB 垃圾回收器中使用的大小
    idle: #      MB 垃圾回收器中空闲等待归还的大小
    sys: #       MB 垃圾回收器中系统映射内存的大小
    released: #  MB 归还给系统的大小
    consumed: #  MB 从系统申请的大小
```

打印如下信息：
```sh
2019/04/06 14:28:26 start.
2019/04/06 14:28:26 > loop.
gc 1 @0.001s 0%: 0.005+0.92+0.004 ms clock, 0.011+0.027/0/0.13+0.009 ms cpu, 4->6->2 MB, 5 MB goal, 2 P
gc 2 @0.003s 0%: 0.002+0.43+0.002 ms clock, 0.004+0.013/0/0.32+0.005 ms cpu, 5->5->1 MB, 6 MB goal, 2 P
gc 3 @0.003s 1%: 0.002+0.47+0.003 ms clock, 0.004+0.027/0/0.44+0.006 ms cpu, 4->4->2 MB, 5 MB goal, 2 P
gc 4 @0.004s 1%: 0.002+0.50+0.003 ms clock, 0.004+0.022/0/0.48+0.007 ms cpu, 5->5->2 MB, 6 MB goal, 2 P
gc 5 @0.004s 1%: 0.001+1.2+0.003 ms clock, 0.003+0.070/0/1.1+0.006 ms cpu, 6->6->3 MB, 7 MB goal, 2 P
gc 6 @0.006s 1%: 0.002+1.8+0.004 ms clock, 0.004+0.027/0.001/1.8+0.008 ms cpu, 8->8->4 MB, 9 MB goal, 2 P
gc 7 @0.008s 1%: 0.002+2.4+0.005 ms clock, 0.005+0.042/0/2.4+0.010 ms cpu, 10->10->5 MB, 11 MB goal, 2 P
gc 8 @0.010s 1%: 0.002+1.0+0.004 ms clock, 0.005+0.025/0/0.99+0.008 ms cpu, 12->12->6 MB, 13 MB goal, 2 P
gc 9 @0.012s 1%: 0.011+1.8+0.005 ms clock, 0.022+0.025/0/1.7+0.010 ms cpu, 15->15->8 MB, 16 MB goal, 2 P
gc 10 @0.014s 1%: 0.002+3.8+0.004 ms clock, 0.005+0.014/0/3.8+0.009 ms cpu, 19->19->10 MB, 20 MB goal, 2 P
gc 11 @0.018s 1%: 0.003+2.0+0.004 ms clock, 0.006+0.026/0/2.0+0.008 ms cpu, 24->24->13 MB, 25 MB goal, 2 P
gc 12 @0.020s 1%: 0.002+3.0+0.005 ms clock, 0.005+0.028/0/3.0+0.011 ms cpu, 30->30->16 MB, 31 MB goal, 2 P
gc 13 @0.024s 0%: 0.003+9.0+0.004 ms clock, 0.006+0.028/0/9.0+0.009 ms cpu, 38->38->21 MB, 39 MB goal, 2 P
gc 14 @0.033s 0%: 0.002+4.6+0.005 ms clock, 0.005+0.036/0/4.6+0.011 ms cpu, 47->47->26 MB, 48 MB goal, 2 P
gc 15 @0.039s 0%: 0.003+13+0.004 ms clock, 0.007+0.024/0/13+0.009 ms cpu, 59->59->33 MB, 60 MB goal, 2 P
gc 16 @0.053s 0%: 0.002+17+0.005 ms clock, 0.005+0.030/0.027/17+0.011 ms cpu, 74->74->41 MB, 75 MB goal, 2 P
gc 17 @0.072s 0%: 0.049+29+0.004 ms clock, 0.098+0.015/0.091/29+0.009 ms cpu, 93->93->51 MB, 94 MB goal, 2 P
gc 18 @0.103s 0%: 0.003+29+0.005 ms clock, 0.007+0.031/0.029/29+0.010 ms cpu, 116->116->64 MB, 117 MB goal, 2 P
gc 19 @0.134s 0%: 0.003+41+0.004 ms clock, 0.006+0.016/0.030/41+0.009 ms cpu, 145->145->80 MB, 146 MB goal, 2 P
gc 20 @0.178s 0%: 0.003+44+0.005 ms clock, 0.006+0.016/0.045/44+0.010 ms cpu, 181->181->101 MB, 182 MB goal, 2 P
gc 21 @0.223s 0%: 0.003+55+0.004 ms clock, 0.006+0.015/0.044/55+0.008 ms cpu, 227->227->126 MB, 228 MB goal, 2 P
gc 22 @0.281s 0%: 0.004+67+0.004 ms clock, 0.009+0.048/0.023/67+0.008 ms cpu, 284->284->157 MB, 285 MB goal, 2 P
gc 23 @0.352s 0%: 0.004+90+0.005 ms clock, 0.008+0.035/0.042/90+0.011 ms cpu, 355->355->197 MB, 356 MB goal, 2 P
2019/04/06 14:28:27 < loop.
2019/04/06 14:28:27 force gc.
gc 24 @0.446s 0%: 0.005+107+0.007 ms clock, 0.010+0.015/0.050/107+0.014 ms cpu, 444->444->0 MB, 445 MB goal, 2 P (forced)
2019/04/06 14:28:27 done.
gc 25 @0.554s 0%: 0.077+0.071+0.002 ms clock, 0.15+0/0.078/0.036+0.004 ms cpu, 0->0->0 MB, 8 MB goal, 2 P (forced)
```

这里顺便一提，gc的打印信息和demo程序log打印的信息是并行往标准错误输出打印的，所以可能会乱，上面所贴的打印信息的倒数第3、4行是我自己手动重排了，重排前的信息如下：
```sh
gc 24 @0.446s 0%: 2019/04/06 14:28:27 force gc.
0.005+107+0.007 ms clock, 0.010+0.015/0.050/107+0.014 ms cpu, 444->444->0 MB, 445 MB goal, 2 P (forced)
```

demo程序之后会每隔一段时间打印一些gc信息，汇总如下：
```sh
GC forced
gc 26 @120.562s 0%: 0.008+0.18+0.005 ms clock, 0.016+0/0.051/0.10+0.010 ms cpu, 0->0->0 MB, 8 MB goal, 2 P
scvg0: inuse: 0, idle: 959, sys: 959, released: 447, consumed: 512 (MB)
GC forced
gc 27 @240.562s 0%: 0.005+0.19+0.005 ms clock, 0.010+0/0.063/0.13+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
scvg1: 512 MB released
scvg1: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
gc 28 @360.564s 0%: 0.007+0.099+0.004 ms clock, 0.014+0/0.036/0.13+0.008 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 29 @480.565s 0%: 0.006+0.30+0.005 ms clock, 0.013+0/0.048/0.12+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
scvg2: 0 MB released
scvg2: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
GC forced
gc 30 @600.566s 0%: 0.004+0.11+0.005 ms clock, 0.009+0/0.045/0.15+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
scvg3: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
GC forced
gc 31 @720.566s 0%: 0.004+0.081+0.004 ms clock, 0.009+0/0.024/0.10+0.008 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 32 @840.567s 0%: 0.006+0.12+0.005 ms clock, 0.012+0/0.039/0.17+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
scvg4: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
```

分析：

先看在f()函数执行完后立即打印的gc 24那行的信息。444->444->0 MB, 445 MB goal表示垃圾回收器已经把444M的内存标记为非活跃的内存。
再看0.1秒之后的gc 25。0->0->0 MB, 8 MB goal表示垃圾回收器中的全局堆内存大小由445M下降为8M。

结论：在f()函数执行完后，demo程序中的切片容器所申请的堆空间都被垃圾回收器回收了。


但是此时top显示内存依然占用470M。

结论：垃圾回收器回收了应用层的内存后，（可能）并不会立即将内存归还给系统。



接下来看scvg相关的信息。该信息在demo程序每运行一段时间后打印一次。

scvg0时consumed为512M。此时内存还没有归还给系统。
scvg1时consumed为0，并且scvg1的released=(scvg0 released + scvg0 consumed)。此时内存已归还给系统。
我们通过top命令查看，内存占用下降为38M。
之后打印的scvg信息不再有变化。

结论：垃圾回收器在一段时间后，（可能）会将回收的内存归还给系统。


到这里，我们对GODEBUG中的gctrace的用法已经介绍完毕了。
实时上，我们最前面的疑问也解决了。

但是我们接下来依然会使用另外几次方法来分析我们的domo程序。



runtime.ReadMemStats
=====================
我们稍微修改一下demo程序，在一些执行流程上以及f()函数执行完后每10秒使用runtime.ReadMemStats获取内存使用情况。
```golang
package main

import (
    "log"
    "runtime"
    "time"
)

func traceMemStats() {
    var ms runtime.MemStats
    runtime.ReadMemStats(&ms)
    log.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
}

func f() {
    container := make([]int, 8)
    log.Println("> loop.")
    for i := 0; i < 32*1000*1000; i++ {
        container = append(container, i)
        if i == 16*1000*1000 {
            traceMemStats()

        }
    }
    log.Println("< loop.")
}

func main() {
    log.Println("start.")
    traceMemStats()
    f()
    traceMemStats()

    log.Println("force gc.")
    runtime.GC()

    log.Println("done.")
    traceMemStats()

    go func() {
        for {
            traceMemStats()
            time.Sleep(10 * time.Second)
        }
    }()

    time.Sleep(1 * time.Hour)
}
```

打印如下信息：
```sh
2019/04/06 17:37:52 start.
2019/04/06 17:37:52 Alloc:49328(bytes) HeapIdle:66494464(bytes) HeapReleased:0(bytes)
2019/04/06 17:37:52 > loop.
2019/04/06 17:37:52 Alloc:238510080(bytes) HeapIdle:364863488(bytes) HeapReleased:334856192(bytes)
2019/04/06 17:37:52 < loop.
2019/04/06 17:37:52 Alloc:207053496(bytes) HeapIdle:664731648(bytes) HeapReleased:396263424(bytes)
2019/04/06 17:37:52 force gc.
2019/04/06 17:37:52 done.
2019/04/06 17:37:52 Alloc:49864(bytes) HeapIdle:871768064(bytes) HeapReleased:396255232(bytes)
2019/04/06 17:37:52 Alloc:51056(bytes) HeapIdle:871727104(bytes) HeapReleased:396222464(bytes)
// ... 省略部分日志
2019/04/06 17:42:32 Alloc:52304(bytes) HeapIdle:871718912(bytes) HeapReleased:396214272(bytes)
2019/04/06 17:42:42 Alloc:52416(bytes) HeapIdle:871718912(bytes) HeapReleased:396214272(bytes)
2019/04/06 17:42:52 Alloc:52528(bytes) HeapIdle:871718912(bytes) HeapReleased:603217920(bytes)
2019/04/06 17:43:02 Alloc:52640(bytes) HeapIdle:871718912(bytes) HeapReleased:871653376(bytes)
2019/04/06 17:43:12 Alloc:52752(bytes) HeapIdle:871718912(bytes) HeapReleased:871653376(bytes)
2019/04/06 17:43:22 Alloc:52864(bytes) HeapIdle:871718912(bytes) HeapReleased:871653376(bytes)
```

可以看到，打印done.之后那条trace信息，Alloc已经下降，即内存已被垃圾回收器回收。在2019/04/06 17:42:52和2019/04/06 17:43:02的两条trace信息中，HeapReleased开始上升，即垃圾回收器把内存归还给系统。距离打印done.时有5分钟时间间隔。

另外，MemStats还可以获取其它哪些信息以及字段的含义可以参见官方文档：
http://golang.org/pkg/runtime/#MemStats




使用pprof工具
=====================
在网页上查看内存使用情况，需在代码中添加两行代码
```golang
import(
    "net/http"
    _ "net/http/pprof"
)

go func() {
    log.Println(http.ListenAndServe("0.0.0.0:10000", nil))
}()
```
然后就可以使用浏览器打开以下地址查看内存信息 http://127.0.0.1:10000/debug/pprof/heap?debug=1

使用此方法，除了有MemStats的信息，还有申请内存发生在哪些函数的信息。



总结
=====================
golang的垃圾回收器在回收了应用层的内存后，有可能并不会立即将回收的内存归还给操作系统。

如果我们要观察应用层代码使用的内存大小，可以观察Alloc字段。
如果我们要观察程序从系统申请的内存以及归还给系统的情况，可以观察HeapIdle和HeapReleased字段。

以上3种方法，都是获取了程序的MemStats信息。区别是：第一种完全不用修改程序，第二种可以在指定位置获取信息，第三种可以查看具体哪些函数申请了内存。

