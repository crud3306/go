

redigo 库
==============


开源库redigo的使用
-------------
github地址：   
https://github.com/gomodule/redigo


获取：   
> go get github.com/gomodule/redigo/redis


连接redis
-------------
```golang
package main
import (
    "fmt"
    "github.com/gomodule/redigo/redis"
)

func main() {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()
}
```


读写 
----------

```golang
package main

import (
    "fmt"
    "time"

    "github.com/gomodule/redigo/redis"
)

func main() {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()


    // 不设过期
    key := "mykey"
    _, err = c.Do("SET", key, "superWang")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }

    username, err := redis.String(c.Do("GET", key))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        fmt.Printf("Get %s: %v \n", key, username)
    }


    // 设置过期时间
    exprieKey := "exprie_key"
    _, err = c.Do("SET", exprieKey, "superWang", "EX", "5")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }

    username, err := redis.String(c.Do("GET", exprieKey))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        fmt.Printf("Get %s: %v \n", exprieKey, username)
    }

    time.Sleep(8 * time.Second)

    username, err = redis.String(c.Do("GET", exprieKey))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        fmt.Printf("Get %s: %v \n", exprieKey, username)
    }
}


//输出： 
Get mykey: superWang 
redis get failed: redigo: nil returned
```


批量写入读取
-----------
MGET key [key …] 
MSET key value [key value …]


批量写入读取对象(Hashtable) 
-----------
HMSET key field value [field value …] 
HMGET key field [field …]


检测值是否存在 
-----------
EXISTS key

```golang
package main

import (
    "fmt"

    "github.com/gomodule/redigo/redis"
)

func main() {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()

    _, err = c.Do("SET", "mykey", "superWang")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }

    is_key_exit, err := redis.Bool(c.Do("EXISTS", "mykey1"))
    if err != nil {
        fmt.Println("error:", err)
    } else {
        fmt.Printf("exists or not: %v \n", is_key_exit)
    }

}

//输出： 
exists or not: false
```


删除 
------------
DEL key [key …]
```golang
package main

import (
    "fmt"

    "github.com/gomodule/redigo/redis"
)

func main() {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()

    _, err = c.Do("SET", "mykey", "superWang")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }

    username, err := redis.String(c.Do("GET", "mykey"))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        fmt.Printf("Get mykey: %v \n", username)
    }

    _, err = c.Do("DEL", "mykey")
    if err != nil {
        fmt.Println("redis delelte failed:", err)
    }

    username, err = redis.String(c.Do("GET", "mykey"))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        fmt.Printf("Get mykey: %v \n", username)
    }
}

//输出： 
Get mykey: superWang 
redis get failed: redigo: nil returned
```


读写json到redis
-----------
```golang
package main

import (
    "encoding/json"
    "fmt"

    "github.com/gomodule/redigo/redis"
)

func main() {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()

    key := "profile"
    imap := map[string]string{"username": "666", "phonenumber": "888"}
    value, _ := json.Marshal(imap)

    n, err := c.Do("SETNX", key, value)
    if err != nil {
        fmt.Println(err)
    }
    if n == int64(1) {
        fmt.Println("success")
    }

    var imapGet map[string]string

    valueGet, err := redis.Bytes(c.Do("GET", key))
    if err != nil {
        fmt.Println(err)
    }

    errShal := json.Unmarshal(valueGet, &imapGet)
    if errShal != nil {
        fmt.Println(err)
    }
    fmt.Println(imapGet["username"])
    fmt.Println(imapGet["phonenumber"])
}
```


设置过期时间 
-----------
EXPIRE key seconds
```golang
// 设置过期时间为24小时  
n, _ := rs.Do("EXPIRE", key, 24*3600)  
if n == int64(1) {  
    fmt.Println("success")  
}  
```



列表操作 
-------------
命令：
```
redis 127.0.0.1:6379> LPUSH runoobkey redis
(integer) 1
redis 127.0.0.1:6379> LPUSH runoobkey mongodb
(integer) 2
redis 127.0.0.1:6379> LPUSH runoobkey mysql
(integer) 3
redis 127.0.0.1:6379> LRANGE runoobkey 0 10

1) "mysql"
2) "mongodb"
3) "redis"
```

代码实现：
```golang
package main

import (
    "fmt"

    "github.com/gomodule/redigo/redis"
)

func main() {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()

    _, err = c.Do("lpush", "runoobkey", "redis")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }

    _, err = c.Do("lpush", "runoobkey", "mongodb")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }
    _, err = c.Do("lpush", "runoobkey", "mysql")
    if err != nil {
        fmt.Println("redis set failed:", err)
    }

    values, _ := redis.Values(c.Do("lrange", "runoobkey", "0", "100"))

    for _, v := range values {
        fmt.Println(string(v.([]byte)))
    }
}

//输出： 
mysql 
mongodb 
redis
```


管道
------------
请求/响应服务可以实现持续处理新请求，即使客户端没有准备好读取旧响应。这样客户端可以发送多个命令到服务器而无需等待响应，最后在一次读取多个响应。这就是管道化(pipelining)，这个技术在多年就被广泛使用了。距离，很多POP3协议实现已经支持此特性，显著加速了从服务器下载新邮件的过程。 
Redis很早就支持管道化，所以无论你使用任何版本，你都可以使用管道化技术

连接支持使用Send()，Flush()，Receive()方法支持管道化操作
```golang
Send(commandName string, args ...interface{}) error
Flush() error
Receive() (reply interface{}, err error)
```

Send向连接的输出缓冲中写入命令。Flush将连接的输出缓冲清空并写入服务器端。Recevie按照FIFO顺序依次读取服务器的响应。下例展示了一个简单的管道：
```golang
c.Send("SET", "foo", "bar")
c.Send("GET", "foo")
c.Flush()
c.Receive() // reply from SET
v, err = c.Receive() // reply from GET
```
Do方法组合了Send,Flush和 Receive方法。Do方法先写入命令，然后清空输出buffer，最后接收全部挂起响应包括Do方发出的命令的结果。如果任何响应中包含一个错误，Do返回错误。如果没有错误，Do方法返回最后一个响应。




非v8版
============

使用
------------
```golang
import "github.com/gomodule/redigo/redis"

func main() {

    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        // handle error
    }
    defer conn.Close()


    // base
     _, err = conn.Do("Set", "name", "tomjerry")
    if err != nil {
        fmt.Println("set err=", err)
        return
    }

    r, err := redis.String(conn.Do("Get", "name"))
    if err != nil {
        fmt.Println("set err=", err)
        return
    }



    // 操作 reids 中的 string 结构
    c.Do("SET", "hello", "world")
    s, _ := redis.String(c.Do("GET", "hello"))

    fmt.Printf("%#v\n", s)

    

    // 操作 hash结构，HMSET 一次写入多个属性值
    m := map[string]string{
        "title":  "Example2",
        "author": "Steve",
        "body":   "Map",
    }

    if _, err := c.Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...); err != nil {
        fmt.Println(err)
        return
    }



    // 管道
    c.Send("SET", "foo", "bar")
    c.Send("GET", "foo")
    c.Flush()
    c.Receive() // reply from SET
    v, err = c.Receive() // reply from GET
}
```


连接池
------------
```golang
package main
import (
    "fmt"
    redigo "github.com/gomodule/redigo/redis"
    "time"
)

func main() {
    var addr = "127.0.0.1:6379"
    var password = ""

    pool := PoolInitRedis(addr, password)

    c1 := pool.Get()
    c2 := pool.Get()
    c3 := pool.Get()
    c4 := pool.Get()
    c5 := pool.Get()
    fmt.Println(c,c2,c3,c4,c5)
    time.Sleep(time.Second * 5)//redis一共有多少个连接？？

    c1.Close()
    c2.Close()
    c3.Close()
    c4.Close()
    c5.Close()
    time.Sleep(time.Second*5) //redis一共有多少个连接？？


    //下次是怎么取出来的？？
    b1 := pool.Get()
    b2 := pool.Get()
    b3 := pool.Get()
    fmt.Println(b1,b2,b3)
    time.Sleep(time.Second*5)

    b1.Close()
    b2.Close()
    b3.Close()

    //redis目前一共有多少个连接？？
    for{
        fmt.Println("主程序运行中....")
        time.Sleep(time.Second*1) 
    }
}

// redis pool
func PoolInitRedis(server string, password string) *redigo.Pool {
    return &redigo.Pool{
        MaxIdle:     2,//空闲数
        IdleTimeout: 240 * time.Second,
        MaxActive:   3,//最大数
        Dial: func() (redigo.Conn, error) {
            c, err := redigo.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            if password != "" {
                if _, err := c.Do("AUTH", password); err != nil {
                    c.Close()
                    return nil, err
                }
            }
            return c, err
        },
        TestOnBorrow: func(c redigo.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}
```





v8版
============
> go get github.com/go-redis/redis/v8


使用
------------
```golang
import (
    "context"
    "github.com/go-redis/redis/v8"  
)

var ctx = context.Background()

func ExampleNewClient() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := rdb.Ping(ctx).Result()
    fmt.Println(pong, err)
    // Output: PONG <nil>
}


func ExampleClient() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }



    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)



    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exist
}
```