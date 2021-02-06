viper


简介
------------
viper 是一个配置解决方案，拥有丰富的特性：

- 支持 JSON/TOML/YAML/HCL/envfile/Java properties 等多种格式的配置文件；
- 可以设置监听配置文件的修改，修改时自动加载新的配置；
- 从环境变量、命令行选项和io.Reader中读取配置；
- 从远程配置系统中读取和监听修改，如 etcd/Consul；
- 代码逻辑中显示设置键值。



快速使用
============
安装：
> go get github.com/spf13/viper


使用：
------------
配置文件 config.toml
```sh
app_name = "awesome web"

# possible values: DEBUG, INFO, WARNING, ERROR, FATAL
log_level = "DEBUG"

[mysql]
ip = "127.0.0.1"
port = 3306
user = "dj"
password = 123456
database = "awesome"

[redis]
ip = "127.0.0.1"
port = 7381
```

```golang
package main

import (
  "fmt"
  "log"

  "github.com/spf13/viper"
)

func main() {
  viper.SetConfigName("config")
  viper.SetConfigType("toml")
  viper.AddConfigPath(".")
  viper.SetDefault("redis.port", 6381)
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatal("read config failed: %v", err)
  }

  fmt.Println(viper.Get("app_name"))
  fmt.Println(viper.Get("log_level"))

  fmt.Println("mysql ip: ", viper.Get("mysql.ip"))
  fmt.Println("mysql port: ", viper.Get("mysql.port"))
  fmt.Println("mysql user: ", viper.Get("mysql.user"))
  fmt.Println("mysql password: ", viper.Get("mysql.password"))
  fmt.Println("mysql database: ", viper.Get("mysql.database"))

  fmt.Println("redis ip: ", viper.Get("redis.ip"))
  fmt.Println("redis port: ", viper.Get("redis.port"))
}
```

viper 的使用非常简单，它需要很少的设置。
设置文件名（SetConfigName）、配置类型（SetConfigType）和搜索路径（AddConfigPath），然后调用ReadInConfig。

viper会自动根据类型来读取配置。使用时调用viper.Get方法获取键值。

编译、运行程序：
```sh
awesome web
DEBUG
mysql ip:  127.0.0.1
mysql port:  3306
mysql user:  dj
mysql password:  123456
mysql database:  awesome
redis ip:  127.0.0.1
redis port:  7381
```

有几点需要注意：

- 设置文件名时不要带后缀；
- 搜索路径可以设置多个，viper 会根据设置顺序依次查找；
- viper 获取值时使用section.key的形式，即传入嵌套的键名；
- 默认值可以调用viper.SetDefault设置。



读取键
------------
viper 提供了多种形式的读取方法。在上面的例子中，我们看到了Get方法的用法。Get方法返回一个interface{}的值，使用有所不便。

GetType系列方法可以返回指定类型的值。
其中，Type 可以为Bool/Float64/Int/String/Time/Duration/IntSlice/StringSlice。
但是请注意，如果指定的键不存在或类型不正确，GetType方法返回对应类型的零值。

如果要判断某个键是否存在，使用IsSet方法。
另外，GetStringMap和GetStringMapString直接以 map 返回某个键下面所有的键值对，前者返回map[string]interface{}，后者返回map[string]string。
AllSettings以map[string]interface{}返回所有设置。
```golang
// 省略包名和 import 部分

func main() {
  viper.SetConfigName("config")
  viper.SetConfigType("toml")
  viper.AddConfigPath(".")
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatal("read config failed: %v", err)
  }

  fmt.Println("protocols: ", viper.GetStringSlice("server.protocols"))
  fmt.Println("ports: ", viper.GetIntSlice("server.ports"))
  fmt.Println("timeout: ", viper.GetDuration("server.timeout"))

  fmt.Println("mysql ip: ", viper.GetString("mysql.ip"))
  fmt.Println("mysql port: ", viper.GetInt("mysql.port"))

  if viper.IsSet("redis.port") {
    fmt.Println("redis.port is set")
  } else {
    fmt.Println("redis.port is not set")
  }

  fmt.Println("mysql settings: ", viper.GetStringMap("mysql"))
  fmt.Println("redis settings: ", viper.GetStringMap("redis"))
  fmt.Println("all settings: ", viper.AllSettings())
}
```

我们在配置文件 config.toml 中添加protocols和ports配置：
```sh
[server]
protocols = ["http", "https", "port"]
ports = [10000, 10001, 10002]
timeout = 3s
```
编译、运行程序，输出：
```sh
protocols:  [http https port]
ports:  [10000 10001 10002]
timeout:  3s
mysql ip:  127.0.0.1
mysql port:  3306
redis.port is set
mysql settings:  map[database:awesome ip:127.0.0.1 password:123456 port:3306 user:dj]
redis settings:  map[ip:127.0.0.1 port:7381]
all settings:  map[app_name:awesome web log_level:DEBUG mysql:map[database:awesome ip:127.0.0.1 password:123456 port:3306 user:dj] redis:map[ip:127.0.0.1 port:7381] server:map[ports:[10000 10001 10002] protocols:[http https port]]]
```
如果将配置中的redis.port注释掉，将输出redis.port is not set。

上面的示例中还演示了如何使用time.Duration类型，只要是time.ParseDuration接受的格式都可以，例如3s、2min、1min30s等。



设置键值
------------
viper 支持在多个地方设置，使用下面的顺序依次读取：

- 调用Set显示设置的；
- 命令行选项；
- 环境变量；
- 配置文件；
- 默认值。

viper.Set  
如果某个键通过viper.Set设置了值，那么这个值的优先级最高。

viper.Set("redis.port", 5381)  
如果将上面这行代码放到程序中，运行程序，输出的redis.port将是 5381。



命令行选项  
------------
如果一个键没有通过viper.Set显示设置值，那么获取时将尝试从命令行选项中读取。
如果有，优先使用。viper 使用 pflag 库来解析选项。
我们首先在init方法中定义选项，并且调用viper.BindPFlags绑定选项到配置中：
```golang
func init() {
  pflag.Int("redis.port", 8381, "Redis port to connect")

  // 绑定命令行
  viper.BindPFlags(pflag.CommandLine)
}
```
然后，在main方法开头处调用pflag.Parse解析选项。

编译、运行程序：
```sh
$ ./main.exe --redis.port 9381
awesome web
DEBUG
mysql ip:  127.0.0.1
mysql port:  3306
mysql user:  dj
mysql password:  123456
mysql database:  awesome
redis ip:  127.0.0.1
redis port:  9381
```

如何不传入选项：
```sh
$ ./main.exe
awesome web
DEBUG
mysql ip:  127.0.0.1
mysql port:  3306
mysql user:  dj
mysql password:  123456
mysql database:  awesome
redis ip:  127.0.0.1
redis port:  7381
```
注意，这里并不会使用选项redis.port的默认值。

但是，如果通过下面的方法都无法获得键值，那么返回选项默认值（如果有）。试试注释掉配置文件中redis.port看看效果。



环境变量
------------
如果前面都没有获取到键值，将尝试从环境变量中读取。我们既可以一个个绑定，也可以自动全部绑定。

在init方法中调用AutomaticEnv方法绑定全部环境变量：
```golang
func init() {
  // 绑定环境变量
  viper.AutomaticEnv()
}
为了验证是否绑定成功，我们在main方法中将环境变量 GOPATH 打印出来：

func main() {
  // 省略部分代码

  fmt.Println("GOPATH: ", viper.Get("GOPATH"))
}
```

通过 系统 -> 高级设置 -> 新建 创建一个名为redis.port的环境变量，值为 10381。
运行程序，输出的redis.port值为 10381，并且输出中有 GOPATH 信息。

也可以单独绑定环境变量：
```golang
func init() {
  // 绑定环境变量
  viper.BindEnv("redis.port")
  viper.BindEnv("go.path", "GOPATH")
}

func main() {
  // 省略部分代码
  fmt.Println("go path: ", viper.Get("go.path"))
}
```
调用BindEnv方法，如果只传入一个参数，则这个参数既表示键名，又表示环境变量名。
如果传入两个参数，则第一个参数表示键名，第二个参数表示环境变量名。

还可以通过viper.SetEnvPrefix方法设置环境变量前缀，这样一来，通过AutomaticEnv和一个参数的BindEnv绑定的环境变量，
在使用Get的时候，viper 会自动加上这个前缀再从环境变量中查找。

如果对应的环境变量不存在，viper 会自动将键名全部转为大写再查找一次。所以，使用键名gopath也能读取环境变量GOPATH的值。



配置文件
------------
如果经过前面的途径都没能找到该键，viper 接下来会尝试从配置文件中查找。
为了避免环境变量的影响，需要删除redis.port这个环境变量。

看快速使用中的示例。




读取配置
=============

从io.Reader中读取
-------------
viper 支持从io.Reader中读取配置。这种形式很灵活，来源可以是文件，也可以是程序中生成的字符串，甚至可以从网络连接中读取的字节流。
```golang
package main

import (
  "bytes"
  "fmt"
  "log"

  "github.com/spf13/viper"
)

func main() {
  viper.SetConfigType("toml")
  tomlConfig := []byte(`
app_name = "awesome web"

# possible values: DEBUG, INFO, WARNING, ERROR, FATAL
log_level = "DEBUG"

[mysql]
ip = "127.0.0.1"
port = 3306
user = "dj"
password = 123456
database = "awesome"

[redis]
ip = "127.0.0.1"
port = 7381
`)
  err := viper.ReadConfig(bytes.NewBuffer(tomlConfig))
  if err != nil {
    log.Fatal("read config failed: %v", err)
  }

  fmt.Println("redis port: ", viper.GetInt("redis.port"))
}
```


Unmarshal 到一个结构体
================
viper 支持将配置Unmarshal到一个结构体中，为结构体中的对应字段赋值。
```golang
package main

import (
  "fmt"
  "log"

  "github.com/spf13/viper"
)

type Config struct {
  AppName  string
  LogLevel string

  MySQL    MySQLConfig
  Redis    RedisConfig
}

type MySQLConfig struct {
  IP       string
  Port     int
  User     string
  Password string
  Database string
}

type RedisConfig struct {
  IP   string
  Port int
}

func main() {
  viper.SetConfigName("config")
  viper.SetConfigType("toml")
  viper.AddConfigPath(".")
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatal("read config failed: %v", err)
  }

  var c Config
  viper.Unmarshal(&c)

  fmt.Println(c.MySQL)
}
```
编译，运行程序，输出：
```sh
{127.0.0.1 3306 dj 123456 awesome}
```




监听文件修改
================
viper 可以监听文件修改，热加载配置。因此不需要重启服务器，就能让配置生效。
```golang
package main

import (
  "fmt"
  "log"
  "time"

  "github.com/spf13/viper"
)

func main() {
  viper.SetConfigName("config")
  viper.SetConfigType("toml")
  viper.AddConfigPath(".")
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatal("read config failed: %v", err)
  }

  // 监听文件修改
  viper.WatchConfig()
  viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Printf("Config file:%s Op:%s\n", e.Name, e.Op)
  })

  fmt.Println("redis port before sleep: ", viper.Get("redis.port"))
  time.Sleep(time.Second * 10)
  fmt.Println("redis port after sleep: ", viper.Get("redis.port"))
}
```
只需要调用viper.WatchConfig，viper 会自动监听配置修改。如果有修改，重新加载的配置。

上面程序中，我们先打印redis.port的值，然后Sleep 10s。在这期间修改配置中redis.port的值，Sleep结束后再次打印。
发现打印出修改后的值：
```sh
redis port before sleep:  7381
redis port after sleep:  73810
```

另外，还可以为配置修改增加一个回调：
```golang
viper.OnConfigChange(func(e fsnotify.Event) {
  fmt.Printf("Config file:%s Op:%s\n", e.Name, e.Op)
})
```
这样文件修改时会执行这个回调。


```golang
package main

import (
  "fmt"
  "net/http"

  "github.com/fsnotify/fsnotify"

  "github.com/gin-gonic/gin"
  "github.com/spf13/viper"
)

type Config struct {
  Port    int    `mapstructure:"port"`
  Version string `mapstructure:"version"`
}

var Conf = new(Config)

func main() {
  viper.SetConfigFile("./conf/config.yaml") // 指定配置文件路径
  err := viper.ReadInConfig()               // 读取配置信息
  if err != nil {                           // 读取配置信息失败
    panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }
  // 将读取的配置信息保存至全局变量Conf
  if err := viper.Unmarshal(Conf); err != nil {
    panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
  }
  
  // 监控配置文件变化
  viper.WatchConfig()
  // 注意！！！配置文件发生变化后要同步到全局变量Conf
  viper.OnConfigChange(func(in fsnotify.Event) {
    fmt.Println("夭寿啦~配置文件被人修改啦...")
    if err := viper.Unmarshal(Conf); err != nil {
      panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
    }
  })

  r := gin.Default()
  // 访问/version的返回值会随配置文件的变化而变化
  r.GET("/version", func(c *gin.Context) {
    c.String(http.StatusOK, Conf.Version)
  })

  if err := r.Run(fmt.Sprintf(":%d", Conf.Port)); err != nil {
    panic(err)
  }
}
```






保存配置
================
有时候，我们想要将程序中生成的配置，或者所做的修改保存下来。viper 提供了接口！

- WriteConfig：将当前的 viper 配置写到预定义路径，如果没有预定义路径，返回错误。将会覆盖当前配置；
- SafeWriteConfig：与上面功能一样，但是如果配置文件存在，则不覆盖；
- WriteConfigAs：保存配置到指定路径，如果文件存在，则覆盖；
- SafeWriteConfig：与上面功能一样，但是入股配置文件存在，则不覆盖。

下面我们通过程序生成一个config.toml配置：
```golang
package main

import (
  "log"

  "github.com/spf13/viper"
)

func main() {
  viper.SetConfigName("config")
  viper.SetConfigType("toml")
  viper.AddConfigPath(".")

  viper.Set("app_name", "awesome web")
  viper.Set("log_level", "DEBUG")
  viper.Set("mysql.ip", "127.0.0.1")
  viper.Set("mysql.port", 3306)
  viper.Set("mysql.user", "root")
  viper.Set("mysql.password", "123456")
  viper.Set("mysql.database", "awesome")

  viper.Set("redis.ip", "127.0.0.1")
  viper.Set("redis.port", 6381)

  err := viper.SafeWriteConfig()
  if err != nil {
    log.Fatal("write config failed: ", err)
  }
}
```

编译、运行程序，生成的文件如下：
```sh
app_name = "awesome web"
log_level = "DEBUG"

[mysql]
  database = "awesome"
  ip = "127.0.0.1"
  password = "123456"
  port = 3306
  user = "root"

[redis]
  ip = "127.0.0.1"
  port = 6381
```

