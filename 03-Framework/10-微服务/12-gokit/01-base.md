
go-kit 入门

go-kit 是一个分布式的开发工具集，在大型的组织（业务）中可以用来构建微服务。其解决了分布式系统中的大多数常见问题，因此，使用者可以将精力集中在业务逻辑上。

go-kit就是一个go语言相关的微服务工具包。它自身称为toolkit，并不是framework。也就是gokit是将一系列的服务集合在一起，提供接口，从而让开发者自由组合搭建自己的微服务项目。基本上看完gokit的例子就可以动手模仿着写一个类似的小项目。

go-kit的结构分为：  
- 传输层
- 端点层
- 服务层


transport（传输层）
--------------
当你构建基于微服务的分布式系统时，服务通常使用HTTP或gRPC等具体传输或使用NATS等pub/sub系统相互通信。Go套件中的传输层绑定到具体运输。Go套件支持使用HTTP，gRPC，NATS，AMQP和Thrift提供服务的各种传输。由于Go kit服务仅专注于实现业务逻辑，并且不了解具体传输，因此你可以为同一服务提供多个传输。例如，可以使用HTTP和gRPC公开单个Go工具包服务。

决定用哪种方式提供服务请求，一般就是 http,rpc


endpoint（端点层）
--------------
端点是服务器和客户端的基本构建块。在Go kit中，主要消息传递模式是RPC。端点表示单个RPC方法。Go工具包服务中的每个服务方法都转换为端点，以便在服务器和客户端之间进行RPC样式通信。每个端点使用传输层通过使用HTTP或gRPC等具体传输将服务方法公开给外部世界。可以使用多个传输来公开单个端点。

是gokit最重要的一个层，是一个抽象的接收请求返回响应的函数类型。在这个定义的类型里面会去调用service层的方法，组装成response返回。而gokit中的所有中间件组件都是通过装饰者设计模式注入的。

//原型：
// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
 
//使用方法：
func(log Logger, in endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req interface{}) (interface{}, error) {
            logger.Log("input", toJSON(req))
            resp, err := in(ctx, req)
            logger.Log("output", toJSON(resp), "err", err)
            return resp, err
    }
}


service（服务层）
--------------
业务逻辑在服务层中实现。Go kit服务被建模为接口。服务中的业务逻辑包含核心业务逻辑，它不应具有端点或HTTP或gRPC等具体传输的任何知识，或者请求和响应消息类型的编码和解码。这将鼓励你遵循基于Go套件的服务的干净架构。每种服务方法都通过使用适配器转换为端点，并使用具体传输进行公开。由于结构简洁，可以使用多个传输来公开单个Go工具包服务。

所有的具体方法写在这里，可以理解为单体web框架中的控制器部分。




Go套件中的中间件
==============
Go kit通过强制分离关注点来鼓励良好的设计原则。使用中间件实现服务和端点的交叉组件。Go kit中的中间件是一种强大的机制，可以包装服务和端点以添加功能（交叉组件），例如日志记录，断路器，速率限制，负载平衡或分布式跟踪。



工具包
==============
这三个层组成一个gokit微服务应用。此外，作为一个工具包，gokit为此提供了很多微服务工具组件。

认证组件（basic, jwt）


回路熔断器  
	Circuitbreaker（回路断路器） 模块提供了很多流行的回路断路lib的端点（endpoint）适配器。回路断路器可以避免雪崩，并且提高了针对间歇性错误的弹性。每一个client的端点都应该封装（wrapped）在回路断路器中。

限流器  
	ratelimit模块提供了到限流器代码包的端点适配器。限流器对服务端（server-client）和客户端（client-side）同等生效。使用限流器可以强制进、出请求量在阈值上限以下。

日志组件  
	服务产生的日志是会被延迟消费（使用）的，或者是人或者是机器（来使用）。人可能会对调试错误、跟踪特殊的请求感兴趣。机器可能会对统计那些有趣的事件，或是对离线处理的结果进行聚合。这两种情况，日志消息的结构化和可操作性是很重要的。Go kit的 log 模块针对这些实践提供了最好的设计。

普罗米修斯监控系统  

Metrics（Instrumentation）度量/仪表盘  
	直到服务经过了跟踪计数、延迟、健康状况和其他的周期性的或针对每个请求信息的仪表盘化，才能被认为是“生产环境”完备的。Go kit 的 metric 模块为你的服务提供了通用并健壮的接口集合。可以绑定到常用的后端服务，比如 expvar 、statsd、Prometheus。

Request tracing（请求跟踪）  
	随着你的基础设施的增长，能够跟踪一个请求变得越来越重要，因为它可以在多个服务中进行穿梭并回到用户。Go kit的 tracing 模块提供了为端点和传输的增强性的绑定功能，以捕捉关于请求的信息，并把它们发送到跟踪系统中。（当前支持 Zipkin，计划支持Appdash


服务发现系统接口（etcd, consul等）  
	如果你的服务调用了其他的服务，需要知道如何找到它（另一个服务），并且应该智能的将负载在这些发现的实例上铺开（即，让被发现的实例智能的分担服务压力）。Go kit的loadbalancer模块提供了客户端端点的中间件来解决这类问题，无论你是使用的静态的主机名还是IP地址，或是 DNS的 SRV 记录，Consul，etcd 或是 Zookeeper。并且，如果你使用定制的系统，也可以非常容易的编写你自己的 Publisher，以使用 Go kit 提供的负载均衡策略。（目前，支持静态主机名、etcd、Consul、Zookeeper）


路由跟踪  



这些组件大大方便了我们开发一个微服务应用。








入口demo  
==============
https://github.com/mycodesmells/gokit-example  

https://github.com/FengGeSe/demo


https://github.com/kplcloud/kplcloud   一个基于了kubernetes的应用管理系统  
