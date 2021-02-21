

使用Air实现Go程序实时热重载

今天我们要介绍一个神器——Air能够实时监听项目的代码文件，在代码发生变更之后自动重新编译并执行，大大提高gin框架项目的开发效率。


为什么需要实时加载？
===================
之前使用Python编写Web项目的时候，常见的Flask或Django框架都是支持实时加载的，你修改了项目代码之后，程序能够自动重新加载并执行（live-reload），这在日常的开发阶段是十分方便的。

在使用Go语言的gin框架在本地做开发调试的时候，经常需要在变更代码之后频繁的按下Ctrl+C停止程序并重新编译再执行，这样就不是很方便。


Air介绍
===================
怎样才能在基于gin框架开发时实现实时加载功能呢？像这种烦恼肯定不会只是你一个人的烦恼，所以我报着肯定有现成轮子的心态开始了全网大搜索。果不其然就在Github上找到了一个工具：Air。它支持以下特性：

- 彩色日志输出
- 自定义构建或二进制命令
- 支持忽略子目录
- 启动后支持监听新目录
- 更好的构建过程


安装Air
===================
Go
这也是最经典的安装方式：
```sh
go get -u github.com/cosmtrek/air
```

MacOS
```sh
curl -fLo air https://git.io/darwin_air
```

Linux
```sh
curl -fLo air https://git.io/linux_air
```

Windows
```sh
curl -fLo air.exe https://git.io/windows_air
```

Dcoker
```sh
docker run -it --rm \
    -w "<PROJECT>" \
    -e "air_wd=<PROJECT>" \
    -v $(pwd):<PROJECT> \
    -p <PORT>:<APP SERVER PORT> \
    cosmtrek/air
    -c <CONF>
```
然后按照下面的方式在docker中运行你的项目：
```sh
docker run -it --rm \
    -w "/go/src/github.com/cosmtrek/hub" \
    -v $(pwd):/go/src/github.com/cosmtrek/hub \
    -p 9090:9090 \
    cosmtrek/air
```


使用Air
===================

首先进入你的项目目录：
```sh
cd /path/to/your_project
```

最简单的用法就是直接执行下面的命令：
```sh
air 

#或者 指定air配置文件
#air -c .air.conf
```

执行air前提，你的web服务是可以跑起来的。 执行air后，它会自动执行类似go run main.go的命令



推荐的使用方法是：
```sh
# 1. 在当前目录创建一个新的配置文件.air.conf
touch .air.conf

# 2. 复制 `air.conf.example` 中的内容到这个文件，然后根据你的需要去修改它

# 3. 使用你的配置运行 air, 如果文件名是 `.air.conf`，只需要执行 `air`。
air
```


air_example.conf示例
------------------
完整的air_example.conf示例配置如下，可以根据自己的需要修改。
```sh
# [Air](https://github.com/cosmtrek/air) TOML 格式的配置文件

# 工作目录
# 使用 . 或绝对路径，请注意 `tmp_dir` 目录必须在 `root` 目录下
root = "."
tmp_dir = "tmp"

[build]
# 只需要写你平常编译使用的shell命令。你也可以使用 `make`
# Windows平台示例: cmd = "go build -o tmp\main.exe ."
cmd = "go build -o ./tmp/main ."
# 由`cmd`命令得到的二进制文件名
# Windows平台示例：bin = "tmp\main.exe"
bin = "tmp/main"
# 自定义执行程序的命令，可以添加额外的编译标识例如添加 GIN_MODE=release
# Windows平台示例：full_bin = "tmp\main.exe"
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
# 监听以下文件扩展名的文件.
include_ext = ["go", "tpl", "tmpl", "html"]
# 忽略这些文件扩展名或目录
exclude_dir = ["assets", "tmp", "vendor", "frontend/node_modules"]
# 监听以下指定目录的文件
include_dir = []
# 排除以下文件
exclude_file = []
# 如果文件更改过于频繁，则没有必要在每次更改时都触发构建。可以设置触发构建的延迟时间
delay = 1000 # ms
# 发生构建错误时，停止运行旧的二进制文件。
stop_on_error = true
# air的日志文件名，该日志文件放置在你的`tmp_dir`中
log = "air_errors.log"

[log]
# 显示日志时间
time = true

[color]
# 自定义每个部分显示的颜色。如果找不到颜色，使用原始的应用程序日志。
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[misc]
# 退出时删除tmp目录
clean_on_exit = true
```

