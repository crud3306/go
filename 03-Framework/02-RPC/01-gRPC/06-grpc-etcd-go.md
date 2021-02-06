
gRPC接入etcd
　　
本篇介绍gRPC接入etcd，实现服务注册与服务发现。


需要先安装Go语言的etcd客户端包：
----------------
> go get go.etcd.io/etcd/clientv3



目录结构：
$GOPATH/src/go-git/etcd-demo：
```sh
etcd-demo
	client
		main.go
	proto
		greet.pb.go
		greet.proto
	server
		main.go
	go.mod
	go.sum
```

 

一. 协议制定（proto/greet.proto）
================
```sh
syntax = "proto3";
 
option go_package = ".;greet";
 
service Greet {
    rpc Morning(GreetRequest)returns(GreetResponse){}
    rpc Night(GreetRequest)returns(GreetResponse){}
}
 
message GreetRequest {
    string name = 1;
}
 
message GreetResponse {
    string message = 1;
    string from = 2;
}
```

生成代码：（proto子目录下执行）
```sh
protoc --go_out=plugins=grpc:. *.proto
```

执行完成，proto子目录生成文件greet.pb.go。

 

二. 服务端（server/main.go）
================

服务端主要有以下步骤：

- 监听网络端口  
- 创建gRPC句柄，注册gRPC服务  
- 将服务地址注册到etcd，并持续发送心跳    
- 监听并处理服务请求  


这里主要介绍一下将服务地址注册到etcd的过程（双保险）：

- 一方面，由于服务端无法保证自身是一直可用的，所以与etcd的租约是有时间期限的，租约一旦过期，服务端存储在etcd上的服务地址信息就会消失。  
- 另一方面，服务端可用时又必须保证调用方能发现自己，即保证自己在etcd上的服务地址信息不消失，所以需要发送心跳检测，一旦发现etcd上没有自己的服务地址时，请求重新添加（续租）。


代码
```golang
package main
 
import (
    "flag"
    "fmt"
    proto "go-git/etcd-demo/proto"
    "net"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"
 
    "go.etcd.io/etcd/clientv3"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
)
 
const schema = "ns"
 
var host = "127.0.0.1" //服务器主机
var (
    Port        = flag.Int("Port", 3000, "listening port")                           //服务器监听端口
    ServiceName = flag.String("ServiceName", "greet_service", "service name")        //服务名称
    EtcdAddr    = flag.String("EtcdAddr", "127.0.0.1:2379", "register etcd address") //etcd的地址
)
var cli *clientv3.Client
 
//rpc服务接口
type greetServer struct{}
 
func (gs *greetServer) Morning(ctx context.Context, req *proto.GreetRequest) (*proto.GreetResponse, error) {
    fmt.Printf("Morning 调用: %s\n", req.Name)
    return &proto.GreetResponse{
        Message: "Good morning, " + req.Name,
        From:    fmt.Sprintf("127.0.0.1:%d", *Port),
    }, nil
}
 
func (gs *greetServer) Night(ctx context.Context, req *proto.GreetRequest) (*proto.GreetResponse, error) {
    fmt.Printf("Night 调用: %s\n", req.Name)
    return &proto.GreetResponse{
        Message: "Good night, " + req.Name,
        From:    fmt.Sprintf("127.0.0.1:%d", *Port),
    }, nil
}
 
//将服务地址注册到etcd中
func register(etcdAddr, serviceName, serverAddr string, ttl int64) error {
    var err error
 
    if cli == nil {
        //构建etcd client
        cli, err = clientv3.New(clientv3.Config{
            Endpoints:   strings.Split(etcdAddr, ";"),
            DialTimeout: 15 * time.Second,
        })
        if err != nil {
            fmt.Printf("连接etcd失败：%s\n", err)
            return err
        }
    }
 
    //与etcd建立长连接，并保证连接不断(心跳检测)
    ticker := time.NewTicker(time.Second * time.Duration(ttl))
    go func() {
        key := "/" + schema + "/" + serviceName + "/" + serverAddr
        for {
            resp, err := cli.Get(context.Background(), key)
            //fmt.Printf("resp:%+v\n", resp)
            if err != nil {
                fmt.Printf("获取服务地址失败：%s", err)

            } else if resp.Count == 0 { //尚未注册
                err = keepAlive(serviceName, serverAddr, ttl)
                if err != nil {
                    fmt.Printf("保持连接失败：%s", err)
                }
            }

            <-ticker.C
        }
    }()
 
    return nil
}
 
//保持服务器与etcd的长连接
func keepAlive(serviceName, serverAddr string, ttl int64) error {
    //创建租约
    leaseResp, err := cli.Grant(context.Background(), ttl)
    if err != nil {
        fmt.Printf("创建租期失败：%s\n", err)
        return err
    }
 
    //将服务地址注册到etcd中
    key := "/" + schema + "/" + serviceName + "/" + serverAddr
    _, err = cli.Put(context.Background(), key, serverAddr, clientv3.WithLease(leaseResp.ID))
    if err != nil {
        fmt.Printf("注册服务失败：%s", err)
        return err
    }
 
    //建立长连接
    ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
    if err != nil {
        fmt.Printf("建立长连接失败：%s\n", err)
        return err
    }
 
    //清空keepAlive返回的channel
    go func() {
        for {
            <-ch
        }
    }()
    return nil
}
 
//取消注册
func unRegister(serviceName, serverAddr string) {
    if cli != nil {
        key := "/" + schema + "/" + serviceName + "/" + serverAddr
        cli.Delete(context.Background(), key)
    }
}

func main() {
    flag.Parse()
 
    //监听网络
    serverAddr := fmt.Sprintf("%s:%d", host, *Port)
    listener, err := net.Listen("tcp", serverAddr)
    if err != nil {
        fmt.Println("监听网络失败：", err)
        return
    }
    defer listener.Close()
 
    //创建grpc句柄
    srv := grpc.NewServer()
    defer srv.GracefulStop()
 
    //将greetServer结构体注册到grpc服务中
    proto.RegisterGreetServer(srv, &greetServer{})
 

    //将服务地址注册到etcd中
    fmt.Printf("greeting server address: %s\n", serverAddr)
    register(*EtcdAddr, *ServiceName, serverAddr, 5)
 

    //关闭信号处理
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

    go func() {
        s := <-ch
        unRegister(*ServiceName, serverAddr)
        if i, ok := s.(syscall.Signal); ok {
            os.Exit(int(i))
        } else {
            os.Exit(0)
        }
    }()
 

    //监听服务
    err = srv.Serve(listener)
    if err != nil {
        fmt.Println("监听异常：", err)
        return
    }
}
```

 

三. 客户端（client/main.go）
================
客户端首先需要实现接口resolver.Resolver，其中方法Build()用于创建一个etcd解析器，grpc.Dial()会同步调用该方法，解析器需要根据key前缀监听etcd中服务地址列表的变化并更新本地列表。   
然后注册解析器，创建gRPC句柄，使用轮询负载均衡请求服务。

代码逻辑：
```golang
package main
 
import (
    "flag"
    "fmt"
    proto "go-git/etcd-demo/proto"
    "log"
    "strings"
    "time"
 
    "github.com/coreos/etcd/mvcc/mvccpb"
    "go.etcd.io/etcd/clientv3"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/resolver"
)
 
const schema = "ns"
 
var (
    ServiceName = flag.String("ServiceName", "greet_service", "service name")        //服务名称
    EtcdAddr    = flag.String("EtcdAddr", "127.0.0.1:2379", "register etcd address") //etcd的地址
)
 
var cli *clientv3.Client
 
//etcd解析器
type etcdResolver struct {
    etcdAddr   string
    clientConn resolver.ClientConn
}
 
//初始化一个etcd解析器
func newResolver(etcdAddr string) resolver.Builder {
    return &etcdResolver{etcdAddr: etcdAddr}
}
 
func (r *etcdResolver) Scheme() string {
    return schema
}
 
//watch有变化以后会调用
func (r *etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
    log.Println("ResolveNow")
    fmt.Println(rn)
}
 
//解析器关闭时调用
func (r *etcdResolver) Close() {
    log.Println("Close")
}
 
//构建解析器 grpc.Dial()同步调用
func (r *etcdResolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    var err error
 
    //构建etcd client
    if cli == nil {
        cli, err = clientv3.New(clientv3.Config{
            Endpoints:   strings.Split(r.etcdAddr, ";"),
            DialTimeout: 15 * time.Second,
        })
        if err != nil {
            fmt.Printf("连接etcd失败：%s\n", err)
            return nil, err
        }
    }
 
    r.clientConn = clientConn
 
    go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")
 
    return r, nil
}
 
//监听etcd中某个key前缀的服务地址列表的变化
func (r *etcdResolver) watch(keyPrefix string) {
    //初始化服务地址列表
    var addrList []resolver.Address
 
    resp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
    if err != nil {
        fmt.Println("获取服务地址列表失败：", err)
    } else {
        for i := range resp.Kvs {
            addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(resp.Kvs[i].Key), keyPrefix)})
        }
    }
 
    r.clientConn.NewAddress(addrList)
 
    //监听服务地址列表的变化
    rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
    for n := range rch {
        for _, ev := range n.Events {
            addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
            switch ev.Type {
            case mvccpb.PUT:
                if !exists(addrList, addr) {
                    addrList = append(addrList, resolver.Address{Addr: addr})
                    r.clientConn.NewAddress(addrList)
                }
            case mvccpb.DELETE:
                if s, ok := remove(addrList, addr); ok {
                    addrList = s
                    r.clientConn.NewAddress(addrList)
                }
            }
        }
    }
}
 
func exists(l []resolver.Address, addr string) bool {
    for i := range l {
        if l[i].Addr == addr {
            return true
        }
    }
    return false
}
 
func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
    for i := range s {
        if s[i].Addr == addr {
            s[i] = s[len(s)-1]
            return s[:len(s)-1], true
        }
    }
    return nil, false
}
 
func main() {
    flag.Parse()
 
    //注册etcd解析器
    r := newResolver(*EtcdAddr)
    resolver.Register(r)
 
    //客户端连接服务器(负载均衡：轮询) 会同步调用r.Build()
    conn, err := grpc.Dial(r.Scheme()+"://author/"+*ServiceName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
    if err != nil {
        fmt.Println("连接服务器失败：", err)
    }
    defer conn.Close()
 

    //获得grpc句柄
    c := proto.NewGreetClient(conn)
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        fmt.Println("Morning 调用...")
        resp1, err := c.Morning(
            context.Background(),
            &proto.GreetRequest{Name: "JetWu"},
        )
        if err != nil {
            fmt.Println("Morning调用失败：", err)
            return
        }
        fmt.Printf("Morning 响应：%s，来自：%s\n", resp1.Message, resp1.From)
 

        fmt.Println("Night 调用...")
        resp2, err := c.Night(
            context.Background(),
            &proto.GreetRequest{Name: "JetWu"},
        )
        if err != nil {
            fmt.Println("Night调用失败：", err)
            return
        }
        fmt.Printf("Night 响应：%s，来自：%s\n", resp2.Message, resp2.From)
    }
}
```

 

四. 运行验证
================

1 先确保启动了etcd  

2 使用3个不同端口运行三个服务端：
```sh
go run main.go -Port 3000

go run main.go -Port 3001

go run main.go -Port 3002
```


3 启动客户端：
```sh
go run main.go
```

可以看到，客户端使用轮询的方式对三个服务端进行请求，从而实现负载均衡。