

go-ini 是第三方库
=============

文档：  
https://ini.unknwon.io/docs/intro/getting_started



使用前需要安装：
-------------
```sh
$ go get gopkg.in/ini.v1

#也可以使用 GitHub 上的仓库：
$ go get github.com/go-ini/ini
```



首先，创建一个my.ini配置文件：
==============
在 ini 文件中，每个键值对占用一行，中间使用=隔开。以#开头的内容为注释。ini 文件是以分区（section）组织的。

分区以[name]开始，在下一个分区前结束。所有分区前的内容属于默认分区，如my.ini文件中的app_name和log_level。
```sh
app_name = awesome web
# possible values: DEBUG, INFO, WARNING, ERROR, FATAL
log_level = DEBUG
[mysql]
ip = 127.0.0.1
port = 3306
user = dj
password = 123456
database = awesome
[redis]
ip = 127.0.0.1
port = 6381
```


使用go-ini读取配置文件的步骤如下：

- 首先调用ini.Load加载文件，得到配置对象cfg；
- 然后以分区名调用配置对象的Section方法得到对应的分区对象section，默认分区的名字为""，也可以使用ini.DefaultSection；
- 以键名调用分区对象的Key方法得到对应的配置项key对象；

由于文件中读取出来的都是字符串，key对象需根据类型调用对应的方法返回具体类型的值使用，如上面的String、MustInt方法。

```golang
package main

import (
  "fmt"
  "log"
  "gopkg.in/ini.v1"
)

func main() {
  cfg, err := ini.Load("my.ini")
  if err != nil {
    log.Fatal("Fail to read file: ", err)
  }
  fmt.Println("App Name:", cfg.Section("").Key("app_name").String())
  fmt.Println("Log Level:", cfg.Section("").Key("log_level").String())
  fmt.Println("MySQL IP:", cfg.Section("mysql").Key("ip").String())
  mysqlPort, err := cfg.Section("mysql").Key("port").Int()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("MySQL Port:", mysqlPort)
  fmt.Println("MySQL User:", cfg.Section("mysql").Key("user").String())
  fmt.Println("MySQL Password:", cfg.Section("mysql").Key("password").String())
  fmt.Println("MySQL Database:", cfg.Section("mysql").Key("database").String())
  fmt.Println("Redis IP:", cfg.Section("redis").Key("ip").String())
  redisPort, err := cfg.Section("redis").Key("port").Int()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Redis Port:", redisPort)
}
```

运行以下程序，得到输出：
```sh
App Name: awesome web
Log Level: DEBUG
MySQL IP: 127.0.0.1
MySQL Port: 3306
MySQL User: dj
MySQL Password: 123456
MySQL Database: awesome
Redis IP: 127.0.0.1
Redis Port: 6381
```
配置文件中存储的都是字符串，所以类型为字符串的配置项不会出现类型转换失败的，故String()方法只返回一个值。

但如果类型为Int/Uint/Float64这些时，转换可能失败。所以Int()/Uint()/Float64()返回一个值和一个错误。

要留意这种不一致！如果我们将配置中 redis 端口改成非法的数字 x6381，那么运行程序将报错：

> 2020/01/14 22:43:13 strconv.ParseInt: parsing "x6381": invalid syntax



MustX便捷方法
==============
如果每次取值都需要进行错误判断，那么代码写起来会非常繁琐。为此，go-ini也提供对应的MustType（Type 为Init/Uint/Float64等）方法，这个方法只返回一个值。

同时它接受可变参数，如果类型无法转换，取参数中第一个值返回，并且该参数设置为这个配置的值，下次调用返回这个值：
```golang
package main
import (
  "fmt"
  "log"
  "gopkg.in/ini.v1"
)
func main() {
  cfg, err := ini.Load("my.ini")
  if err != nil {
    log.Fatal("Fail to read file: ", err)
  }
  redisPort, err := cfg.Section("redis").Key("port").Int()
  if err != nil {
    fmt.Println("before must, get redis port error:", err)
  } else {
    fmt.Println("before must, get redis port:", redisPort)
  }
  fmt.Println("redis Port:", cfg.Section("redis").Key("port").MustInt(6381))
  redisPort, err = cfg.Section("redis").Key("port").Int()
  if err != nil {
    fmt.Println("after must, get redis port error:", err)
  } else {
    fmt.Println("after must, get redis port:", redisPort)
  }
}
```
配置文件还是 redis 端口为非数字 x6381 时的状态，运行程序：
```sh
before must, get redis port error: strconv.ParseInt: parsing "x6381": invalid syntax
redis Port: 6381
after must, get redis port: 6381
```

我们看到第一次调用Int返回错误，以 6381 为参数调用MustInt之后，再次调用Int，成功返回 6381。MustInt源码也比较简单：
```golang
// gopkg.in/ini.v1/key.go
func (k *Key) MustInt(defaultVal ...int) int {
  val, err := k.Int()
  if len(defaultVal) > 0 && err != nil {
    k.value = strconv.FormatInt(int64(defaultVal[0]), 10)
    return defaultVal[0]
  }
  return val
}
```


分区操作
==============
获取信息

在加载配置之后，可以通过Sections方法获取所有分区，SectionStrings()方法获取所有分区名。
```golang
sections := cfg.Sections()
names := cfg.SectionStrings()
fmt.Println("sections: ", sections)
fmt.Println("names: ", names)

//运行输出 3 个分区：
[DEFAULT mysql redis]
```


调用Section(name)获取名为name的分区，如果该分区不存在，则自动创建一个分区返回：
```golang
newSection := cfg.Section("new")
fmt.Println("new section: ", newSection)
fmt.Println("names: ", cfg.SectionStrings())

//创建之后调用SectionStrings方法，新分区也会返回：
names:  [DEFAULT mysql redis new]
```


也可以手动创建一个新分区，如果分区已存在，则返回错误：
```
err := cfg.NewSection("new")
```


占位符
================
在配置文件中，可以使用占位符%(name)s表示用之前已定义的键name的值来替换，这里的s表示值为字符串类型：
```sh
NAME = ini
VERSION = v1
IMPORT_PATH = gopkg.in/%(NAME)s.%(VERSION)s
[package]
CLONE_URL = https://%(IMPORT_PATH)s
[package.sub]
```
上面在默认分区中设置IMPORT_PATH的值时，使用了前面定义的NAME和VERSION。

在package分区中设置CLONE_URL的值时，使用了默认分区中定义的IMPORT_PATH。




父子分区
================
我们还可以在分区名中使用.表示两个或多个分区之间的父子关系，例如package.sub的父分区为package，package的父分区为默认分区。

如果某个键在子分区中不存在，则会在它的父分区中再次查找，直到没有父分区为止：
```golang
cfg, err := ini.Load("parent_child.ini")
if err != nil {
  fmt.Println("Fail to read file: ", err)
  return
}
fmt.Println("Clone url from package.sub:", cfg.Section("package.sub").Key("CLONE_URL").String())

//运行程序输出：
Clone url from package.sub: https://gopkg.in/ini.v1
//子分区中package.sub中没有键CLONE_URL，返回了父分区package中的值。
```



加载多个配置文件(数据覆写)
================
在加载多个数据源时，如果某一个键在一个或多个数据源中出现，则会出现数据覆写。该键从前一个数据源读取的值会被下一个数据源覆写。

举例来说，如果加载两个配置文件 my.ini 和 my.ini.local（ 开始使用 中的输入和输出文件），app_mode 的值会是 production 而不是 development。
```sh
#my.ini
app_mode = development

#my.ini.local
app_mode = production
```

```golang
cfg, err := ini.Load("my.ini", "my.ini.local")
//...

cfg.Section("").Key("app_mode").String() // production
```
数据覆写只有在一种情况下不会触发，即使用 ShadowLoad 加载数据源。



保存配置
================
有时候，我们需要将生成的配置写到文件中。例如在写工具的时候。保存有两种类型的接口，一种直接保存到文件，另一种写入到io.Writer中：
```golang
err = cfg.SaveTo("my.ini")
err = cfg.SaveToIndent("my.ini", "\t")
cfg.WriteTo(writer)
cfg.WriteToIndent(writer, "\t")
```

下面我们通过程序生成前面使用的配置文件my.ini并保存：
```golang
package main
import (
  "fmt"
  "os"
  "gopkg.in/ini.v1"
)
func main() {
  cfg := ini.Empty()
  defaultSection := cfg.Section("")
  defaultSection.NewKey("app_name", "awesome web")
  defaultSection.NewKey("log_level", "DEBUG")
  mysqlSection, err := cfg.NewSection("mysql")
  if err != nil {
    fmt.Println("new mysql section failed:", err)
    return
  }
  mysqlSection.NewKey("ip", "127.0.0.1")
  mysqlSection.NewKey("port", "3306")
  mysqlSection.NewKey("user", "root")
  mysqlSection.NewKey("password", "123456")
  mysqlSection.NewKey("database", "awesome")
  redisSection, err := cfg.NewSection("redis")
  if err != nil {
    fmt.Println("new redis section failed:", err)
    return
  }
  redisSection.NewKey("ip", "127.0.0.1")
  redisSection.NewKey("port", "6381")
  err = cfg.SaveTo("my.ini")
  if err != nil {
    fmt.Println("SaveTo failed: ", err)
  }
  err = cfg.SaveToIndent("my-pretty.ini", "\t")
  if err != nil {
    fmt.Println("SaveToIndent failed: ", err)
  }
  cfg.WriteTo(os.Stdout)
  fmt.Println()
  cfg.WriteToIndent(os.Stdout, "\t")
}
```

运行程序，生成两个文件my.ini和my-pretty.ini，同时控制台输出文件内容。

my.ini：
```sh
app_name  = awesome web
log_level = DEBUG
[mysql]
ip       = 127.0.0.1
port     = 3306
user     = root
password = 123456
database = awesome
[redis]
ip   = 127.0.0.1
port = 6381
```

my-pretty.ini：
```sh
app_name  = awesome web
log_level = DEBUG
[mysql]
	ip       = 127.0.0.1
	port     = 3306
	user     = root
	password = 123456
	database = awesome
[redis]
	ip   = 127.0.0.1
	port = 6381
```
XXIndent方法会对子分区下的键增加缩进，看起来美观一点。




分区与结构体字段映射
===============
定义结构变量，加载完配置文件后，调用MapTo将配置项赋值到结构变量的对应字段中。
```golang
package main

import (
  "fmt"
  "gopkg.in/ini.v1"
)

type Config struct {
  AppName   string `ini:"app_name"`
  LogLevel  string `ini:"log_level"`
  MySQL     MySQLConfig `ini:"mysql"`
  Redis     RedisConfig `ini:"redis"`
}

type MySQLConfig struct {
  IP        string `ini:"ip"`
  Port      int `ini:"port"`
  User      string `ini:"user"`
  Password  string `ini:"password"`
  Database  string `ini:"database"`
}

type RedisConfig struct {
  IP      string `ini:"ip"`
  Port    int `ini:"port"`
}

func main() {
  cfg, err := ini.Load("my.ini")
  if err != nil {
    fmt.Println("load my.ini failed: ", err)
  }
  c := Config{}
  cfg.MapTo(&c)
  fmt.Println(c)
}
```
MapTo内部使用了反射，所以结构体字段必须都是导出的。如果键名与字段名不相同，那么需要在结构标签中指定对应的键名。

这一点与 Go 标准库encoding/json和encoding/xml不同。标准库json/xml解析时可以将键名app_name对应到字段名AppName。

或许这是go-ini库可以优化的点？


先加载，再映射有点繁琐，直接使用ini.MapTo将两步合并：
```golang
err = ini.MapTo(&c, "my.ini")
```


也可以只映射一个分区：
```golang
mysqlCfg := MySQLConfig{}
err = cfg.Section("mysql").MapTo(&mysqlCfg)
```


还可以通过结构体生成配置：
```golang
cfg := ini.Empty()
c := Config {
  AppName: 	"awesome web",
  LogLevel: 	"DEBUG",
  MySQL: MySQLConfig {
    IP: 	"127.0.0.1",
    Port:	3306,
    User:	"root",
    Password:"123456",
    Database:"awesome",
  },
  Redis: RedisConfig {
    IP:		"127.0.0.1",
    Port:	6381,
  },
}
err := ini.ReflectFrom(cfg, &c)
if err != nil {
  fmt.Println("ReflectFrom failed: ", err)
  return
}
err = cfg.SaveTo("my-copy.ini")
if err != nil {
  fmt.Println("SaveTo failed: ", err)
  return
}
```


