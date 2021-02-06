
目前市面上常见的服务治理有consul，etcd，zookeeper，euerka，我们需要根据自己的服务特点选择自己相对合适的服务治理工具。
```sh
Feature				Consul					zookeeper				etcd			euerka
服务健康检查			服务状态，内存，硬盘等	(弱)长连接，keepalive	连接心跳			可配支持
多数据中心			支持						—						—				—
kv存储服务			支持						支持						支持				—
一致性				raft					paxos					raft			—
cap					ca						cp						cp				ap
使用接口(多语言能力)	支持http和dns			客户端					http/grpc		http（sidecar）

watch支持			全量/支持long polling	支持						支持 long polling	支持 long polling/大部分增量

自身监控				metrics					—						metrics			metrics
安全					acl /https				acl						https支持(弱)	—
spring cloud集成		已支持					已支持					已支持			已支持
```
调研一个工具需要看到其优点，更需要看到其缺点，当服务优点大于自身业务需求缺点，且缺点有对应的解决方案时，我们可以倾向于考虑。


euerka 据说现在已停止维护，决定不考虑使用。

zookeeper 为java开发的，需要java环境，相对比较复杂，优先级较低。

etcd 与consul为go开发，部署简单，功能相对更符合自身业务的需求。

consul监控检查更为丰富，支持多数据中心，webui查看等，配合consul-template实现nginx动态负载均衡等特点更符合自身业务需求，因为决定选用consul作为业务的服务治理工具。



使用consul，其主要有四大特性：
------------------
1. 服务发现：利用服务注册，服务发现功能来实现服务治理。

2. 健康检查：利用consul注册的检查检查函数或脚本来判断服务是否健康，若服务不存在则从注册中心移除该服务，减少故障服务请求。

3. k/v数据存储：存储kv数据，可以作为服务配置中心来使用。

4. 多数据中心：可以建立多个consul集群通过inter网络进行互联，进一步保证数据可用性。


用的最多的：服务注册发现，配置共享。