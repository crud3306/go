

```golang
package main
 
import (
    "github.com/gin-gonic/gin"
 
    consulapi "github.com/hashicorp/consul/api"
    "net/http"
    "fmt"
    "log"
)
 
func main() {
    r := gin.Default()
 
    // consul健康检查回调函数
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "ok",
        })
    })
 
 
    // 注册服务到consul
    ConsulRegister()
 
    // 从consul中发现服务
    ConsulFindServer()
 
    // 取消consul注册的服务
    //ConsulDeRegister()
 
 
    http.ListenAndServe(":8081", r)
}
 
/**
 consulapi.DefaultConfig()的源代码显示默认采用的是http方式连接"127.0.0.1:8500"，consul开发模式默认提供的http服务是在127.0.0.1:8500，在实际使用中需要设置为实际的参数。

    consulapi.AgentServiceCheck中的HTTP指定了健康检查的接口地址，健康检查还有其他几种方式，具体可以参考官方文档。

　　consulapi.AgentServiceCheck中的DeregisterCriticalServiceAfter指定检查不通过后多长时间注销本服务，这里设置为30秒。
 */
// 注册服务到consul
func ConsulRegister()  {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = "xxx.xx.xx.xx:8500"
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 创建注册到consul的服务
    registration := new(consulapi.AgentServiceRegistration)
    registration.ID = "111"
    registration.Name = "go-consul-test"
    registration.Port = 8081
    registration.Tags = []string{"go-consul-test"}
    registration.Address = "xx1.xx2.xx.xx"
 
    // 增加consul健康检查回调函数
    check := new(consulapi.AgentServiceCheck)
    check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
    check.Timeout = "5s"
    check.Interval = "5s"
    check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
    registration.Check = check
 
    // 注册服务到consul
    err = client.Agent().ServiceRegister(registration)
}
 
// 取消consul注册的服务
func ConsulDeRegister()  {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = "xxx.xx.xx.xx:8500"
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 参数是registration.ID
    client.Agent().ServiceDeregister("111")
}
 
// 从consul中发现服务
func ConsulFindServer()  {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = "xxx.xx.xx.xx:8500"
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 获取所有service
    services, _ := client.Agent().Services()
    for _, value := range services{
        fmt.Println(value.Address)
        fmt.Println(value.Port)
    }
 
    fmt.Println("=================================")

    // 仅获取指定service
    service, _, err := client.Agent().Service("111", nil)
    if err == nil{
        fmt.Println(service.Address)
        fmt.Println(service.Port)
    }
 
}
 
func ConsulCheckHeath()  {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = "xxx.xx.xx.xx:8500"
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 健康检查
    a, b, _ := client.Agent().AgentHealthServiceByID("111")
    fmt.Println(a)
    fmt.Println(b)
}
 
func ConsulKVTest()  {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = "xxx.xx.xx.xx:8500"
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // KV, put值
    values := "test"
    key := "go-consul-test/xxx.xx.xx.xx:8100"
    client.KV().Put(&consulapi.KVPair{Key:key, Flags:0, Value: []byte(values)}, nil)
 
    // KV get值
    data, _, _ := client.KV().Get(key, nil)
    fmt.Println(string(data.Value))
 
    // KV list
    datas, _ , _:= client.KV().List("go", nil)
    for _ , value := range datas{
        fmt.Println(value)
    }
    keys, _ , _ := client.KV().Keys("go", "", nil)
    fmt.Println(keys)
}
```



另外一个grpc的例子  
https://github.com/nosixtools/LearnGrpc 注意：发现时用的google.golang.org/grpc/resolver

