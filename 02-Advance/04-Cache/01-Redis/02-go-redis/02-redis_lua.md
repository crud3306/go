

Lua+Redis

秒杀项目的高并发的核心是lua和redis。lua脚本把用户请求的处理过程都包装起来，作为一个原子操作。lua脚本是用C语言写的，体积小、运行速度快，作为一个原子事务来执行。


使用Lua脚本的好处
--------------
- 减少网络开销：可以将多个请求通过脚本的形式一次发送，减少网络时延和请求次数。比如说在抢购事务中要访问两次Redis，第一次读Redis看是否还有优惠券可以抢，第二次抢购成功写Redis使优惠券数量减一，原本是两个请求，使用lua后将他们整合起来一起发送，就变成只有一次请求了。

- 原子性的操作：Redis会将整个脚本作为一个整体执行，中间不会被其他命令插入。因此在编写脚本的过程中无需担心会出现竞态条件（两个线程竞争同一资源时，如果对资源的访问顺序敏感，就称存在竞态条件）。

- 代码复用：客户端发送的脚本会一直存在redis中，这样，其他客户端可以复用这一脚本来完成相同的逻辑。
体积小，加载快，运行快。


使用Lua的例子
--------------
秒杀项目中使用Lua脚本，最重要的原因是为了他的原子性。redis本身是单线程应用，就是说读或者写都是各自独立的，是串行的，每次只能执行一个操作，而除读写之外的其他操作是可以并行的。在这个项目里，用户抢优惠券的请求是先读取redis查看优惠券信息，然后再抢优惠券（就是写redis让优惠券数量减一）。就是先读后写。

正常情况下是没什么问题的，但是在高并发的访问情况下，就可能会出现临界区的问题（也叫竞态条件）。比如说，假设当前Redis中优惠券的数量为1，此时两个用户A和B几乎同时发来两个抢购请求，A先读Redis，得知还有一张优惠券，因此允许抢购，此时原本应该是A继续写Redis，抢走优惠券使优惠券数量变0的。但是意外出现了，此时B插入其中，B抢占了Redis的访问权限，B读Redis发现此时还有一张优惠券，于是他也被允许抢购。那么问题就出现了，无论后面是A还是B先拿回Redis的写权限，都会出现超卖的问题，因为优惠券实际只有一张，却有两个用户抢购成功了。

为了避免抢购事务中redis的读写分离导致超卖的问题，就需要用到lua脚本，把用户抢优惠券这个过程的读操作和写操作捆绑起来，作为一个整体的原子操作，就是说一个用户读redis后要马上写redis，都处理完了，才允许其他用户再来读。这样的话，就相当于是串行处理用户抢优惠券的请求了，所以说这里高并发的实现其实主要是用redis来加快速度，然后lua来保证不会出现临界区问题的错误。




redis 调用Lua脚本
================
redis确保一条Lua script脚本执行期间，其它任何脚本或者命令都无法执行，这样保证了脚本执行的原子性。

从Redis 2.6 版本开始，内嵌支持 Lua 环境。通过使用EVAL或EVALSHA命令可以使用 Lua 解释器来执行脚本。 EVAL和EVALSHA的使用是差不多的。

- EVAL命令
- SCRIPT命令



EVAL命令
================
redis调用Lua脚本需要使用EVAL命令。

redis EVAL命令格式：

> EVAL script numkeys key [key ...] arg [arg ...]  

script： 参数是一段 Lua 5.1 脚本程序。脚本不必(也不应该)定义为一个 Lua 函数。  
numkeys： 用于指定键名参数的个数。   
key [key ...]： 从 EVAL 的第三个参数开始算起，表示在脚本中所用到的那些 Redis 键(key)，这些键名参数可以在 Lua 中通过全局变量 KEYS 数组，用 1 为基址的形式访问( KEYS[1] ， KEYS[2] ，以此类推)。  
arg [arg ...]： 附加参数，在 Lua 中通过全局变量 ARGV 数组访问，访问的形式和 KEYS 变量类似( ARGV[1] 、 ARGV[2] ，诸如此类)。  

```sh
127.0.0.1:6379> eval "return {KEYS[1],ARGV[1],ARGV[2]}" 1 1 ONE TWO
1) "1"
2) "ONE"
3) "TWO"
```


最简单的例子：
```sh
127.0.0.1:6379> eval "return {'Hello, GrassInWind!'}" 0
1) "Hello, GrassInWind!"
127.0.0.1:6379> eval "return redis.call('set',KEYS[1],'bar')" 1 foo
OK
```
使用redis-cli调用lua脚本示例(若在windows系统下，则需要在git bash中执行，在powershell中无法读取value)：
```sh
***@LAPTOP-V7V47H0L MINGW64 /d/study/code/lua
$ redis-cli.exe -a 123 --eval test.lua  testkey , hello
hello
```


test.lua如下(redis log打印在server的日志中)：
```sh
local key,value = KEYS[1],ARGV[1]
redis.log(redis.LOG_NOTICE, "key=", key, "value=", value)
redis.call('SET', key, value)
local a = redis.call('GET', key)
return a
```



SCRIPT命令
================
redis提供了以下几个script命令，用于对于脚本子系统进行控制：

script flush：清除所有的脚本缓存

script load：将脚本装入脚本缓存，不立即运行并返回其校验和

script exists：根据指定脚本校验和，检查脚本是否存在于缓存

script kill：杀死当前正在运行的脚本（防止脚本运行缓存，占用内存）


主要优势： 
----------------
- 减少网络开销：多个请求通过脚本一次发送，减少网络延迟

- 原子操作：将脚本作为一个整体执行，中间不会插入其他命令，无需使用事务

- 复用：客户端发送的脚本永久存在redis中，其他客户端可以复用脚本

- 可嵌入性：可嵌入JAVA，C#等多种编程语言，支持不同操作系统跨平台交互


通过script命令加载及执行lua脚本示例：
```sh
127.0.0.1:6379> script load "return 'Hello GrassInWind'"
"c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b"
127.0.0.1:6379> script exists "c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b"
1) (integer) 1
127.0.0.1:6379> evalsha "c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b" 0
"Hello GrassInWind"
127.0.0.1:6379> script flush
OK
127.0.0.1:6379> script exists "c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b"
1) (integer) 0
```


简单例子
================
生成一段Lua脚本
```golang
// KEYS: key for record
// ARGV: fieldName, currentUnixTimestamp, recordTTL
// Update expire field of record key to current timestamp, and renew key expiration
var updateRecordExpireScript = redis.NewScript(`
redis.call("EXPIRE", KEYS[1], ARGV[3])
redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
return 1
`)

//该变量创建时，Lua代码不会被执行，也不需要有已存的Redis连接。
//Redis提供的Lua脚本支持，默认有KEYS、ARGV两个数组，KEYS代表脚本运行时传入的若干键值，ARGV代表传入的若干参数。由于Lua代码需要保持简洁，难免难以读懂，最好为这些参数写一些注释
//注意：上面一段代码使用``跨行，`所在的行虽然空白回车，也会被认为是一行，报错时不要看错代码行号。
```

运行一段Lua脚本
```golang
updateRecordExpireScript.Run(c.Client, []string{recordKey(key)}, 
                                    expireField,
                                    time.Now().UTC().UnixNano(), int64(c.opt.RecordTTL/time.Second)).Err()
```
运行时，Run将会先通过EVALSHA尝试通过缓存运行脚本。如果没有缓存，则使用EVAL运行，这时Lua脚本才会被整个传入Redis。


Lua脚本的限制
----------------
- Redis不提供引入额外的包，例如os等，只有redis这一个包可用。
- Lua脚本将会在一个函数中运行，所有变量必须使用local声明
- return返回多个值时，Redis将会只给你第一个


脚本中的类型限制
----------------
- 脚本返回nil时，Go中得到的是err = redis.Nil（与Get找不到值相同）
- 脚本返回false时，Go中得到的是nil，脚本返回true时，Go中得到的是int64类型的1
- 脚本返回{"ok": ...}时，Go中得到的是redis的status类型（true/false)
- 脚本返回{"err": ...}时，Go中得到的是err值，也可以通过return redis.error_reply("My Error")达成
- 脚本返回number类型时，Go中得到的是int64类型
- 传入脚本的KEYS/ARGV中的值一律为string类型，要转换为数字类型应当使用to_number


如果脚本运行了很久会发生什么？
----------------
Lua脚本运行期间，为了避免被其他操作污染数据，这期间将不能执行其它命令，一直等到执行完毕才可以继续执行其它请求。当Lua脚本执行时间超过了lua-time-limit时，其他请求将会收到Busy错误，除非这些请求是SCRIPT KILL（杀掉脚本）或者SHUTDOWN NOSAVE（不保存结果直接关闭Redis）





一个简单的秒杀例子
==============

逻辑：  
一个key 存储商品数量  
一个SET 存储已抢到的用户id  

```golang
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"sync"
)


func createScript() *redis.Script {
	script := redis.NewScript(`
		local goodsSurplus
		local flag
		local existUserIds    = tostring(KEYS[1])
		local memberUid       = tonumber(ARGV[1])
		local goodsSurplusKey = tostring(KEYS[2])

		-- 是否已抢到过
		local hasBuy = redis.call("sIsMember", existUserIds, memberUid)
		if hasBuy ~= 0 then
		  return 0
		end
		
		-- 获取剩余商品数量
		goodsSurplus =  redis.call("GET", goodsSurplusKey)
		if goodsSurplus == false then
		  return 0
		end
		
		-- 没有剩余可抢购物品
		goodsSurplus = tonumber(goodsSurplus)
		if goodsSurplus <= 0 then
		  return 0
		end
		
		-- 添加用户至已抢列表
		flag = redis.call("SADD", existUserIds, memberUid)
		-- 商品数量减1
		flag = redis.call("DECR", goodsSurplusKey)
		
		return 1
	`)
	return script
}


func evalScript(client *redis.Client, userId string, wg *sync.WaitGroup){
	defer wg.Done()

	script := createScript()
	sha, err := script.Load(client.Context(), client).Result()
	if err != nil {
		log.Fatalln(err)
	}


	ret := client.EvalSha(client.Context(), sha, []string{
		"hadBuyUids",
		"goodsSurplus",
	}, userId)
	if result, err := ret.Result();err!= nil {
		log.Fatalf("Execute Redis fail: %v", err.Error())
	} else {
		fmt.Println("")
		fmt.Printf("userid: %s, result: %d", userId, result)
	}
}

func main() {
	var wg sync.WaitGroup
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for _, v := range []string{"5824742984", "5824742984", "5824742983", "5824742983", "5824742982", "5824742980"}{
		wg.Add(1)
		go evalScript(client, v, &wg)
	}
	wg.Wait()

}
```




redis+lua 实现评分排行榜实时更新
=============
使用redis的zset保存排行数据，使用lua脚本实现评分排行更新的原子操作。

lua 脚本

相关redis命令：

ZCARD key 获取有序集合的成员数
ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT] 通过分数返回有序集合指定区间内的成员(从小到大的顺序)
ZREMRANGEBYRANK key start stop 移除有序集合中给定的排名区间的所有成员
ZADD key score1 member1 [score2 member2] 向有序集合添加一个或多个成员，或者更新已存在成员的分数

主要思路是维护一个zset，将评分前N位保存到redis中，当成员的评分发生变化时，动态更新zset的成员信息。

lua脚本如下，其中
```sh
KEYS[1] 表示zset的key，
ARGV[1] 表示期望的zset最大存储成员数量，
ARGV[2] 表示评分上限，默认评分下限是0，
ARGV[3] 表示待添加的评分，
ARGV[4] 表示待添加的成员名称。
```

```sh
-- redis zset operations
-- argv[capacity maxScore newMemberScore member]
-- 执行示例 redis-cli.exe --eval zsetop.lua mtest , 3 5 5 test1
-- 获取键和参数

local key,cap,maxSetScore,newMemberScore,member = KEYS[1],ARGV[1],ARGV[2],ARGV[3],ARGV[4]
redis.log(redis.LOG_NOTICE, "key=", key,",cap=", cap,",maxSetScore=", maxSetScore,",newMemberScore=", newMemberScore,",member=", member)
local len = redis.call('zcard', key);

-- len need not nil, otherwise will occur "attempt to compare nil with number"
if len then

    if tonumber(len) >= tonumber(cap)
    then
        local num = tonumber(len)-tonumber(cap)+1
        local list = redis.call('zrangebyscore',key,0,maxSetScore,'limit',0,num)
        redis.log(redis.LOG_NOTICE,"key=",key,"maxSetScore=",maxSetScore, "num=",num)

        for k,lowestScoreMember in pairs(list) do
            local lowestScore = redis.call('zscore', key,lowestScoreMember)
            redis.log(redis.LOG_NOTICE, "list: ", lowestScore, lowestScoreMember)

            if tonumber(newMemberScore) > tonumber(lowestScore)
            then
                local rank = redis.call('zrevrank',key,member)

                -- rank is nil indicate new member is not exist in set, need remove the lowest score member
                if not rank then
                    local index = tonumber(len) - tonumber(cap);
                    redis.call('zremrangebyrank',key, 0, index)
                end
                
                redis.call('zadd', key, newMemberScore, member);
                break
            end
        end
    else
        redis.call('zadd', key, newMemberScore, member);
    end
end
```


Golang调用redis+lua示例
----------------
init函数中读取Lua脚本并通过redisgo包的NewScript函数加载这个脚本，在使用时通过返回的指针调用lua.Do()即可。
```golang
func init() {
    ...
    file, err := os.Open(zsetopFileName)
    if err != nil {
        panic("open"+zsetopFileName+" "+err.Error())
    }
    bytes,err := ioutil.ReadAll(file)
    if err != nil {
        panic(err.Error())
    }
    zsetopScript = utils.UnsafeBytesToString(bytes)
    logs.Debug(zsetopScript)
    lua =redis.NewScript(1,zsetopScript)
}

func ZaddWithCap(key,member string, score float32, maxScore, cap int) (reply interface{}, err error) {
    c := pool.Get()
    //Do optimistically evaluates the script using the EVALSHA command. If script not exist, will use eval command.
    reply, err= lua.Do(c,key,cap,maxScore,score,member)
    return
}
```

redisgo包对Do方法做了优化，会检查这个脚本的SHA是否存在，若不存在，会通过EVAL命令执行即会加载脚本，下次执行就可以通过
EVALSHA来执行了。
```sh
func (s *Script) Do(c Conn, keysAndArgs ...interface{}) (interface{}, error) {
    v, err := c.Do("EVALSHA", s.args(s.hash, keysAndArgs)...)
    if e, ok := err.(Error); ok && strings.HasPrefix(string(e), "NOSCRIPT ") {
        v, err = c.Do("EVAL", s.args(s.src, keysAndArgs)...)
    }
    return v, err
}
```




