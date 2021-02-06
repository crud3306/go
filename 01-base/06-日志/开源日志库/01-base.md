
在Github中stars数最多的Go日志库集合


在Go语言世界中，日志库并不像Java世界那里有一个具有统治力的日志库。在做新项目技术选型的时候，难免会遇到日志库的选择问题，今天笔者就为大家介绍一下在Github中stars数最多的Go日志库。

logrus是我已知的Go语言日志库中在Github中stars数最多的日志库，功能强大，性能也不错。另外值得一提的是logrus的一个fork库gogap/logrus，可以配置输出到文件和graylog日志系统，基础的level、调用链、行号，文件切分都有，特色是结构化日志可以被统计和索引，借助graylog就可以做日志分析，不过这个库很久没有维护了。

zap是Go中一个快速，结构化，分级日志库，Uber出品。能够激活强大的ad-hoc分析, 灵活的仪表板, 其功能超过其他工具ELK, Splunk, 和 Sentry. 所有日志消息都是JSON-serialized。也是目前笔者使用的日志库之一。

zerolog也是一款性能相当好的日志库，有趣的是zerolog和zap都在自家的Github库首页中的性能对比数据里打败了对方：）

Seelog是一个原生Go日志库,提供了灵活的异步调度、过滤和格式化。这也是笔者较早期使用的Go日志库之一。

blog4go是高性能日志库。创新地使用“边解析边输出”方法进行日志输出，同时支持回调函数、日志淘汰和配置文件。可以解决高并发，调用日志函数频繁的情境下，日志库造成的性能问题。

有一些Github库star数并不多，但是功能却很多，例如mkideal/log、go-log、alog等就是功能十分齐全的日志库，简单易用。还有一些是对开发工程中有些小帮助的库，例如happierall/l。


下面是Go语言日志库列表
```sh
Name	Stars	Forks	Description

logrus	16.5k	775	Structured, pluggable logging for Go.
zap		11.3k	175	Blazing fast, structured, leveled logging in Go.

oklog	1914	61	A distributed and coördination-free log management system
glog	1496	307	Leveled execution logs for Go
Seelog	960	173	Seelog is a native Go logging library that provides flexible asynchronous dispatching, filtering, and formatting.
log15	625	79	Structured, composable logging for Go
zerolog	444	18	Zero Allocation JSON Logger
apex/log	433	38	Structured logging package for Go.
log		215	14	Simple, configurable and scalable Structured Logging for Go.
blog4go	189	31	BLog4go is an efficient logging library written in the Go programming language, providing logging hook, log rotate, filtering and formatting log message.
logutils	176	18	Utilities for slightly better logging in Go (Golang).
log4go	161	109	Logging package similar to log4j for the Go programming language
fileLogger	80	27	fileLogger是一个基于Go开发的可自动分割文件进行备份的异步日志库
gogap/logrus	75	775	Obsolete, Please refer to gogap/logrus_mate
ozzo-log	74	17	A Go (golang) package providing high-performance asynchronous logging, message filtering by severity and category, and multiple message targets.
azer/logger	74	9	Minimalistic logging library for Go.
alexcesaro/log	42	4	Logging packages for Go
happierall/l	33	3	Golang Pretty Logger.Custom go logger for pretty print, log, debug, warn, error with colours and levels.
mkideal/log	33	2	pluginable, fast,structrued and leveled logging package
slf		33	1	Structured Logging Facade (SLF) for Go
logex	29	6	An golang log lib, supports tracking and level, wrap by standard log lib
gologger	26	5	Simple Logger for golang. Logs Into console, file or ElasticSearch. Simple, easy to use.
go-log	24	9	A logger, for Go
slog	22	1	The reference SLF (structured logging facade) implementation for Go
cxr29/log	18	14	log - Go level and rotate log
ulog	15	2	ulog - Structured and context based logging for golang
siddontang/go-log	14	5	a golang log lib supports level and multi handlers
ccpaging/log4go	13	109	Logging package similar to log4j for the Go programming language
mlog	11	9	A simple logging module for go, with a rotating file feature and console logging.
alog	6	2	Golang async log package
golog	5	1	golog is a multilayer & leveled & structured logger for golang.
szxp/log	5	0	A small structured logging library for Golang
go-async-log	4	2	Golang异步日志库，支持异步批量写入，按天或者小时自动切割，错误等级，多文件等
```